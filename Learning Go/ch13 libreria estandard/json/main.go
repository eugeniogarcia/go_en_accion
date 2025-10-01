package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	err := ProcessPerson()
	if err != nil {
		slog.Error("error in processPerson", "msg", err)
	}
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func ProcessPerson() error {
	toFile := Person{
		Name: "Fred",
		Age:  40,
	}

	// Crea un archivo temporal
	tmpFile, err := os.CreateTemp(os.TempDir(), "ejemplo-")
	if err != nil {
		return err
	}
	// borrar el archivo al final
	defer os.Remove(tmpFile.Name())
	// Escribir la estructura en el archivo como JSON
	err = json.NewEncoder(tmpFile).Encode(toFile)
	if err != nil {
		return err
	}
	//ceramos el archivo
	err = tmpFile.Close()
	if err != nil {
		return err
	}

	// Leemos el archivo
	tmpFile2, err := os.Open(tmpFile.Name())
	if err != nil {
		return err
	}
	var fromFile Person
	// Decodificamos el JSON del archivo a una estructura
	err = json.NewDecoder(tmpFile2).Decode(&fromFile)
	if err != nil {
		return err
	}
	err = tmpFile2.Close()
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", fromFile)
	return nil
}
