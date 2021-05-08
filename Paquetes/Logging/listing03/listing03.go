// This sample program demonstrates how to use the base log package.
package main

import (
	"log"
)

//Especifica la configuración. Prefijos y flags
func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func main() {
	// Println writes to the standard logger.
	log.Println("message")

	// Fatalln is Println() followed by a call to os.Exit(1).
	//Escribe el log y termina el programa
	log.Fatalln("fatal message")

	// Panicln is Println() followed by a call to panic().
	//Escribe el log y termina con un Panic
	log.Panicln("panic message")
}
