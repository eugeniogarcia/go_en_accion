El contexto nos permite compartir datos entre diferentes gorutinas. El contexto esta definido como un interface:

```go
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
}
```

En el contexto podemos guardar key-values. Los key y los values serán de tipo `any` (`any` es una interface definida como `type any = interface{}`, de modo que cualquier tipo la _implementa_). Recuperamos el valor de un key usando el método `Value`.

Con el método `Done` se obtiene un canal de lectura. Cuando el contexto ha _terminado_ se señala cerrando el canal.

Con el método `Deadline` se recupera la deadline definida para el contexto, y si se ha alcanzado ya o no.

La forma de operar con un contexto es que:

- se crea. Para ello hay una serie de funciones helper en el paquete `context`
- sobre un contexto ya existente, se construye otro que es un wrapper del original y aporta un _extra_


tipicamente podemos crear un contexto _pelado_ con:

```go
context.Background()
```

la forma de crear un wrapper alrededor de un contexto es usando un método que por convenio se llama `WithXXXX`. Por ejemplo, para crear un contexto cancelable haremos:

```go
ctx, cancelFunc := context.WithCancel(context.Background())
```

la función nos devuelve el contexto - `ctx` -, y una función helper para iterar con él - en este ejemplo `cancelFunc`, que se encargaría de liberar los recursos, y cerrar el canal que devuelve `Done()`.

Típicamente el contexto se pasará como primer argumento a aquellas funciones que lo usen, y típicamente por convenio se suele denominar al argumento `ctx`

```go
func logic(ctx context.Context, info string) (string, error) {
    // do some interesting stuff here
    return "", nil
}
```

## Uso con http

### cliente

Usaremos el servicio `http://httpbin.org` para probar nuestro cliente http.

En `cancel_http` tenemos un ejemplo de cliente http que utiliza un contexto. Creamos una request con un contexto determinado (primero creamos el contexto, y luego se lo pasamos a la funcion helper `http.NewRequestWithContext`). Si durante la ejecución de la petición cancelaramos el contexto, por ejemplo, el cliente http se detendría:

```go
func makeRequest(ctx context.Context, url string) (*http.Response, error) {
	// Creamos una petición HTTP con el contexto proporcionado
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// Realizamos la petición HTTP usando el cliente por defecto
	return http.DefaultClient.Do(req)
}
```

en `cancel_error_http` tenemos una evolución del ejemplo anterior en la que usamos un contexto cancelable, pero en el que podemos indicar el motivo de cancelación:

```go
	// Creamos un contexto que es cancelable y en el que podemos indicar la razón de la cancelación. cancelFunc adminte un argumento en el que se indica la razón de cancelación
	ctx, cancelFunc := context.WithCancelCause(context.Background())
	// Aseguramos que se llame a cancelFunc siepre para liberar recursos
	defer cancelFunc(nil)
```

cuando cancelamos pasamos la razón a la función de cancelación:

```go
cancelFunc(fmt.Errorf("in status goroutine: %w", err))
```

y cuando necesitamos saber cual fue la razón por la que se canceló el contexto usamos la función `context.Cause(ctx)`:

```go
case <-ctx.Done():
    // Podemos recuperar la razón de la cancelación usando context.Cause
    fmt.Println("in main: cancelled with error", context.Cause(ctx))
```

### servidor

