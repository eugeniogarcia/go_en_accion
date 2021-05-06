// Package pool manages a user defined set of resources.
package pool

import (
	"errors"
	"io"
	"log"
	"sync"
)

// Pool manages a set of resources that can be shared safely by
// multiple goroutines. The resource being managed must implement
// the io.Closer interface.
type Pool struct {
	//Para sincronizar el acceso
	m sync.Mutex
	//El pool será un channel de interfaces io.Closer
	resources chan io.Closer
	//La factoría que creará el recurso cuando no haya uno disponible en el pool
	factory func() (io.Closer, error)
	//Indica si el pool esta o no operativo
	closed bool
}

// ErrPoolClosed is returned when an Acquire returns on a
// closed pool.
//Error comun que los métodos del pool pueden arrojar
var ErrPoolClosed = errors.New("Pool has been closed.")

// New creates a pool that manages resources. A pool requires a
// function that can allocate a new resource and the size of
// the pool.
//Funcion para construir un Pool - devuelve un puntero
func New(fn func() (io.Closer, error), size uint) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("Size value too small.")
	}

	//Construye el pool
	return &Pool{
		factory: fn,
		//Es un buffered channel con el tamaño de nuestro pool
		resources: make(chan io.Closer, size),
	}, nil
}

// Acquire retrieves a resource	from the pool.
// Este método le pide al pool un recuros
func (p *Pool) Acquire() (io.Closer, error) {
	//Si hay un recurso disponible en el buffered channel, lo recupera. Sino crea uno. No usamos el mutex para sincronizar, el propio channel actua de mecanismo de sincronización
	select {
	// Check for a free resource.
	case r, ok := <-p.resources:
		log.Println("Acquire:", "Shared Resource")
		if !ok {
			return nil, ErrPoolClosed
		}
		return r, nil

	// Provide a new resource since there are none available.
	default:
		log.Println("Acquire:", "New Resource")
		//Crea un nuevo recuro, que no viene del channel
		return p.factory()
	}
}

// Release places a new resource onto the pool.
// Método para devolver un recurso al pool cuando ya no se necesita
func (p *Pool) Release(r io.Closer) {
	// Secure this operation with the Close operation.

	//El metodo esta sincronizado con el metodo Close
	p.m.Lock()
	defer p.m.Unlock()

	// If the pool is closed, discard the resource.
	//Si el pool ya se cerro - necesitamos sincronizar este metodo con el método close -, cierra el recurso
	if p.closed {
		r.Close()
		return
	}

	//Guarda el recurso en el buffered channel para que otro lo pueda utilizar más tarde
	select {
	// Attempt to place the new resource on the queue.
	case p.resources <- r:
		log.Println("Release:", "In Queue")

	// If the queue is already at cap we close the resource.
	//Si ya tenemos en el pool el máximo de recursos, lo cerramos - no se reutilizara más
	default:
		log.Println("Release:", "Closing")
		r.Close()
	}
}

// Close will shutdown the pool and close all existing resources.
func (p *Pool) Close() {

	// Secure this operation with the Release operation.
	//El metodo esta sincronizado con el metodo Close
	p.m.Lock()
	defer p.m.Unlock()

	// If the pool is already close, don't do anything.
	if p.closed {
		return
	}

	// Set the pool as closed.
	p.closed = true

	// Close the channel before we drain the channel of its
	// resources. If we don't do this, we will have a deadlock.
	//Cerramos el channel. Impedira que nadie pueda añadir más recursos al channell
	close(p.resources)

	// Close the resources
	//Y nosotros mismo sacamos los recursos que ya hubiera en el channel, y los vamos cerrando
	for r := range p.resources {
		r.Close()
	}
}
