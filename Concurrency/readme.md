# Paralelismo y concurrencia

- listing01. Un solo procesador, go rutinas rapidas
- listing04. Un solo procesador, go rutinas lentas
- listing07. Más de un procesador - __ejecución en paralelo__. El programa es el mismo que teníamos en _listing01_

# Problema

- listing09. Programa con un defecto que demuestra que sucede cuando se accede a un recurso común de forma no sincronizada, provocando una _race condition_. Demuestra también el uso de la función _Gosched_. Podemos usar _go build -race ._ para comprobar que hay una race condition

## Atomic

- listing13. Soluciona el problema de _listing09_ usando una función __atomic.AddInt64(&counter, 1)__

## Load & Store

- listing15. Soliciona el problema de _listing09_ usando las funciones __atomic.StoreInt64(&shutdown, 1)__ y __atomic.LoadInt64(&shutdown)__. Estas funciones permiten escribir y leer de forma sincronizada. Cuando usemos esta función podemos tener la garantía de que nadie más esta escribiendo o leyendo al mismo tiempo en la variable.

## Mutex
 
- listing16. Declara una _critical section_ usando un __mutex__. Con _mutex.Lock()_ indicamos que empieza la seccion critica, y con _mutex.Unlock()_ que termina. La seccion crítica solo puede ejecutarla una go rutina. Notese que hay _runtime.Gosched()_ en la seccion crítica. En este caso el go scheduler volverá a asignar el slot de tiempo a la misma go rutina

## Canales

- listing20. Demuestra el uso de un __unbuffered channel__. Las dos gorutinas estan _enlazadas_ con el canal. Cuando una envia datos por el canal, se mantiene a la espera hasta que los datos son consumidos del canal. Cuando el canal se cierra, cualquier go rutina que estuviera esperando a recibir datos deja de esperar por ellos.
- listing22. Usa las mismas construcciones que en el ejemplo anterior. En este caso una go mism runtina recibe por el canal, se lanza aís misma, y escribe en el canal. Cuando escribe en el canal despierta al hijo que ha creado, que repite el patron hasta cuatro veces.
- listing24. Demuestra el uso de __buffered channels__. El emisor bloqueara si no hay hueco en el channel para escribir. El lector bloqueara si no hay nada en el channel que leer.
    - Tenemos un método _init()_ para iniciar la semilla del generador de números aleatorios
    - Se crea un buffered channel _tasks := make(chan string, taskLoad)_
    - Lanza varias go-rutinas. Se ejecutaran en el procesdor lógico

    ```go
    wg.Add(numberGoroutines)
	for gr := 1; gr <= numberGoroutines; gr++ {
		go worker(tasks, gr)
	}
    ```
    
    Al cerrar el canal, las go rutinas pueden aún recibir datos por él, pero no podrán ya escribir en él. En este ejemplo demostramos esta caracteristica; Se escriben 10 tareas en el canal, y tan pronto se han escrito se cierra. Vemos que las go rutinas tienen un loop en el que leen del canal, procesan, vuelven a escuchar. Mientras haya datos, la go rutina los procesara; __si no hay datos, y el canal esta cerrado__, la go rutina sale

