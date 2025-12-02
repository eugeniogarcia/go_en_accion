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