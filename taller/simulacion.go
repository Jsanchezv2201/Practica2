package taller

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"practica2/crud"
)

type Configuracion struct {
	UsarDatosExistentes bool
	NumCoches           int
	TiposCoches         []TipoIncidencia
	MecanicosIniciales  []struct {
		ID           string
		Especialidad TipoIncidencia
	}
}

func EjecutarSimulacion(config Configuracion) (*Estadisticas, time.Duration) {
	start := time.Now()
	rand.Seed(time.Now().UnixNano())
	
	t := NuevoTaller()
	
	if config.UsarDatosExistentes {
		fmt.Println("\n=== SIMULACIÓN CON DATOS ACTUALES ===")
		fmt.Println("📊 Preparando simulación con datos existentes...")
		
		// Cargar mecánicos desde CRUD
		mecanicos, err := crud.GetMecanicos()
		if err != nil {
			fmt.Printf("❌ Error cargando mecánicos: %v\n", err)
			return nil, time.Since(start)
		}
		
		// Convertir mecánicos CRUD a mecánicos del taller
		for _, m := range mecanicos {
			// Convertir especialidad string a TipoIncidencia
			var especialidad TipoIncidencia
			switch strings.ToLower(m.Especialidad) {
			case "mecánica":
				especialidad = Mecanica
			case "eléctrica":
				especialidad = Electrica
			case "carrocería":
				especialidad = Carroceria
			default:
				especialidad = Mecanica // Por defecto
			}
			
			mecanico := NuevoMecanico(m.ID, especialidad)
			t.AgregarMecanico(mecanico)
			fmt.Printf("   • Mecánico %s (%s)\n", m.ID, m.Especialidad)
		}
		
		// Cargar vehículos desde CRUD
		vehiculos, err := crud.GetVehiculos()
		if err != nil {
			fmt.Printf("❌ Error cargando vehículos: %v\n", err)
			return nil, time.Since(start)
		}
		
		// Cargar incidencias para determinar el tipo de cada vehículo
		incidencias, err := crud.GetIncidencias()
		if err != nil {
			fmt.Printf("❌ Error cargando incidencias: %v\n", err)
			return nil, time.Since(start)
		}
		
		// Crear mapa de incidencias por ID de vehículo
		tiposPorVehiculo := make(map[string]TipoIncidencia)
		for _, inc := range incidencias {
			if inc.VehiculoID != "" { // ← VERIFICAR QUE TENGA VehiculoID
				var tipo TipoIncidencia
				switch strings.ToLower(inc.Tipo) {
				case "mecánica":
					tipo = Mecanica
				case "eléctrica":
					tipo = Electrica
				case "carrocería":
					tipo = Carroceria
				default:
					tipo = Mecanica // Por defecto
				}
				tiposPorVehiculo[inc.VehiculoID] = tipo
			}
		}
		
		// Crear coches para la simulación
		cochesCreados := 0
		for _, v := range vehiculos {
			tipo, exists := tiposPorVehiculo[v.Matricula]
			if !exists {
				// Si no hay incidencia, usar tipo aleatorio
				tipos := []TipoIncidencia{Mecanica, Electrica, Carroceria}
				tipo = tipos[rand.Intn(len(tipos))]
			}
			
			coche := NuevoCoche(v.Matricula, tipo)
			t.LlegadaCoche(coche)  
			cochesCreados++
			time.Sleep(time.Duration(200 + rand.Intn(300)) * time.Millisecond)
		}
		
		fmt.Printf("\n✅ Configuración cargada: %d mecánico(s), %d vehículo(s)\n", 
			len(mecanicos), cochesCreados)
			
	} else {
		fmt.Println("\n=== SIMULACIÓN AUTOMÁTICA ===")
		fmt.Println("🎯 Usando configuración automática...")
		
		// Configuración automática por defecto
		if len(config.TiposCoches) == 0 {
			config.TiposCoches = []TipoIncidencia{Mecanica, Electrica, Carroceria}
		}
		
		if config.NumCoches == 0 {
			config.NumCoches = 8
		}
		
		// Agregar mecánicos base o los proporcionados
		if len(config.MecanicosIniciales) > 0 {
			for _, m := range config.MecanicosIniciales {
				mecanico := NuevoMecanico(m.ID, m.Especialidad)
				t.AgregarMecanico(mecanico)
				fmt.Printf("   • Mecánico %s (%s)\n", m.ID, m.Especialidad.String()) 
			}
		} else {
			// Mecánicos base por defecto
			mecanicosBase := []struct {
				ID           string
				Especialidad TipoIncidencia
			}{
				{"M1", Mecanica},
				{"E1", Electrica},
				{"C1", Carroceria},
			}
			
			for _, m := range mecanicosBase {
				mecanico := NuevoMecanico(m.ID, m.Especialidad)
				t.AgregarMecanico(mecanico)
				fmt.Printf("   • Mecánico %s (%s)\n", m.ID, m.Especialidad.String()) 
			}
		}
		
		// Generar coches
		distribucion := make(map[TipoIncidencia]int)
		for i := 0; i < config.NumCoches; i++ {
			tipo := config.TiposCoches[rand.Intn(len(config.TiposCoches))]
			coche := NuevoCoche(fmt.Sprintf("C%d", i+1), tipo)
			t.LlegadaCoche(coche)  
			distribucion[tipo]++
			time.Sleep(time.Duration(200 + rand.Intn(300)) * time.Millisecond)
		}
		
		fmt.Printf("\n📊 Distribución de coches: ")
		primero := true
		for tipo, count := range distribucion {
			if !primero {
				fmt.Print(", ")
			}
			fmt.Printf("%d %s", count, tipo.String()) 
			primero = false
		}
		fmt.Println()
	}
	
	fmt.Println("\n🚀 Iniciando simulación...")
	t.Iniciar()
	
	// CANAL para controlar la finalización
	done := make(chan bool, 1)
	
	// Goroutine principal de monitoreo
	go func() {
		timeout := time.After(45 * time.Second) // Timeout más largo
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				if t.Stats.CochesAtendidos >= t.Stats.CochesTotales {
					time.Sleep(10 * time.Second)
					t.Detener()
					
					time.Sleep(4 * time.Second)
					
					t.ObtenerYLimpiarBuffer()
					
					fmt.Printf("\n✅ Simulación completada: %d/%d coches\n", 
						t.Stats.CochesAtendidos, t.Stats.CochesTotales)
					done <- true
					return
				}
			case <-timeout:
				fmt.Printf("\n⏰ Timeout: Atendidos %d/%d coches\n", 
					t.Stats.CochesAtendidos, t.Stats.CochesTotales)
				t.Detener()
				time.Sleep(4 * time.Second)
				t.ObtenerYLimpiarBuffer() // Limpiar buffer en timeout
				done <- true
				return
			}
		}
	}()
	
	// Esperar a que termine la goroutine de monitoreo
	<-done
	
	// PAUSA FINAL ADICIONAL
	time.Sleep(2 * time.Second)
	return t.Stats, time.Since(start)
}

