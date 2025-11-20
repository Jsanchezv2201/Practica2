package taller

import (
	"fmt"
	"time"
)

// ESTRUCTURAS
type Estadisticas struct {
	CochesAtendidos      int
	CochesTotales        int
	TiempoTotal          time.Duration
	MecanicosContratados int
	TiemposAtencion      []time.Duration
	CochesPrioritarios   int
}

type Taller struct {
	Cola               *Cola
	Mecanicos          []*Mecanico
	ChanDetener        chan bool
	Stats              *Estadisticas
	running            bool
	mensajesBuffer     []string  
}

// CONSTRUCTOR
func NuevoTaller() *Taller {
	return &Taller{
		Cola:               NuevaCola(),
		ChanDetener:        make(chan bool),
		Stats: &Estadisticas{
			TiemposAtencion: make([]time.Duration, 0),
		},
		running: true,
	}
}

// MÉTODOS DEL TALLER
func (t *Taller) AgregarMecanico(mecanico *Mecanico) {
	mecanico.Iniciar(t)
	t.Mecanicos = append(t.Mecanicos, mecanico)
}

func (t *Taller) Iniciar() {
	go t.coordinator()
}

func (t *Taller) coordinator() {
	for t.running {
		coche := t.Cola.ObtenerCoche()
		if coche == nil {
			// Cola cerrada
			return
		}
		
		if coche.TiempoAtendido > 15*time.Second {
			fmt.Printf("🚨 PRIORIDAD: %s (acumulado: %v)\n", coche, coche.TiempoAtendido)
			t.Stats.CochesPrioritarios++
			t.atiendeCochePrioritario(coche)
		} else {
			t.atiendeCocheNormal(coche)
		}
	}
}


func (t *Taller) LlegadaCoche(coche *Coche) {
    if t.running {
        t.Stats.CochesTotales++
        t.Cola.AgregarCoche(coche)
        fmt.Printf("🚗 %s llegó al taller\n", coche)
    }
}


func (t *Taller) atiendeCochePrioritario(coche *Coche) {
	mecanico := t.buscarMecanicoLibreCualquierEspecialidad()
	
	if mecanico != nil {
		select {
		case mecanico.ChanTrabajo <- coche:
			t.Stats.CochesAtendidos++
		default:
			go func(c *Coche) {
				time.Sleep(100 * time.Millisecond)
				if t.running {
					t.Cola.AgregarCoche(c)
				}
			}(coche)
		}
	} else {
		t.BufferMensaje(fmt.Sprintf("📢 CONTRATANDO: nuevo mecánico para %s", coche))
		nuevoID := fmt.Sprintf("N%d", len(t.Mecanicos)+1)
		nuevoMecanico := NuevoMecanico(nuevoID, coche.TipoIncidencia)
		t.AgregarMecanico(nuevoMecanico)
		
		select {
		case nuevoMecanico.ChanTrabajo <- coche:
			t.Stats.CochesAtendidos++
			t.Stats.MecanicosContratados++
		default:
			go func(c *Coche) {
				time.Sleep(100 * time.Millisecond)
				if t.running {
					t.Cola.AgregarCoche(c)
				}
			}(coche)
		}
	}
}

func (t *Taller) atiendeCocheNormal(coche *Coche) {
	mecanico := t.buscarMecanicoLibre(coche.TipoIncidencia)
	if mecanico != nil {
		select {
		case mecanico.ChanTrabajo <- coche:
			t.Stats.CochesAtendidos++
		default:
			go func(c *Coche) {
				time.Sleep(100 * time.Millisecond)
				if t.running {
					t.Cola.AgregarCoche(c)
				}
			}(coche)
		}
	} else {
		go func(c *Coche) {
			waitTime := 1 * time.Second
			time.Sleep(waitTime)
			
			if t.running {
				c.TiempoAtendido += waitTime
				if c.TiempoAtendido % (5 * time.Second) == 0 {
					fmt.Printf("⏳ %s en cola (acumulado: %v)\n", c, c.TiempoAtendido)
				}
				t.Cola.AgregarCoche(c)
			}
		}(coche)
	}
}

func (t *Taller) buscarMecanicoLibre(tipo TipoIncidencia) *Mecanico {
	for _, m := range t.Mecanicos {
		if !m.Ocupado && m.Especialidad == tipo {
			return m
		}
	}
	return nil
}

func (t *Taller) buscarMecanicoLibreCualquierEspecialidad() *Mecanico {
	for _, m := range t.Mecanicos {
		if !m.Ocupado {
			return m
		}
	}
	return nil
}

func (t *Taller) Detener() {
	t.running = false
	
	// 1. Cerrar cola primero
	t.Cola.Cerrar()
	
	// 2. Luego cerrar canal de detener
	select {
	case <-t.ChanDetener:
		// Ya cerrado
	default:
		close(t.ChanDetener)
	}
	
	// 3. Detener mecánicos
	for _, m := range t.Mecanicos {
		m.Detener()
	}
	
	time.Sleep(1 * time.Second)
}

func (t *Taller) RegistrarTiempoAtencion(tiempo time.Duration) {
	t.Stats.TiemposAtencion = append(t.Stats.TiemposAtencion, tiempo)
}

// MÉTODOS PARA ESTADÍSTICAS
func (e *Estadisticas) TiempoPromedioAtencion() time.Duration {
	if len(e.TiemposAtencion) == 0 {
		return 0
	}
	var total time.Duration
	for _, t := range e.TiemposAtencion {
		total += t
	}
	return total / time.Duration(len(e.TiemposAtencion))
}

func (e *Estadisticas) Eficiencia() float64 {
	if e.CochesTotales == 0 {
		return 0
	}
	return float64(e.CochesAtendidos) / float64(e.CochesTotales) * 100
}

// Añadir mensaje al buffer
func (t *Taller) BufferMensaje(mensaje string) {
	if t.running {
		t.mensajesBuffer = append(t.mensajesBuffer, mensaje)
		fmt.Println(mensaje) 
	}
}

// Obtener y limpiar buffer
func (t *Taller) ObtenerYLimpiarBuffer() []string {
	buffer := make([]string, len(t.mensajesBuffer))
	copy(buffer, t.mensajesBuffer)
	t.mensajesBuffer = make([]string, 0)
	return buffer
}