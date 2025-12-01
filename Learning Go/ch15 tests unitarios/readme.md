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
