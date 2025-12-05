los tests se definen en archivos con el sufijo `_test.go`. Dentro de este archivo cada caso de prueba es una función con la firma `func Test_xxxxxx(t *testing.T)`. Para ejecutar los casos de prueba se hace `go test` en la linea de comandos.

```go
import "testing"

func Test_addNumbers(t *testing.T) {
	result := addNumbers(2, 3)
	if result != 5 {
		t.Error("incorrect result: expected 5, got", result)
	}
}
```

en `testing.T` tenemos varios métodos para indicar el resultado de la prueba

- `t.Error`
- `t.Errorf`
- `t.Fatal`
- `t.Fatalf`

## Setup & Teardown

Si necesitamos realizar una tarea antes y/o después de ejecutar la batería de tests creamos un método `func TestMain(m *testing.M)`. Cuando este método esta presente al hacer `go test` no se lanzan los test, pero se ejecuta este método. Én el método hacemos `m.Run()` cuando queramos lanzar la batería de tests:

```go
var testTime time.Time

func TestMain(m *testing.M) {
	fmt.Println("Set up stuff for tests here")
	testTime = time.Now()
	exitVal := m.Run()
	fmt.Println("Clean up stuff after tests here")
	os.Exit(exitVal)
}

func TestFirst(t *testing.T) {
	fmt.Println("TestFirst uses stuff set up in TestMain", testTime)
}

func TestSecond(t *testing.T) {
	fmt.Println("TestSecond also uses stuff set up in TestMain", testTime)
}
```

## Limpiar recursos

Cuando necesitamos crear recursos - datos, por ejemplo - para ejecutar los tests, podemos crear el dato como parte del test. 

Con `t.Cleanup()` registramos funciones que se ejecutarán cuando termine la ejecución del test - del caso de prueba. Las funciones que se pasen a Cleanup, si se pasa más de una, se ejecutan como el defer, en modo LIFO.

```go
t.Cleanup(func() {
    fmt.Printf("Limpia el archivo que hemos creado\n")
    os.Remove(f.Name())
})
```

cuando el recurso que estamos creando es un archivo, si creamos los archivos dentro de un directorio creado con `t.TempDir()`, no hace falta especificar el `t.Cleanup()` porque go se encarga de registrar automáticamente un cleanup que borra todo el contenido de este directorio.

Tenemos un ejemplo de estas técnicas en `cleanup`

### Variables de entorno

Es una práctica habitual utilizar variables de entorno para configurar como debe comportarse una aplicación. Para crear estas variables de entorno durante los tests - y luego eliminarlas -, se incluye en `testing.T` ña función `t.Setenv`:

```go
func TestEnvVarProcess(t *testing.T) {
    // Crea la variable de entorno
    t.Setenv("OUTPUT_FORMAT", "JSON")


	cfg := ProcessEnvVars()
	if cfg.OutputFormat != "JSON" {
		t.Error("OutputFormat not set correctly")
	}
	
    // La variable de entorno OUTPUT_FORMAT se reseteará automáticamente al terminar el caso de prueba
}
```

### Datos de prueba

Hay una carpeta que se usa por convenio, `testfata` para guardar datos de prueba. Esta carpeta se ubica en el directorio raiz del paquete con los casos de prueba. En el ejemplo `text` podemos ver como se usa.

## Comparar

Ejemplo `cmp`.

Cuando hay que comparar estructuras podemos utilizar una librería de Google, `github.com/google/go-cmp/cmp`. La librería podemos usarla de dos formas:

- comparar directamente. El resultado indica cuales son los deltas entre las dos estructuras

```go
// Comparamos las dos estructuras completas. El resultado es un texto en el que se indican las diferencias.
if diff := cmp.Diff(esperado, deseado); diff != "" {
	t.Error(diff) // Si hay diferencias, fallamos el test mostrando las diferencias, y continuamos con otro test.
}
```

- comparar usando una función de comparación. La función implementa la lógica de comparación

