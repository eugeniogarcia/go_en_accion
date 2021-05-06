// Example is provided with help by Gabriel Aszalos.
// Package runner manages the running and lifetime of a process.
package runner

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

// Runner runs a set of tasks within a given timeout and can be
// shut down on an operating system interrupt.
//Cuando queramos ejecutar varias tareas secuencialmente, las especificamos en un slice de funciones. Aunque no expuestas a los clientes, para gestionar las tareas usaremos dos canales. complete, que sera el mecanismo por el que cada tarea de feedback cuando es procesada. El otro canal, interrupt, sirve para recivir mensajes desde el SSOO. Finalmente, establecemos una ventana de tiempo para ejecutar todas las tareas. El canal timeout se usara para comunicar que se ha superado este tiempo
type Runner struct {
	// interrupt channel reports a signal from the
	// operating system.
	interrupt chan os.Signal

	// complete channel reports that processing is done.
	complete chan error

	// timeout reports that time has run out.
	// Canal solo de recepcion
	timeout <-chan time.Time

	// tasks holds a set of functions that are executed
	// synchronously in index order.
	tasks []func(int)
}

// Errores genericos
// ErrTimeout is returned when a value is received on the timeout channel.
var ErrTimeout = errors.New("received timeout")

// ErrInterrupt is returned when an event from the OS is received.
var ErrInterrupt = errors.New("received interrupt")

// Funcion que Construye un runner
// New returns a new ready-to-use Runner.
func New(d time.Duration) *Runner {
	return &Runner{
		//Es un buffered channel de tamaño 1. Esto hace que quien publique una interrupción, no tenga que esperar a que se consuma, a no ser que ya hubiera una interrupción pendiente de procesar
		interrupt: make(chan os.Signal, 1),
		complete:  make(chan error),
		timeout:   time.After(d),
	}
}

// Add attaches tasks to the Runner. A task is a function that
// takes an int ID.
// Metodo para añadir tareas...
func (r *Runner) Add(tasks ...func(int)) {
	// Añade las tareas al slice
	r.tasks = append(r.tasks, tasks...)
}

// Start runs all tasks and monitors channel events.
// Inicia la ejecución de las tareas
func (r *Runner) Start() error {
	// We want to receive all interrupt based signals.
	//Lo que hacemos es suscribirnos a la recepción de interrupciones. Es non-blocking, y hace que las señales del SSOO se publiquen en el canal que indicamos. En este caso solo queremos que se publiquen las os.Interrupt, es decir, cuando se esta matando el proceso desde el OOSS
	signal.Notify(r.interrupt, os.Interrupt)

	// Run the different tasks on a different goroutine.
	//Lanza la go-rutina. La go-rutina se bloquea hasta que run() termina, y publica el resultado en el canal
	go func() {
		r.complete <- r.run()
	}()

	//Si hay datos en el canal los devuelbe y termina
	//Si hay datos en el canal timeout, termina
	//Mientras no haya nada de lo anterior, se bloquea a la espera
	select {
	// Signaled when processing is done.
	case err := <-r.complete:
		return err

	// Signaled when we run out of time.
	case <-r.timeout:
		return ErrTimeout
	}
}

// run executes each registered task.
// Implementa la lógica para ejecutar las tareas
func (r *Runner) run() error {
	//La tareas se ejecutan en serie...
	for id, task := range r.tasks {
		// Check for an interrupt signal from the OS.
		// primero comprobamos si hay que matar la ejecución
		if r.gotInterrupt() {
			return ErrInterrupt
		}

		// Execute the registered task.
		//la ejecuta
		task(id)
	}

	return nil
}

// gotInterrupt verifies if the interrupt signal has been issued.
func (r *Runner) gotInterrupt() bool {
	select {
	// Signaled when an interrupt event is sent.
	case <-r.interrupt:
		// Stop receiving any further signals.
		// Ya no recibimos más notificaciones
		signal.Stop(r.interrupt)
		return true

	// Continue running as normal.
	default:
		return false
	}
}
