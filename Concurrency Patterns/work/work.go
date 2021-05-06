// Package work manages a pool of goroutines to perform work.
package work

import "sync"

// Worker must be implemented by types that want to use
// the work pool.
//Usamos un interface para definir que se puede ejecutar
type Worker interface {
	Task()
}

// Pool provides a pool of goroutines that can execute any Worker
// tasks that are submitted.
//Por un lado tenemos el canal por el que iran llegando a las go-rutinas el trabajo a realizar, y por otro lado el mecanismo para coordinar que el programa termine cuando acaben las gorutinas
type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

// New creates a new work pool.
func New(maxGoroutines int) *Pool {
	//Se crea una instancia del Pool
	p := Pool{
		work: make(chan Worker),
	}

	//Lanzamos las gorutinas
	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			//Espera hasta que haya trabajo en el canal. Hasta que no se cierre el canal la gorutina seguirá trabajando
			for w := range p.work {
				w.Task()
			}
			p.wg.Done()
		}()
	}

	return &p
}

// Run submits work to the pool.
//Alimentamos el canal con trabajo
func (p *Pool) Run(w Worker) {
	p.work <- w
}

// Shutdown waits for all the goroutines to shutdown.
//Cerramos el mecanismo. Cerramos el canal, y esperamos a que terminen todas las gorutinas
func (p *Pool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}
