package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// espera una interface io.Reader y devuelve un mapa con la cuenta de letras
func cuentaLetras(r io.Reader) (map[string]int, error) {
	buf := make([]byte, 2048) //crea un slice
	out := map[string]int{}   //crea el mapa de salida
	for {
		n, err := r.Read(buf) //leemos en el slice
		//tratamos los bytes leidos...
		for _, b := range buf[:n] {
			//...actualiza el mapa
			if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') {
				out[string(b)]++
			}
		}
		//hasta terminar de leer
		if err == io.EOF {
			return out, nil
		}
		if err != nil {
			return nil, err
		}
	}
}

func cuentaLetrasSimple() error {
	s := "The quick brown fox jumped over the lazy dog"
	//obtenemos un reader para el string
	sr := strings.NewReader(s)
	counts, err := cuentaLetras(sr)
	if err != nil {
		return err
	}
	fmt.Println(counts)
	return nil
}

func buildGZipReader(fileName string) (*gzip.Reader, func(), error) {
	r, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, nil, err
	}
	return gr, func() {
		gr.Close()
		r.Close()
	}, nil
}

func cuentaLetrasGzip() error {
	//obtenemos un reader para el gzip
	r, closer, err := buildGZipReader("my_data.txt.gz")
	if err != nil {
		return err
	}
	defer closer()

	counts, err := cuentaLetras(r)
	if err != nil {
		return err
	}
	fmt.Println(counts)
	return nil
}

func main() {
	err := cuentaLetrasSimple()
	if err != nil {
		slog.Error("error with simpleCountLetters", "msg", err)
	}

	err = cuentaLetrasGzip()
	if err != nil {
		slog.Error("error with gzipCountLetters", "msg", err)
	}
}