```go
// Comparamos las dos estructuras, pero esta vez usando una función de comparación que no tiene en cuenta el campo DateAdded.
comparer := cmp.Comparer(func(x, y Person) bool {
	return x.Name == y.Name && x.Age == y.Age
})

// el resultado de la comparación es un texto en el que se indican las diferencias.
if diff := cmp.Diff(esperado, deseado, comparer); diff != "" {
	t.Error(diff) // Si hay diferencias, fallamos el test mostrando las diferencias, y continuamos con otro test.
}
```

## Definir tests programáticamente

En el ejemplo `table` mostramos como crear programáticamente casos de prueba - una tabla de casos de prueba. Podemos usar `t.Run(d.name, func(t *testing.T) {...})` para ejecutar un caso de prueba definido en la función. Combinando esto con _closures_ podemos definir de forma flexible una batería de casos de prueba que son _iguales_ excepto por los parámetros/combinaciones que probamos. Por ejemplo, en este caso definimos una slice con todas las combinaciones a probar y su resultado esperado:

```go
data := []struct {
	name     string //nombre del caso de prueba
	num1     int	//argumento 1
	num2     int	//argumento 2
	op       string	//operacion	
	expected int	//resultado esperado
	errMsg   string	//mensaje de error en caso de no obtener el resultado esperado
}{
	{"addition", 2, 2, "+", 4, ""},
	{"subtraction", 2, 2, "-", 0, ""},
	{"multiplication", 2, 2, "*", 4, ""},
	{"division", 2, 2, "/", 1, ""},
	{"bad_division", 2, 0, "/", 0, `division by zero`},
}
```

luego en un loop lanzamos los casos de prueba:

```go
for _, d := range data {
	t.Run(d.name, func(t *testing.T) {
		result, err := DoMath(d.num1, d.num2, d.op)
		if result != d.expected {
			t.Errorf("Expected %d, got %d", d.expected, result)
		}
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		}
		if errMsg != d.errMsg {
			t.Errorf("Expected error message `%s`, got `%s`",
				d.errMsg, errMsg)
		}
	})
}
```

## Lanzar casos en paralelo

Por defecto los casos de prueba se ejecutan de forma secuencial. Si queremos que el caso de prueba se ejecute en paralelo bastará con marcarlo en la primera línea con la llamada a `t.Parallel()`:

```go
t.Run(d.name, func(t *testing.T) {
	t.Parallel() // marcamos el caso de prueba, de modo que se ejecuta en paralelo. Se lanza lo que sigue a continuación, y sin esperar se ejecuta el siguiente caso de prueba
	fmt.Println(d.input, d.output)
	out := toTest(d.input)
	if out != d.output {
		t.Error("didn't match", out, d.output)
	}
})
```

## Covertura de código

Para comprobar la cobertura de código ejecutaremos los tests de la siguiente forma:

```ps
go test -v -cover
```

podemos guardar los resultados de la covertura con:

```ps
go test -v -cover -coverprofile c.out
```

con go se incluye la `cover tool`, que nos muestra formateado en una página html los resultados de la covertura:

```ps
go tool cover -html c.out
```

## Fuzzing

Fuzzing es una técnica que trata de ejecutar el código con diferentes permutaciones de datos de prueba, con la finalidad de detectar comportamientos anómalos que tienen su origen en la naturaleza del dato (por ejemplo, podemos tener una cobertura del 100% y seguir teniendo errores en el código, errores que solo se manifiestan con ciertos datos).

En go podemos incluir un caso, y solo uno, de Fuzzing, en un método con la siguiente firma: `func Fuzzxxxx(f *testing.F) {`. Típicamente lo que haremos es utilizar `*testing.F` de la siguiente forma:

- crear un seed de datos (_seed corpus_). Se alimentan con el método `f.Add(arg1, arg2, arg3)` juegos de datos. Tomando este juego como punto de partida, el Fuzzer creará diferentes permutaciones de datos

```go
// crea un slice de slices de bytes como see de datos
testcases := [][]byte{
	[]byte("3\nhello\ngoodbye\ngreetings\n"),
	[]byte("0\n"),
}
// crea el seed de datos
for _, tc := range testcases {
	f.Add(tc) // Use f.Add to provide a seed corpus
}
```

- Definir el caso de prueba con `f.Fuzz(func(t *testing.T, arg1, arg2, arg3) {`

