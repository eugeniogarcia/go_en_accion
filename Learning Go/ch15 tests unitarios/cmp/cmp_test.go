package cmp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreatePerson(t *testing.T) {
	esperado := Person{
		Name: "Dennis",
		Age:  37,
	}
	deseado := CreatePerson("Dennis", 37)

	// Comparamos las dos estructuras completas. El resultado es un texto en el que se indican las diferencias.
	if diff := cmp.Diff(esperado, deseado); diff != "" {
		t.Error(diff) // Si hay diferencias, fallamos el test mostrando las diferencias, y continuamos con otro test.
	}
}

func TestCreatePersonIgnoreDate(t *testing.T) {
	esperado := Person{
		Name: "Dennis",
		Age:  37,
	}
	deseado := CreatePerson("Dennis", 37)

	// Comparamos las dos estructuras, pero esta vez usando una función de comparación que no tiene en cuenta el campo DateAdded.
	comparer := cmp.Comparer(func(x, y Person) bool {
		return x.Name == y.Name && x.Age == y.Age
	})

	// el resultado de la comparación es un texto en el que se indican las diferencias.
	if diff := cmp.Diff(esperado, deseado, comparer); diff != "" {
		t.Error(diff) // Si hay diferencias, fallamos el test mostrando las diferencias, y continuamos con otro test.
	}

	// Además, comprobamos que el campo DateAdded se ha asignado correctamente.
	if deseado.DateAdded.IsZero() {
		t.Error("DateAdded wasn't assigned")
	}
}
