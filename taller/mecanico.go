package taller

import (
	"fmt"
	"time"
)

type Mecanico struct {
	ID           string
	Especialidad TipoIncidencia
	Ocupado      bool
	ChanTrabajo  chan *Coche
	Trabajando   bool
	taller       *Taller
}

func NuevoMecanico(id string, especialidad TipoIncidencia) *Mecanico {
	return &Mecanico{
		ID:           id,
		Especialidad: especialidad,
		ChanTrabajo:  make(chan *Coche, 1),
		Trabajando:   false,
	}
}

func (m *Mecanico) Iniciar(taller *Taller) {
    m.taller = taller
    go func() {
        for coche := range m.ChanTrabajo {
            m.Trabajando = true
            m.Ocupado = true
            
            // FORZAR impresión inmediata
            fmt.Printf("🔧 %s → %s (acumulado: %v)\n", m.ID, coche, coche.TiempoAtendido)
            
            tiempoAtencion := coche.TiempoAtencion()
            time.Sleep(tiempoAtencion)
            
            coche.TiempoAtendido += tiempoAtencion
            
            if m.taller != nil {
                m.taller.RegistrarTiempoAtencion(tiempoAtencion)
            }
            
            // FORZAR impresión inmediata - SIN buffer
            fmt.Printf("✅ %s completó %s (sesión: %v, total: %v)\n", 
                m.ID, coche, tiempoAtencion, coche.TiempoAtendido)
            
            m.Trabajando = false
            m.Ocupado = false
        }
    }()
}

func (m *Mecanico) Detener() {
	close(m.ChanTrabajo)
}