```go
f.Fuzz(func(t *testing.T, in []byte) {
	r := bytes.NewReader(in)
	out, err := ParseData(r)
	if err != nil {
		t.Skip("invalid number")
	}
	roundTrip := ToData(out)
	rtr := bytes.NewReader(roundTrip)
	out2, err := ParseData(rtr)
	if diff := cmp.Diff(out, out2); diff != "" {
		t.Error(diff)
	}
})
```

**Fuzzing consume muchos recursos**

Para ejecutar el fuzzer hacemos:

```ps
go test -fuzz=FuzzParseData
```

donde `FuzzParseData` es el nombre de nuestro caso de Fuzzer. Al ejecutar se creará en la carpeta `testdata\fuzz\[xxxxxx]` un juego de datos - archivo - por cada combinación que el Fuzzer ha detectado KO. `[xxxxxx]` es el nombre del caso de Fuzzer, ``FuzzParseData` en nuestro ejemplo. 

Si ahora hicieramos 

```ps
 go test -run=FuzzParseData/
```

se ejecutarán, a modo de regresión, los casos que el fuzzer detecto como erróneos (usa el dato que se guardo en `testdata\fuzz\`).

## Benchmarks

Los benchmarks se definen en métodos con la siguiente firma: `func BenchmarkXxxxxx(b *testing.B) {`. La estructura típicamente tiene un loop en el que iteramos hasta `b.N`. El objeto del loop es ponderar los resultados en varios experimentos:

```go
func BenchmarkFileLen1(b *testing.B) {
	// loop principal de benchmark
	for i := 0; i < b.N; i++ {
		result, err := FileLen("testdata/data.txt", 1)
		if err != nil {
			b.Fatal(err)
		}
		blackhole = result // usamos el resultado para evitar optimizaciones (si hacemos una llamada a una funcion, por ejemplo FileLen, y no usamos su resultado, el compilador puede optimizar y eliminar la llamada)
	}
}
```

también podemos definir los benchmark de forma programática. Por ejemplo aquí probamos con diferentes valores de longitud de palabra:

```go
func BenchmarkFileLen(b *testing.B) {
	// vamos a hacer el becnkmark para varios tamaños de palabra
	for _, v := range []int{1, 10, 100, 1000, 10000, 100000} {
		// lanzamos un benchmark para cada configuración
		b.Run(fmt.Sprintf("FileLen-%d", v), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				result, err := FileLen("testdata/data.txt", v)
				if err != nil {
					b.Fatal(err)
				}
				blackhole = result
			}
		})
	}
}
```

```ps
go test -bench=BenchmarkFileLen -benchmem 

goos: windows
goarch: amd64
pkg: github.com/learning-go-book-2e/ch15/bench
cpu: 13th Gen Intel(R) Core(TM) i7-1365U
BenchmarkFileLen1-12                   6         206920217 ns/op           65906 B/op      65208 allocs/op
BenchmarkFileLen/FileLen-1-12                  5         204911320 ns/op           65905 B/op      65208 allocs/op
BenchmarkFileLen/FileLen-10-12                54          19837691 ns/op          105048 B/op       6525 allocs/op
BenchmarkFileLen/FileLen-100-12              525           2179823 ns/op           73944 B/op        657 allocs/op
BenchmarkFileLen/FileLen-1000-12            2668            451986 ns/op           69304 B/op         70 allocs/op
BenchmarkFileLen/FileLen-10000-12           8931            141645 ns/op           82616 B/op         11 allocs/op
BenchmarkFileLen/FileLen-100000-12         10000            140110 ns/op          213688 B/op          5 allocs/op
PASS
ok      github.com/learning-go-book-2e/ch15/bench       15.310s
```

podemos ver los datos para cada valor del parámetro. Para ejecutar los benchmarks haremos:

```ps
go test -bench=[benchmark a ejecutar]
```

podemos ejecutar todos los benchmarks:

```ps
go test -bench=[benchmark a ejecutar] -benchmem
```

adicionalmente podemos incluir en el benchmark el uso de memoria:

```ps
go test -bench=.
```