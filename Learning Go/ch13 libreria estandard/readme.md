## `io.Reader`, `io.Writer`

`io_friends`. Muestra como el reader y el writer son una interface simple, y como en la librería estandar hay diferentes métodos que aceptan este tipo de interfaces

## JSON

`custom_json`. Muestra como serializar y deserializar un tipo Go a un JSON usando `Marshal` y `Unmarshal` para convertir un slice de bytes a un tipo y viceversa. Hemos usado las versiones incluidas en el paquete `json` pero de forma análoga podríamos usar las de `xml`, o otros formatos. La estructura a la que se deserializa o serializan los datos tiene anotaciones con `json:"[nombre campo]"` en la que se indica el nombre del atributo en el json, si debe o no filtrarse del json, `json:"-"`, y si debe o no incluirse cuando el dato esta vacio, `json:"[nombre campo],omitempty"`. 

```go
type Order struct {
	ID          string    `json:"id"`
	Items       []Item    `json:"items"`
	NoLoquiero  bool      `json:"-"` //ignoramos este campo
	DateOrdered time.Time `json:"date_ordered"`
	CustomerID  string    `json:"customer_id,omitempty"` //si el campo esta vacio lo ignoramos
}

[...]

var o Order
err := json.Unmarshal([]byte(data), &o) //deserializamos, en el objeto de tipo Order
if err != nil {
    panic(err)
}

[...]

out, err := json.Marshal(o) //serializamos, de nuevo a JSON
if err != nil {
    panic(err)
}
```

En el ejemplo `json` se muestra una alternativa a usar byte slices, usar `io.Reader` o `io.Writer` usando `json.NewDecoder` o `json.NewEncoder`. Estos métodos pueden usarse para llamar `Decode(&)` o `Encode()` respectivamente. Estos métodos pueden llamarse una vez o de forma iterativa hasta llegar al `io.EOF`. En `encode_decode` tenemos un ejemplo para leer una lista de jsons.

```go
err = json.NewEncoder([implemente io.Writer]).Encode([tipo])
if err != nil {
    return err
}

[...]

err = json.NewDecoder([implemente io.Reader]).Decode(&[tipo])
```

o podemos usarlo para procesar un stream de valores:

```go
decoder := json.NewDecoder([implemente io.Reader]) //decoder
for { //bucle infinito
    err:=decoder.Decode(&[tipo]) //recuperamos el siguiente elemento
    if (err=!nil){
        if errors.Is(err, io.EOF) { //comprobamos si el error es o tiene wrappeado io.EOF
		    break
		}
		panic(err)
    }
    [hacer algo con tipo]
}
```

Comentar que habitualmente usaremos estos tipos como fuente de datos:

- `bytes.NewBuffer([]byte)`. Crea un `io.Reader`, `io.Writer`
- `bytes.NewReader([]byte)`. Crea un `io.Reader`
- `string.NewBuffer(string)`. Crea un `io.Reader`, `io.Writer`
- `string.NewReader(string)`. Crea un `io.Reader`

## http

### Cliente

En `client` tenemos un ejemplo de cliente http/2. Esta estructura esta diseñada para se compartida a lo largo del programa y sus gorutinas. Típicamente se crea una vez, y se usa el método `Do` o alguno de sus métodos espciales `Get`, `Post`, `PostForm`, `Head` o `CloseIdleConnections`. `Get`, `Post`, `PostForm` envian `http.NewRequestWithContext`. El `http.NewRequestWithContext` se construye con cabeceras, uri, y opcionalmente un tipo que se envia como payload. 

### Servidor

En `server` tenemos un ejemplo de servidor http/2. Un servidor es una estructura de tipo `http.Server`:

```go
s := http.Server{
    Addr:         ":8080", //host y puerto en el que escuchamos
    ReadTimeout:  30 * time.Second,
    WriteTimeout: 90 * time.Second,
    IdleTimeout:  120 * time.Second,
    Handler:      HolaMundoHandler{}, //handler
}
err := s.ListenAndServe()
```

