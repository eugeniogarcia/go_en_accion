# imagen base, con Go y Alpine Linux
# Start from golang alpine base image
FROM golang:1.25-alpine

# Creamos el directorio de trabajo en la imagen
WORKDIR /app

# copiamos todo al directorio de trabajo
COPY . .

# descargamos las dependencias
RUN go mod download

# construimos la aplicación
RUN go build -o runners-app main.go

# exponemos el puerto 8080
EXPOSE 8080

# definimos una variable de entorno llamada ENV con el valor k8s
ENV ENV=k8s

# Comando para ejecutar cuando se inicie el contenedor
CMD ["/app/runners-app"]
