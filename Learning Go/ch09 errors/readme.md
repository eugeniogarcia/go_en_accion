## sentinel errors

Se trata de errores que estan definidos a nivel de paquete y que señalan que la ejecución debe interrumpirse. En el paquete zip, por ejemplo, se definen algunos sentinel errors, estos cuatro:

```go

package zip

var (
	ErrFormat       = errors.New("zip: not a valid zip file")
	ErrAlgorithm    = errors.New("zip: unsupported compression algorithm")
	ErrChecksum     = errors.New("zip: checksum error")
	ErrInsecurePath = errors.New("zip: insecure file path")
)

[...]

func (r *checksumReader) Read(b []byte) (n int, err error) {
	if r.err != nil {
		return 0, r.err
	}
	n, err = r.rc.Read(b)
	r.hash.Write(b[:n])
	r.nread += uint64(n)
	if r.nread > r.f.UncompressedSize64 {
		return 0, ErrFormat
	}

[...]

```

si en nuestro código usamos este paquete, podemos ver si el error obtenido es uno de los sentinel:

```go
_, err := zip.NewReader(notAZipFile, int64(len(data)))
if err == zip.ErrFormat {
    fmt.Println("Told you so")
}
```

## Wrap Errors

### Teoría

Esta feature consiste en "añadir" a un error otro error o errores. Supongamos que tenemos un error A, y queremos producir otro error B que "incorpore" el error A. Lo que haremos es wrappear A con el nuevo error. `B:=fmt.Errorf("[mensaje de error correspondiente a B] %w", A)`. También es posible wrappear varios errores. Supongamos que tenemos A1, A2, A3, haríamos `B:=fmt.Errorf("[mensaje de error correspondiente a B] %w %w %w", A1, A2, A3)`. El convenio es incluir los errores, los `%w` al final. 

Cuando wrapeamos un error tenemos que incluir en el error un método `Unwrap() error` que deshaga la operación. El error `B` que se ha generado con el `fmt.Errorf` al usar `%w` incluye ya este método. Podemos ver como esta implementado en el paquete `errors`. Efectivamente se definen estos dos tipos:

```go
type wrapError struct {
	msg string
	err error
}

func (e *wrapError) Error() string {
	return e.msg
}

func (e *wrapError) Unwrap() error {
	return e.err
}

type wrapErrors struct {
	msg  string
	errs []error
}

func (e *wrapErrors) Error() string {
	return e.msg
}

func (e *wrapErrors) Unwrap() []error {
	return e.errs
}
```

__Notese que como el receptor es `(e *wrapErrors)`, el tipo que implementa el interface `error` es el puntero. Veremos que efectivamente `fmt.Errorf` devuelve un puntero a estos tipos.__

La funcion `fmt.Errorf` devuelve:

- Si no hay `%w` `err = errors.New(s)`, donde s es el mensaje. Es decir, retorna una implementacion de la interface `error` _normal_.

- Si hay un `%w`:

    ```go
    w := &wrapError{msg: s} //el puntero de wrapError es un error
    w.err, _ = a[p.wrappedErrs[0]].(error) //comprobamos que lo que estamos pasando con argumento sea un error
    err = w
    ```

    Es decir, retorna un puntero a `wrapError` de modo que se implementa la interface `error` al tiempo que se incluye el método `Unwrap() error`.

- Si hay más de un `%w`:

    ```go
    err = &wrapErrors{s, errs}
    ```

    Es decir, retorna un puntero a `wrapErrors` de modo que se implementa la interface `error` al tiempo que se incluye el método `Unwrap() []error`.

    Añadir que hay una funcion en el paquete `errors` que nos permite "fusionar" varios errores, `errors.Join(slice_de_errores)`. En nuestro ejemplo la utilizamos para implementar `Error() string` el unwrapp del slice de errores:

    ```go
    func (m MyError) Error() string {
        return errors.Join(m.Errors...).Error() //usamos errors.Join para unir los mensajes de los errores wrappeados
    }
    ``` 


Y para completar la feature, dos cosas:

- Para extraer un error empaquetado podemos usar `A:=errors.Unwrap(B);`. La implementacion del método `Unwrap` que se incluye en el paquete errors es la siguiente:

```go
func Unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})// se verifica que err implemente una interface anónima que incluye el método `Unwrap() error`, y si lo hace se llama a dicho método para desempaquetar
	if !ok {
		return nil
	}
	return u.Unwrap()
}
```

- Si solo queremos incorporar el mensaje del error `B` con `fmt.Errorf`, pero no el error propiamente dicho, en lugar de usar el selector `%w` usaremos el selector `%v`. Lo que tenemos es un error que incluye el mensaje de error del error `B`, pero no se ha _empaquetado_ `B` como tal.

### New

Por completar, esta es la implementación del método `New(string) error` del paquete `errors`. 

```go
func New(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}
```

### Ejemplo

Tenemos un ejemplo en `wrapped`.

## errors.As, errors.Is

Estos métodos sirben para comprobar si un error **o entre los errores _empaquetados_ en el error** tenemos un error determinado (`errors.Is`) o uno que implemente un determinado tipo (`errors.As`).

### errors.Is

Con `errors.Is` podemos buscar un error concreto en otro. La función utiliza `==` para encontrar un error entre los errores "empaquetados" dentro de un error. En el ejemplo `is_error` podemos ver un ejemplo de uso de esta función. La función se usa típicamente para buscar __sentinel errors__.

La limitación que tiene esta función es que utiliza `==` para buscar un error. Esto hace que si el error tiene algun elemento no comparable en su definición, por ejemplo un slice, no vayamos a encontrar el error con `errors.Is`. Podemos entonces usar dos técnicas:

- emplear `errors.As`
- implementar un método `Is(target error) bool` que evalue la igualdad. Si tenemos este método implementado cuando cuando hacemos `errors.Is` se usará este método para evaluer en lugar de `==`

### errors.As

Con `errors.As` podemos comprobar si un error implementa un tipo determinado. La forma de usar la función es pasar un puntero a una variable. El método comprobará si el error implementa el tipo de la variable. __Esto abre las posibilidades, no solo para buscar si nuestro error incluye un tipo determinado, como para verificar si en nuestro error hay algun error que implemente una interface concreta__.

## Panic y recover

Este mecanismo es similar a la gestión de excepciones disponible en otros lenguajes, pero la filosofía de diseño de Go es utilizar este mecanismo para el apagado controlado, esto es, cuando se crea un panic la intención en la casi totalidad de los casos debe ser la de señalar que hay un problema que impide que continue el procesamiento y que se debe detener el programa. En el caso de que estemos construyendo una librería publica si debemos interceptar el panic y convertirlo en un error para que el consumidor de la librería - del que no tenemos control - gestione lo que debe suceder a continuación.

Podemos también generar panics desde nuestro código con la instrucción `panic(objeto)`.

Un panic puede interceptarse en un `defer()` utilizando la función `recover()`. Esta función obtiene el objeto que se paso a `panic`.

En los ejemplos `panic` y `panic_recover` podemos ver como se usa