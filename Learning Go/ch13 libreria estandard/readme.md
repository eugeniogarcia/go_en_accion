- `io_friends`. Muestra como el reader y el writer son una interface simple, y como en la librería estandar hay diferentes métodos que aceptan este tipo de interfaces

- `custom_json`. Muestra como serializar y deserializar un tipo Go a un JSON usando `Marshal` y `Unmarshal` para convertir un slice de bytes a un tipo y viceversa. Hemos usado las versiones incluidas en el paquete `json` pero de forma análoga podríamos usar las de `xml`, o otros formatos. La estructura a la que se deserializa o serializan los datos tiene anotaciones con `json:"[nombre campo]"` en la que se indica el nombre del atributo en el json, si debe o no filtrarse del json, `json:"-"`, y si debe o no incluirse cuando el dato esta vacio, `json:"[nombre campo],omitempty"`. 

En el ejemplo `json` se muestra una alternativa a usar byte slices, usar `io.Reader` o `io.Writer` usando `json.NewDecoder` o `json.NewEncoder`. Estos métodos pueden usarse para llamar `Decode(&)` o `Encode()` respectivamente. Estos métodos pueden llamarse una vez o de forma iterativa hasta llegar al `io.EOF`. En `encode_decode` tenemos un ejemplo para leer una lista de jsons.

Comentar que habitualmente usaremos estos tipos como fuente de datos:

- `bytes.NewBuffer([]byte)`. Crea un `io.Reader`, `io.Writer`
- `bytes.NewReader([]byte)`. Crea un `io.Reader`
- `string.NewBuffer(string)`. Crea un `io.Reader`, `io.Writer`
- `string.NewReader(string)`. Crea un `io.Reader`

- En `client` tenemos un ejemplo de cliente http/2. Esta estructura esta diseñada para se compartida a lo largo del programa y sus gorutinas. Típicamente se crea una vez, y se usa el método `Do` o alguno de sus métodos espciales `Get`, `Post`, `PostForm`, `Head` o `CloseIdleConnections`. `Get`, `Post`, `PostForm` envian `http.NewRequestWithContext`. El `http.NewRequestWithContext` se construye con cabeceras, uri, y opcionalmente un tipo que se envia como payload. 

- En `server` tenemos un ejemplo de servidor http/2. Un servidor es una estructura de tipo `http.Server`:

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

- En `server_mux` tenemos un router de peticiones. `http.NewServeMux` es un handler que podemos usar en un `http.server`, pero que nos permite a su vez definir diferentes handlers para diferentes rutas. Podemos asociar a un `http.NewServeMux` otro `http.NewServeMux` de modo que se anidan rutas. La función _helper_ `http.StripPrefix([prefijo],[handler])` quita de la ruta el _prefijo_ de modo que al _hadler_ le llegará una ruta en la que ya no está el prefijo.

Esta funcionalidad es análoga al mecanismos de _Router_ de una aplicación _App Express_ de node.

- También podemos implementar un middleware, en `middleware` tenemos un ejemplo. Un middleware recibe un handler y devuelve otro. El tratamiento que hace el middleware sera aplicar una cierta lógica a resultas de la cual se decide continuar con el pipeline, llamando a `h.ServeHTTP(w, r)`, o se aborta la ejecución con `return` 