El paquete que gestiona los clientes y los servidores http se desarrollo antes de que se incorporase el contexto en go, y por este motivo no sigue el convenio que comentaba antes, de pasar el contexto como argumento en las funciones que lo usan. El patrón que tenemos que aplicar en estos casos podemos observarlo en el ejemplo context_patterns`. Veremos varios casos de uso

- Uso en un Middleware
- Uso en la lógica de negocio de un backend
- Propagar el contexto a otros servicios - llamados por el backend

#### Middleware

Como enriquecer el contexto en un middleware, de modo que pueda usarse ya modificado en el resto de la cadena

- Paso 1. _Extraemos el contexto_ de la request http
- Paso 2. _Incorporamos_ al contexto lo que necesitemos usando el método `context.WithValue(ctx, key, valor)` (**notese que es una helper function, no es parte de la interface que define un contexto**). Esto crea un contexto que wrappea el contexto original. Podemos recuperar valores del contexto con `ctx.Value(key)`
- Paso 3. _Actualizamos la request_ con el nuevo contexto
- Paso 4. Continua la ejecución

```go
// Middleware. Toma un handler y devuelve un handler
func Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Paso 1. Obtenemos de la request el contexto
		ctx := req.Context()
		// Paso 2. Hacemos algo con el contexto
		// Paso 3. Asociamos a la request este contexto
		req = req.WithContext(ctx)

		// Paso 4. Llamamos al siguiente handler de la cadena, pero a partir de este punto el contexto ya esta enriquecido
		handler.ServeHTTP(rw, req)
	})
}
```

#### Backend

Eventualmente la razón de ser del contexto es utilizarlo. 

```go
func handler(rw http.ResponseWriter, req *http.Request) {
	// Paso 1. Obtenemos de la request el contexto
	ctx := req.Context()

[...]
}
```

#### Propagar

Podemos propagar el contexto a otros servicios, de forma que se compartan dataos y se cree, ..., un contexto de ejecución. Por ejemplo, para implementar trazing.

```go
func (sc ServiceCaller) callAnotherService(ctx context.Context, data string) (string, error) {
	// Paso 1. Creamos la request HTTP, asociando el contexto (no creamos un contexto nuevo, sino que usamos el que nos pasan)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"http://example.com?data="+data, nil)
	if err != nil {
		return "", err
	}
	resp, err := sc.client.Do(req)
	if err != nil {
		return "", err
	}

[...]
}
```

## key-values

El key y el value son de tipo `any`, pero para key debemos utilizar un tipo que sea comparable. Cuando añadimos un key value a un contexto lo que hacemos es crear un wrapper alrededor del contexto anterior en el que se añaden los key. Si usamos un key que ya existia lo enmascaremos. Al objeto de garantizar que el key sea único podemos sar estas dos técnicas (extraidas del ejemplo `context_user`):

- definimos una key, que será unica - para evitar el enmascaramiento -, y que no exportamos:

```go
// Declaramos un tipo con base int para usarlo como key en el contexto
type userKey int

// Definimos la lista de keys usando iota. En este caso solo tenemos una key. Con esto garantizamos la unicidad de la key. key no se exporta es solo para uso interno.
const (
	_ userKey = iota
	key
)
```

- usamos el helper `context.WithValue()` para crear una nueva key-value. En este ejemplo la key que se usa es algo interno al paquete. El usuario no vera ni tiene que usar la key:

```go
func ContextWithUser(ctx context.Context, user string) context.Context {
	// Usamos la helper function context.WithValue para crear un nuevo contexto que wrappea ctx e incluye una nueva key/valor
	return context.WithValue(ctx, key, user)
}
```

- recuperamos el valor usado el método del contexto `Value()`:

```go
func UserFromContext(ctx context.Context) (string, bool) {
	// Recuperamos un valor. Como es el valor es un interface - any - podemos hacer type assertion. En este caso esperamos un string. ok sera true si la aserción ha ido bien. En user tendremos el valor asociado a la key
	user, ok := ctx.Value(key).(string)
	return user, ok
}
```

## timeout

Podemos crear un contexto con timeout. Lo que nos permite hacer este contexto es a) cancelar - como un cancellable - y b) definir un timeout. Cuando se alcanza el timeout se como si hubieramos cancelado el contexto:

```go
// Creamos un contexto padre con timeout de 3 segundos
ctx, cancelFuncParent := context.WithTimeout(context.Background(), 3*time.Second)
// Aseguramos que almenos se llame a cancel una vez para liberar recursos
defer cancelFuncParent()
```

podemos encadenar varios contextos. Por ejemplo, si queremos poder indicar un motivo de cancelación, y al mismo tiempo definir un timeout:

```go
// Creamos un contexto padre con timeout de 3 segundos
ctx, cancelFuncParent := context.WithTimeout(context.Background(), 3*time.Second)
// Aseguramos que almenos se llame a cancel una vez para liberar recursos
defer cancelFuncParent()
// Creamos un contexto hijo que puede ser cancelado con un motivo
ctx, cancelFunc := context.WithCancelCause(ctx)
// Aseguramos que almenos se llame a cancel una vez para liberar recursos
defer cancelFunc(nil)
```

podemos ver un ejemplo de esto en `timeout_error_http`.

Si definimos varios timeouts en el mismo contexto, el primero que se dispare cancelará la ejecución:

```go
func main() {
	ctx := context.Background()
	parent, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	child, cancel2 := context.WithTimeout(parent, 3*time.Second)
	defer cancel2()
	start := time.Now()
	<-child.Done()
	end := time.Now()
	fmt.Println(end.Sub(start).Truncate(time.Second))
}
```