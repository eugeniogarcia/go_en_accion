- `varadic`. Vemos como pasar un número variable de argumentos

- `anon_func`. Vemos como usar funciones anónimas

- Closures:

```go
func main() {
    a := 20
    f := func() {
        fmt.Println(a)
        a := 30 //capturamos la variable a
        fmt.Println(a)
    }
    f() //usa la variable capturada
    fmt.Println(a)
}
```

- En `defer_db` se muestra el uso de `defer`

