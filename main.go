package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"practica2/crud"
	"practica2/taller"
	"runtime"
	"strings"
	"time"
)

func main() {
	for {
		clearScreen()
		fmt.Println("=== TALLER MECÁNICO - PRÁCTICA 2 ===")
		fmt.Println("1. Gestión Manual (Clientes, Vehículos, Incidencias, Mecánicos)")
		fmt.Println("2. Ejecutar Simulación Automática")
		fmt.Println("3. Simulación con Datos Actuales")  
		fmt.Println("4. Estado Actual del Taller")
		fmt.Println("5. Ejecutar Tests")
		fmt.Println("0. Salir")
		fmt.Print("\nSeleccione opción: ")
		
		var opcion int
		fmt.Scan(&opcion)
		
		switch opcion {
		case 1:
			crud.MenuPrincipal()
		case 2:
			ejecutarSimulacionAutomatica()
		case 3:  
			ejecutarSimulacionConDatosActuales()
		case 4:
			crud.MostrarEstadoTaller()
		case 5:
			ejecutarTests()
		case 0:
			fmt.Println("¡Hasta luego!")
			return
		default:
			fmt.Println("Opción no válida")
			pausa()
		}
	}
}

func ejecutarSimulacionAutomatica() {
	clearScreen()
	fmt.Println("=== SIMULACIÓN AUTOMÁTICA ===")
	fmt.Println("Comparativas según requisitos del enunciado:")
	fmt.Println("1. Duplicar cantidad de coches")
	fmt.Println("2. Duplicar plantilla de mecánicos")
	fmt.Println("3. Distribuciones desbalanceadas de especialidades")
	fmt.Println("")
	
	escenarios := []struct {
		nombre      string
		descripcion string
		configNum   int
	}{
		{
			"CONFIGURACIÓN BASE (REFERENCIA)",
			"• OBJETIVO: Establecer línea base para comparativas\n• MECÁNICOS: 3 (uno de cada especialidad)\n• COCHES: 8 con distribución equilibrada\n• PROPÓSITO: Medir eficiencia del sistema estándar",
			1,
		},
		{
			"DUPLICAR CANTIDAD DE COCHES", 
			"• OBJETIVO: Test de carga - duplicar coches (16 vs 8)\n• MECÁNICOS: 4 (2 mecánica, 1 eléctrica, 1 carrocería)\n• COCHES: 16 (doble del escenario base)\n• PROPÓSITO: Ver escalabilidad con más demanda\n• RELACIÓN ENUNCIADO: 'cantidad máxima de coches se duplica'",
			2,
		},
		{
			"MECÁNICOS ESPECIALIZADOS",
			"• OBJETIVO: Test de especialización desbalanceada\n• MECÁNICOS: 5 (3 mecánica, 1 eléctrica, 1 carrocería)\n• COCHES: 15 con mayoría eléctrica/carrocería\n• PROPÓSITO: Eficiencia con especialización específica\n• RELACIÓN ENUNCIADO: '3 mecánicos mecánica por cada eléctrica/carrocería'",
			3,
		},
		{
			"DUPLICAR PLANTILLA (6 MECÁNICOS)",
			"• OBJETIVO: Test de recursos - duplicar mecánicos\n• MECÁNICOS: 6 (2 de cada especialidad)\n• COCHES: 12 con distribución equilibrada\n• PROPÓSITO: Medir mejora con más recursos humanos\n• RELACIÓN ENUNCIADO: 'duplicamos la plantilla de 3 a 6 mecánicos'",
			4,
		},
		{
			"DISTRIBUCIÓN EXTREMA 1-3-3", 
			"• OBJETIVO: Test de distribución muy desbalanceada\n• MECÁNICOS: 7 (1 mecánica, 3 eléctrica, 3 carrocería)\n• COCHES: 10 con distribución variada\n• PROPÓSITO: Eficiencia con especialización extrema\n• RELACIÓN ENUNCIADO: '1 mecánico mecánica por cada 3 eléctrica/3 carrocería'",
			5,
		},
	}

	for i, escenario := range escenarios {
		clearScreen()
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Printf("🎯 ESCENARIO %d: %s\n", i+1, escenario.nombre)
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println(escenario.descripcion)
		fmt.Println("\n⏳ Iniciando simulación en 3 segundos...")
		time.Sleep(3 * time.Second)

		// Ejecutar simulación
		config := taller.CrearConfiguracionAutomatica(escenario.configNum)
		stats, duracion := taller.EjecutarSimulacion(config)
		
		mostrarResultados(escenario.nombre, stats, duracion)
		
		// Si no es el último escenario, mostrar mensaje de continuación
		if i < len(escenarios)-1 {
			fmt.Println("\n" + strings.Repeat("-", 50))
			fmt.Println("🔄 Preparando siguiente escenario...")
			time.Sleep(2 * time.Second)
		}
	}

	// MENSAJE FINAL 
	clearScreen()
	fmt.Println("✅ SIMULACIÓN COMPLETADA - TODOS LOS ESCENARIOS PROBADOS")
	fmt.Println("\nResumen de escenarios ejecutados:")
	for i, escenario := range escenarios {
		fmt.Printf("  %d. %s\n", i+1, escenario.nombre)
	}
	fmt.Println("\nPresione Enter para volver al menú principal...")
	pausaForzada() 
}

