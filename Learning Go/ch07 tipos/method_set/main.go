package main

import (
	"fmt"
	"time"
)

type Counter struct {
	total       int
	lastUpdated time.Time
}

func (c *Counter) Increment() {
	c.total++
	c.lastUpdated = time.Now()
}
func (c Counter) String() string {
	return fmt.Sprintf("total: %d, last updated: %v", c.total, c.lastUpdated)
}

// interface
type Incrementer interface {
	Increment()
}

func main() {
	// Dos variables que tienen como tipo sendos interfaces
	var myStringer fmt.Stringer
	var myIncrementer Incrementer

	// Una variable de tipo puntero
	pointerCounter := &Counter{}
	// Una variable de tipo struct
	valueCounter := Counter{}

	// Ambas variables implementan la interfaz fmt.Stringer
	myStringer = pointerCounter // ok
	myStringer = valueCounter   // ok

	// Sólo la variable de tipo puntero implementa la interfaz Incrementer (porque el método Increment tiene un receiver puntero). Si tuviera un método Increment con receiver valor, entonces ambas variables implementarían la interfaz Incrementer
	myIncrementer = pointerCounter // ok
	myIncrementer = valueCounter   // el receiver es func (c *Counter) Increment() {

	fmt.Println(myStringer, myIncrementer)
}
