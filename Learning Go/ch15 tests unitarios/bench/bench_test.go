package bench

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

var rng *rand.Rand

func TestMain(m *testing.M) {
	// creamos datos
	rng = rand.New(rand.NewSource(1))
	makeData()
	//ejecutamos los tests
	exitVal := m.Run()
	// borramos los datos
	os.Remove("testdata/data.txt")
	os.Exit(exitVal)
}

func makeData() {
	file, err := os.Create("testdata/data.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for i := 0; i < 10000; i++ {
		data := makeWord(rng.Intn(10) + 1)
		file.Write(data)
	}
}

func makeWord(l int) []byte {
	out := make([]byte, l+1)
	for i := 0; i < l; i++ {
		out[i] = 'a' + byte(rng.Intn(26))
	}
	out[l] = '\n'
	return out
}

func TestFileLen(t *testing.T) {
	result, err := FileLen("testdata/data.txt", 1)
	if err != nil {
		t.Fatal(err)
	}
	if result != 65204 {
		t.Error("Expected 65204, got", result)
	}
}

var blackhole int

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
