- Los `if`s no tienen nada especial

- Los `for` podemo usarlos de cuatro formas:
    - completa, tipo C. 

    ```go
    func completeFor() {
        for i := 0; i < 10; i++ {
            fmt.Println(i)
        }
    }
    ```
    
    - solo la condicion (podríamos haber hecho también `for ;i<100; {}`>):

    ```go
    func conditionOnlyFor() {
        i := 1
        for i < 100 {
            fmt.Println(i)
            i = i * 2
        }
    }
    ```

    - infinita. Podemos usar `break` para salir

    ```go
    for {
            fmt.Println("Hello")
        }
    ```

    -  tipo foreach:

    ```go
    evenVals := []int{2, 4, 6, 8, 10, 12}
    for i, v := range evenVals {
        fmt.Println(i, v)
    }
    ```

Podemos usar `break` y `continue`. Por ejemplo para hacer un `do while`:

```go
for {
    // things to do in the loop
    if !CONDITION {
        break
    }
}
```

Podemos iterar en mapas, slices, arrays, strings:

```go
func main() {
	m := map[string]int{
		"a": 1,
		"c": 3,
		"b": 2,
	}

	for i := 0; i < 3; i++ {
		fmt.Println("Loop", i)
		for k, v := range m {
			fmt.Println(k, v)
		}
	}
}
```

```go
func main() {
	samples := []string{"hello", "apple_π!"}
	for _, sample := range samples {
		fmt.Println(sample)
		for i, r := range sample {
			fmt.Println(i, r, string(r))
		}
		fmt.Println()
	}

}
```

- Etiquetas como elemento de control. En `good_goto` y en `for_label` podemos ver como usar etiquetas

- `Switch`. Similar al switch de C, aunque no es necesario añadir `break` para evitar el fallback al siguiente case. Podemos usar _blank switchs_ (no se indica una variable a evaluar en el switch, y en cada case se incluye una expresión buleana), y el switch normal (en cada case se incluye un valor de la variable). En un case se pueden poner varios valores - posibles - de la variable (de ahí que no se necesite el `break` que se usaba o no en C). En `blank_switch` tenemos ejemplos de uso de _blank switchs_

