package taller

import (
	"fmt"
	"math/rand"
	"time"
)

type TipoIncidencia string

const (
	Mecanica   TipoIncidencia = "mecánica"
	Electrica  TipoIncidencia = "eléctrica"
	Carroceria TipoIncidencia = "carrocería"
)

type Coche struct {
	Matricula     string         
	ID            string
	TipoIncidencia TipoIncidencia
	TiempoAtendido time.Duration
	ChanTerminado chan bool
	TiempoLlegada time.Time
}

func NuevoCoche(matricula string, tipo TipoIncidencia) *Coche {
	return &Coche{
		Matricula:     matricula, 
		ID:            matricula, 
		TipoIncidencia: tipo,
		ChanTerminado: make(chan bool),
		TiempoLlegada: time.Now(),
	}
}

func (c *Coche) TiempoAtencion() time.Duration {
	switch c.TipoIncidencia {
	case Mecanica:
		return time.Duration(5+rand.Intn(3)) * time.Second // 5-7 segundos
	case Electrica:
		return time.Duration(7+rand.Intn(3)) * time.Second // 7-9 segundos
	case Carroceria:
		return time.Duration(11+rand.Intn(5)) * time.Second // 11-15 segundos
	default:
		return 5 * time.Second
	}
}

func (c *Coche) String() string {
	return fmt.Sprintf("Coche %s [%s]", c.Matricula, c.TipoIncidencia)
}

func (t TipoIncidencia) String() string {
	return string(t)
}