func ejecutarTests() {
	clearScreen()
	fmt.Println("=== EJECUTANDO TESTS COMPARATIVOS ===")
	fmt.Println("Ejecuta: go test -v")
	fmt.Println("Esto probará 5 escenarios diferentes...")
	pausa()
}

func mostrarResultados(nombre string, stats *taller.Estadisticas, duracion time.Duration) {
	fmt.Println("\n⏳ Finalizando simulación y recopilando resultados...")
	time.Sleep(10 * time.Second) 
	
	fmt.Printf("\n📊 RESULTADOS - %s:\n", nombre)
	fmt.Printf("   Duración total:      %v\n", duracion)
	fmt.Printf("   Coches totales:      %d\n", stats.CochesTotales)
	fmt.Printf("   Coches atendidos:    %d\n", stats.CochesAtendidos)
	fmt.Printf("   Eficiencia:          %.1f%%\n", stats.Eficiencia())
	fmt.Printf("   Mecánicos extra:     %d\n", stats.MecanicosContratados)
	fmt.Printf("   Coches prioritarios: %d\n", stats.CochesPrioritarios)
	
	if len(stats.TiemposAtencion) > 0 {
		fmt.Printf("   Tiempo promedio:     %v\n", stats.TiempoPromedioAtencion())
	}
	
	fmt.Println("\n" + strings.Repeat("-", 50))
	pausaForzada()
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func pausaForzada() {
	fmt.Print("Presione Enter para continuar...")
	
	// Leer entrada sin buffer para forzar la espera
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}

func pausa() {
	fmt.Print("Presione Enter para continuar...")
	var discard string
	fmt.Scanln(&discard)
}

func ejecutarSimulacionConDatosActuales() {
	clearScreen()
	fmt.Println("=== SIMULACIÓN CON DATOS ACTUALES ===")
	
	// Verificar que hay datos suficientes
	if len(crud.Mecanicos) == 0 {
		fmt.Println("❌ ERROR: No hay mecánicos creados.")
		fmt.Println("   Crea al menos un mecánico de cada especialidad primero.")
		pausa()
		return
	}
	
	if len(crud.Vehiculos) == 0 {
		fmt.Println("❌ ERROR: No hay vehículos creados.")
		fmt.Println("   Crea algunos vehículos primero.")
		pausa()
		return
	}
	
	fmt.Printf("📊 Preparando simulación con:\n")
	fmt.Printf("   • %d mecánico(s)\n", len(crud.Mecanicos))
	fmt.Printf("   • %d vehículo(s)\n", len(crud.Vehiculos))
	fmt.Printf("   • %d incidencia(s)\n", len(crud.Incidencias))
	fmt.Println("\n⏳ Iniciando simulación en 3 segundos...")
	time.Sleep(3 * time.Second)
	
	// Configurar la simulación con datos del CRUD
	config := taller.Configuracion{
		UsarDatosExistentes: true,
	}
	
	stats, duracion := taller.EjecutarSimulacion(config)
	
	// Mostrar resultados
	fmt.Println("🎯 RESULTADOS DE LA SIMULACIÓN")
	fmt.Println("===============================")
	fmt.Printf("Duración total:        %v\n", duracion)
	fmt.Printf("Vehículos totales:     %d\n", stats.CochesTotales)
	fmt.Printf("Vehículos atendidos:   %d\n", stats.CochesAtendidos)
	fmt.Printf("Eficiencia:            %.1f%%\n", stats.Eficiencia())
	fmt.Printf("Mecánicos contratados: %d\n", stats.MecanicosContratados)
	fmt.Printf("Vehículos prioritarios: %d\n", stats.CochesPrioritarios)
	
	if stats.Eficiencia() == 100.0 {
		fmt.Println("\n✅ ¡Todos los vehículos fueron atendidos!")
	} else {
		fmt.Printf("\n⚠️  Algunos vehículos no pudieron ser atendidos\n")
	}
	
	pausa()
}