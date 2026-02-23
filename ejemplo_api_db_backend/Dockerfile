########################################
# Etapa 1: Build (builder)
########################################
# Imagen base con Go (Alpine, musl). Rápida y pequeña
FROM golang:1.25-alpine AS imagen_constructora

# Herramientas necesarias (git para módulos privados si aplica)
RUN apk add --no-cache git ca-certificates

# creamos el directorio de trabajo
WORKDIR /src

# no copiamos toda la aplicación todavía. Copiamos aquellas partes menos susceptibles a cambios, de modo que es más improbable que esta capa tenga que reconstruirse a menudo. Esto optimiza el uso de la caché de Docker y mejora los tiempos de construcción.
# Optimiza caché de módulos copiando primero los manifests
COPY go.mod go.sum ./

# Descarga dependencias
RUN go mod download

# Ahora ya copiamos el resto del código
COPY . .

# Variables para cross-compilación y binario estático
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ENV CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH

# Compila el binario, reduciendo tamaño
RUN go build -trimpath -ldflags "-s -w" -o /out/runners-app ./main.go


########################################
# Etapa 2: Runtime (mínimo, no-root)
########################################
# Distroless estático, seguro y muy pequeño (sin shell)
FROM gcr.io/distroless/static:nonroot

# Creamos directorio de aplicación
WORKDIR /app

# Copiamos binario
COPY --from=imagen_constructora /out/runners-app /app/runners-app

# Copiamos archivos de configuración necesarios
COPY --from=imagen_constructora /src/runners.toml /app/runners.toml
COPY --from=imagen_constructora /src/runners-k8s.toml /app/runners-k8s.toml

# Puerto para exponer la API REST
EXPOSE 8080
# Puerto para métricas de Prometheus
EXPOSE 9000

# Variable de entorno para seleccionar el archivo de configuración
ENV ENV=k8s

# Ejecuta como usuario no root
USER nonroot

# Entrada del contenedor
ENTRYPOINT ["/app/runners-app"]
