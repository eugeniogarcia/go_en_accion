package main

import (
	"fmt"

	"github.com/shopspring/decimal"
	"gz.com/ch10/paquetes/convert"
	persona "gz.com/ch10/paquetes/person"
)

func main() {
	bob := persona.Person{PetName: "Fluffy"}
	fmt.Println(bob.Pet())

	money := convert.Money{Value: decimal.NewFromFloat(10.0), Currency: "USD"}
	converted, err := convert.Convert(money, "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Converted: %v %s\n", converted.Value, converted.Currency)
	}

}
