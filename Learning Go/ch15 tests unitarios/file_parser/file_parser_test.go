package file_parser

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseData(t *testing.T) {
	// slice con los casos de prueba a utilizar
	data := []struct {
		name   string
		in     []byte
		out    []string
		errMsg string
	}{
		{
			name:   "simple",
			in:     []byte("3\nhello\ngoodbye\ngreetings\n"),  // slice de bytes de entrada
			out:    []string{"hello", "goodbye", "greetings"}, // salida esperada
			errMsg: "",                                        // mensaje de error esperado
		},
		{
			name:   "empty_error",
			in:     []byte(""), // entrada vacía
			out:    nil,
			errMsg: "empty",
		},
		{
			name:   "zero",
			in:     []byte("0\n"), // entrada con cero líneas
			out:    []string{},
			errMsg: "",
		},
		{
			name:   "number_error",
			in:     []byte("asdf\nhello\ngoodbye\ngreetings\n"), // entrada en la que no se indica el número de líneas
			out:    nil,
			errMsg: `strconv.Atoi: parsing "asdf": invalid syntax`,
		},
		{
			name:   "line_count_error",
			in:     []byte("4\nhello\ngoodbye\ngreetings\n"), // entrada con menos líneas de las indicadas
			out:    nil,
			errMsg: "too few lines",
		},
	}
	// iterar sobre los casos de prueba
	for _, d := range data {
		// lanzamos el caso de prueba
		t.Run(d.name, func(t *testing.T) {
			// lector de la entrada
			r := bytes.NewReader(d.in)
			// parseo de los datos
			out, err := ParseData(r)
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			// compara la respuesta con la esperada, y muestra las diferencias, si las hubiera
			if diff := cmp.Diff(d.out, out); diff != "" {
				t.Error(diff)
			}

			// compara la respuesta con la esperada, y muestra las diferencias, si las hubiera
			if diff := cmp.Diff(d.errMsg, errMsg); diff != "" {
				t.Error(diff)
			}

			//if err == nil {
			//	roundTrip := ToData(out)
			//	if diff := cmp.Diff(d.in, roundTrip); diff != "" {
			//		t.Error(diff)
			//	}
			//}
		})
	}
}

func FuzzParseData(f *testing.F) {
	// crea un slice de slices de bytes como see de datos
	testcases := [][]byte{
		[]byte("3\nhello\ngoodbye\ngreetings\n"),
		[]byte("0\n"),
	}
	// crea el seed de datos
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}

	// lanzamos los tests utilizando el juego de datos fuzzeado
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
}
