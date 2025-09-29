## gorutinas y canales

Creamos una gorutina cuando llamamos a una función cualquiera anteponiendo el comando `go`. Cuando llamamos a una gorutina podemos pasar cualquier argumento a la función, pero no se recupera nada de su respuesta.

Típicamente usamos _channels_ para interactuar con una corutina. Un _channel_ es un tipo por referencia (como un puntero) que se declara así

```go
[nombre] chan [tipo]
```

aqui declaramos un _channel_ llamado _nombre_ de tipo _tipo_.

Para leer o escribir en un channel:

```go

var ch chan int
var  val int

[...]

//leer
x:=<-ch

//escribir
val->ch
```

para pasar un channel como argumento de función podemos declarar la función de tres formas:

```go
//toma un channel al que se puede leer o escribir
func f1(ch chan int){}
//toma un channel del que se puede leer
func f2(ch <-chan int){}
//toma un channel al que se puede escribir
func f3(ch chan<- int){}
```

cuando se lee de un channel la llamada es bloqueante hasta que se recupera un valor. Cuando se escribe un channel la llamada es bloqueante hasta que se consume el valor - por el channel. Uso la palabra consume para introducir el concepto de buffer. Cuando escribimos en un channel, el valor se debe registrar en un buffer que tenga espacio disponible, o sino hay un buffer, se tiene que consumir del channel (esto es, cuando no hay un buffer el proceso que escribe y el que lee del channel se sincronizan.). Si tenemos un buffer en el channel, pero el buffer esta lleno, la llamada se bloquea hasta que se libere un espacion en el buffer (que pasará cuando alguien lee del channel). Otro tanto con la lectura. Si leemos de un channel la llamada se bloquea hasta que no recuperemos un dato. Si el channel tiene un buffer se leera del buffer (que pasa a tener un elemento menos).

Para crear un channel usamos `make`, de la misma forma que lo usamos con los slices:

```go
//crea un canal de int sin buffer
ch1:=make(chan int)
//crea un canal de int con un buffer de capacidad 10 (el largo sera 0)
ch2:=make(chan int, 10)
```

podemos usar `len()` y `cap()` para recuperar el tamaño y la capacidad del buffer

Cerramos un canal con `close([canal])`. Lo típico es cerrar el canal desde la gorutina que escribe en el canal. Si queremos escribir o cerrar un canal que ya esta cerrado se genera un _panic_. A la hora de leer un canal cerrado la respuesta no falla y no es bloqueante. Para saber si hemos leido un dato del canal o simplemente este estaba cerrado:

```go
val,ok:=<-ch
```

si ok es `true` se ha leido algo. Si es `false` el canal `ch` estaba cerrado. Lo leido estaría en `val`. El __garbage collector__ no detecta cuando no se va a escribir más en un canal. Por lo tanto, es necesario cerrar explicitamente un canal cuando no se vaya a escribir más.

Tenemos un ejemplo de todo esto en `goroutine`. Aquí se muestra un patrón típico en el que las gorutinas se implementan como _closures_.

En `deadlock` tenemos un ejemplo en el que se muestra como _es fácil_ sino tenemos cuidado provocar un deadlock.

## select

Podemos usar un select para procesar _channels_ con una sintaxis que recuerda _blank switchs_, aunque es diferente. 

```go
select {
    case ch2 <- v: //si hay alguien leyendo sería true
    case v2 = <-ch1: //si hay alguien escribiendo sería true
}
```

Los case se evaluan todos ellos de forma simultanea (no en cascada como en el switch), y en caso de que más de uno satisfaga la condicion, se elige uno de ellos al azar.

Cuando ninguna de las condiciones del select se cumplen, el código queda bloqueado.

Usar un default en un select no es buena idea, porque siempre se ejecutaría. Si tuvieramos incluido el select en un loop, el loop estarían siempre corriendo - en lugar de bloquear y dejar que la cpu se utilice para otras cosas.

```go
for {
    select {
    case <-done: //cuando alguien escriba en el canal done saldremos
        return
    case v := <-ch: //cuando alguien escriba en el canal ch imprimimos el valor
        fmt.Println(v)
    }
}
```

esto se suele llamar for-select loop.

## Cancelar

Este es un ejemplo de __gorutine leak__. Imaginemonos que tenemos que salir del bucle de main antes de iterar 10 veces:

```go
func countTo(max int) <-chan int {
    ch := make(chan int)
    go func() {
        for i := 0; i < max; i++ {
            ch <- i
        }
        close(ch)
    }()
    return ch
}

func main() {
    for i := range countTo(10) {
        if i > 5 {
            break //salimos del loop antes de consumir todos los datos del canal. Esto hace que la gorutina no termine nunca: gorutine leak
        }
        fmt.Println(i)
    }
}
```

Una forma de evitar el _gorutine leak_ es cancelar la ejecución. En el ejemplo `context_cancel` tenemos una demostración del uso de contextos de ejecución.

```go
ctx, cancel := context.WithCancel(context.Background()) //1. creamos un contexto cancelable. Devuelve el contexto y una función para cancelarlo
```

```go
select {
case <-ctx.Done(): //2. este método devuelve un canal que se cierra cuando el contexto se cancela
    return
case ch <- i:
}
```

```go
cancel() //3. llamamos a la función para cancelar el contexto
```