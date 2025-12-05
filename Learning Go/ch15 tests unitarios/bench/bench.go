package bench

import "os"

func FileLen(f string, bufsize int) (int, error) {
	// Abre un archivo para lectura
	file, err := os.Open(f)
	if err != nil {
		return 0, err
	}
	// Asegura que el archivo se cierre al finalizar la función
	defer file.Close()

	count := 0
	for {
		// creamos un slice para depositar lo que vamos leyendo
		buf := make([]byte, bufsize)
		// leemos del archivo
		num, err := file.Read(buf)
		count += num
		if err != nil {
			break
		}
	}
	return count, nil
}
