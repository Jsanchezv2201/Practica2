package taller

import (
	"sync"
)

type Cola struct {
	coches   []*Coche
	mutex    sync.Mutex
	cerrada  bool
	notify   chan struct{} // Para notificar cuando hay elementos nuevos
}

func NuevaCola() *Cola {
	return &Cola{
		coches:  make([]*Coche, 0),
		notify:  make(chan struct{}, 1),
		cerrada: false,
	}
}

func (c *Cola) AgregarCoche(coche *Coche) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	if c.cerrada {
		return
	}
	
	c.coches = append(c.coches, coche)
	
	// Notificar que hay un nuevo coche (non-blocking)
	select {
	case c.notify <- struct{}{}:
	default:
		// Ya hay una notificación pendiente
	}
}

func (c *Cola) ObtenerCoche() *Coche {
	for {
		c.mutex.Lock()
		
		// Si hay coches en la cola, devolver el primero
		if len(c.coches) > 0 {
			coche := c.coches[0]
			c.coches = c.coches[1:]
			c.mutex.Unlock()
			return coche
		}
		
		// Si la cola está cerrada y vacía, retornar nil
		if c.cerrada {
			c.mutex.Unlock()
			return nil
		}
		
		c.mutex.Unlock()
		
		// Esperar a que haya un coche o se cierre la cola
		<-c.notify
	}
}

func (c *Cola) Cerrar() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cerrada = true
	close(c.notify)
}

// Método auxiliar para obtener el tamaño actual (para debugging)
func (c *Cola) Tamaño() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.coches)
}