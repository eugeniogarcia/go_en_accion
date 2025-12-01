package cleanup

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

// creamos un archivo temporal y lo limpiamos después de la prueba
func createFile(t *testing.T) (_ string, err error) {
	// creamos el archivo temporal
	f, err := os.Create("tempFile")
	if err != nil {
		return "", err
	}
	// aseguramos que el archivo este cerrado al final
	defer func() {
		// devuelve como error, los errores que se hayan detectado en la ejecución de la función más el posible error detectado al cerrar el archivo
		err = errors.Join(err, f.Close())
	}()
	// Aqui podemos escribir el contenido que queramos en el archivo
	// ...
	// ...

	// indicamos que instrucciones se han de ejecutar al finalizar la prueba
	t.Cleanup(func() {
		fmt.Printf("Limpia el archivo que hemos creado\n")
		os.Remove(f.Name())
	})

	//devuelve el nombre del archivo creado
	return f.Name(), nil
}

// createFileWithCreateTemp is a helper function called from multiple tests
func createFileWithCreateTemp(tempDir string) (_ string, err error) {
	// creamos un archivo temporal en el directorio especificado
	f, err := os.CreateTemp(tempDir, "tempFile")
	if err != nil {
		return "", err
	}
	defer func() {
		// indicamos que instrucciones se han de ejecutar al finalizar la prueba
		err = errors.Join(err, f.Close())
	}()

	// Aqui podemos escribir el contenido que queramos en el archivo
	// ...
	// ...

	//devuelve el nombre del archivo creado
	return f.Name(), nil
}

func TestFileProcessing(t *testing.T) {
	fmt.Printf("Test iniciado 1\n")
	// creamos un archivo para la prueba
	fName, err := createFile(t)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(fName, "tempFile") {
		t.Error("unexpected name")
	}

	fmt.Printf("Test 1 finalizado\n")
}

func TestFileProcessingWithCreateTemp(t *testing.T) {
	fmt.Printf("Test iniciado 2\n")

	// creamos un directorio temporal para la prueba
	tempDir := t.TempDir()

	// creamos un archivo temporal en el directorio temporal
	fName, err := createFileWithCreateTemp(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(fName, "tempFile") {
		t.Error("unexpected name")
	}

	fmt.Printf("Test 2 finalizado\n")

}