// Función auxiliar para crear configuraciones predefinidas
func CrearConfiguracionAutomatica(escenario int) Configuracion {
	switch escenario {
	case 1:
		// Configuración base
		return Configuracion{
			UsarDatosExistentes: false,
			NumCoches:           8,
			TiposCoches:         []TipoIncidencia{Mecanica, Electrica, Carroceria},
		}
	case 2:
		// Doble de coches
		return Configuracion{
			UsarDatosExistentes: false,
			NumCoches:           16,
			TiposCoches:         []TipoIncidencia{Mecanica, Electrica, Carroceria},
			MecanicosIniciales: []struct {
				ID           string
				Especialidad TipoIncidencia
			}{
				{"M1", Mecanica},
				{"E1", Electrica},
				{"C1", Carroceria},
			},
		}
	case 3:
		// 3 mecánicos mecánica / 1 eléctrica / 1 carrocería 
		return Configuracion{
			UsarDatosExistentes: false,
			NumCoches:           8,  
			TiposCoches:         []TipoIncidencia{Mecanica, Electrica, Carroceria},
			MecanicosIniciales: []struct {
				ID           string
				Especialidad TipoIncidencia
			}{
				{"M1", Mecanica},
				{"M2", Mecanica},
				{"M3", Mecanica},  
				{"E1", Electrica},
				{"C1", Carroceria},
			},
		}
	case 4:
		// Duplicar plantilla (6 mecánicos)
		return Configuracion{
			UsarDatosExistentes: false,
			NumCoches:           8,
			TiposCoches:         []TipoIncidencia{Mecanica, Electrica, Carroceria},
			MecanicosIniciales: []struct {
				ID           string
				Especialidad TipoIncidencia
			}{
				{"M1", Mecanica},
				{"M2", Mecanica},
				{"E1", Electrica},
				{"E2", Electrica},
				{"C1", Carroceria},
				{"C2", Carroceria},
			},
		}
	case 5:
		// Distribución extrema 1-3-3
		return Configuracion{
			UsarDatosExistentes: false,
			NumCoches:           8,
			TiposCoches:         []TipoIncidencia{Mecanica, Electrica, Carroceria},
			MecanicosIniciales: []struct {
				ID           string
				Especialidad TipoIncidencia
			}{
				{"M1", Mecanica},
				{"E1", Electrica},
				{"E2", Electrica},
				{"E3", Electrica},
				{"C1", Carroceria},
				{"C2", Carroceria},
				{"C3", Carroceria},
			},
		}
	default:
		return Configuracion{
			UsarDatosExistentes: false,
			NumCoches:           8,
			TiposCoches:         []TipoIncidencia{Mecanica, Electrica, Carroceria},
		}
	}
}