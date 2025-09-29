package person

import "fmt"

type Person struct {
	Name    string
	Age     int
	PetName string
}

var pets = map[string]Pet{
	"Fluffy": {"Fluffy", "Cat", "Bob"},
	"Rex":    {"Rex", "Dog", "Julia"},
}

func (p Person) Pet() Pet {
	return pets[p.PetName]
}

type Pet struct {
	Name      string
	Type      string
	OwnerName string
}

var owners = map[string]Person{
	"Bob":   {"Bob", 30, "Fluffy"},
	"Julia": {"Julia", 40, "Rex"},
}

func (p Pet) Owner() Person {
	return owners[p.OwnerName]
}

func init() {
	fmt.Println("paquete person inicializado")
}

func init() {
	fmt.Println("paquete person inicializado2")
}
