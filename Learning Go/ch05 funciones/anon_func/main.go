package main

import "fmt"

func main() {
	f := func(j int) {
		fmt.Println("printing", j, "from inside of an anonymous function 1")
	}
	for i := 0; i < 5; i++ {
		f(i)

		func(j int) {
			fmt.Println("printing", j, "from inside of an anonymous function 2")
		}(i)
	}
}