se le pasa un handler para tratar las peticiones. El handler implementa la interface `Handler`: 

```go
type Handler interface {
    ServeHTTP(http.ResponseWriter, *http.Request)
}
```

`http.Request` es el mismo tipo que se envio desde el cliente. `http.ResponseWriter` es la siguiente interface:

```go
type ResponseWriter interface {
        Header() http.Header
        Write([]byte) (int, error)
        WriteHeader(statusCode int)
}
```

Con `Header()` podemos especificar cabeceras en la respuesta. con `WriteHeader` especificamos el _http status code_ de la respuesta. Por último `Write` es el metodo que envia la respuesta al cliente.

### Multiples rutas

En `server_mux` tenemos un router de peticiones. `http.NewServeMux` es un handler que podemos usar en un `http.server`, pero que nos permite a su vez definir diferentes handlers para diferentes rutas. Podemos asociar a un `http.NewServeMux` otro `http.NewServeMux` de modo que se anidan rutas. La función _helper_ `http.StripPrefix([prefijo],[handler])` quita de la ruta el _prefijo_ de modo que al _hadler_ le llegará una ruta en la que ya no está el prefijo.

Esta funcionalidad es análoga al mecanismos de _Router_ de una aplicación _App Express_ de node.

### Middleware

También podemos implementar un middleware, en `middleware` tenemos un ejemplo. Un middleware recibe un handler y devuelve otro. El tratamiento que hace el middleware sera aplicar una cierta lógica a resultas de la cual se decide continuar con el pipeline, llamando a `h.ServeHTTP(w, r)`, o se aborta la ejecución con `return` 

## Logs

Tenemos dos paquetes, el `log` clásico, y el `slog` para trabajar con strcutured logs, es decir, logs en los que enviamos información más o menos estructurada para facilitar su tratamiento posterior.

`slog` permite un uso básico como `log`:

```go
// Métodos simples para crear un log
slog.Debug("debug log message")
slog.Info("info log message")
slog.Warn("warning log message")
slog.Error("error log message")
```

también podemos enviar parejas de key/values:

```go
userID := "fred"
loginCount := 20
slog.Info("user login",
    "id", userID, //primer par key/value
    "login_count", loginCount) //segundo key/value
```

si necesitamos guardar la información de una forma más estructurada usaremos un handler. En el paquete slog se incluye `NewJSONHandler` para crear un handler que guarde la información como un json:

```go
// Si necesitamos enviar la información estructurada, por ejemplo con un json, creamos un handler
//1. definimos las opciones del handler
options := &slog.HandlerOptions{Level: slog.LevelDebug}
//2. creamos el handler, en este caso un JSONHandler. Para crearlos indicamos un io.Writer, y las opciones
handler := slog.NewJSONHandler(os.Stderr, options)
//3. creamos el logger a partir del handler
mySlog := slog.New(handler) //a partir del handler creamos un logger, que pasamos a utilizar con los métodos estadard
```

a partir de aquí ya podemos usar nuestro logger de la forma habitual:

```go
lastLogin := time.Date(2023, 01, 01, 11, 50, 00, 00, time.UTC)
mySlog.Debug("debug message", "id", userID, "last_login", lastLogin)
```

Para optimizar el rendimiento del logger podemos usamos `LogAttrs` en lugar de Info, Debug, etc:

```go
// Toma un contexto, y la información a registrar. Los pares key-value se crear con helpers. slog.Any nos sirve para cualquier tipo
ctx := context.Background()
mySlog.LogAttrs(ctx, slog.LevelInfo, "faster logging", slog.String("id", userID), slog.Time("last_login", lastLogin))
```

es posible crear una logger a medida con `slog.NewLogLogger(mySlog.Handler(), slog.LevelDebug)`.
