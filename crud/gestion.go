package crud

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Estructuras
type Cliente struct {
	ID        string
	Nombre    string
	Telefono  string
	Email     string
	Vehiculos []string
}

type Vehiculo struct {
	Matricula    string
	Marca        string
	Modelo       string
	FechaEntrada string
	FechaSalida  string
	IncidenciaID string
}

type Incidencia struct {
	ID          string
	VehiculoID  string   
	Mecanicos   []string
	Tipo        string
	Prioridad   string
	Descripcion string
	Estado      string
}

type Mecanico struct {
	ID           string
	Nombre       string
	Especialidad string
	AniosExp     int
	Activo       bool
}

type Plaza struct {
	Numero     int
	Ocupada    bool
	Matricula  string
	MecanicoID string
}

// Variables globales
var (
	Clientes    = make(map[string]Cliente)
	Vehiculos   = make(map[string]Vehiculo)
	Incidencias = make(map[string]Incidencia)
	Mecanicos   = make(map[string]Mecanico)
	Plazas      []Plaza
	ProxPlaza   = 1
)

func MenuPrincipal() {
	for {
		clearScreen()
		fmt.Println("=== GESTIÓN MANUAL ===")
		fmt.Println("1. Gestión de Clientes")
		fmt.Println("2. Gestión de Vehículos") 
		fmt.Println("3. Gestión de Incidencias")
		fmt.Println("4. Gestión de Mecánicos")
		fmt.Println("5. Asignar vehículo a plaza")
		fmt.Println("0. Volver al Menú Principal")
		fmt.Print("\nSeleccione opción: ")

		var opcion int
		fmt.Scan(&opcion)

		switch opcion {
		case 1:
			MenuClientes()
		case 2:
			MenuVehiculos()
		case 3:
			MenuIncidencias()
		case 4:
			MenuMecanicos()
		case 5:
			AsignarVehiculoAPlaza()
		case 0:
			return
		default:
			fmt.Println("Opción no válida")
			pausa()
		}
	}
}

func LeerLinea() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
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

func pausa() {
	fmt.Print("Presione Enter para continuar...")
	fmt.Scanln()
}

func MostrarEstadoTaller() {
	clearScreen()
	fmt.Println("=== ESTADO ACTUAL DEL TALLER ===")
	fmt.Printf("Clientes registrados: %d\n", len(Clientes))
	fmt.Printf("Vehículos registrados: %d\n", len(Vehiculos))
	fmt.Printf("Incidencias registradas: %d\n", len(Incidencias))
	fmt.Printf("Mecánicos registrados: %d\n", len(Mecanicos))
	fmt.Printf("Plazas creadas: %d\n", len(Plazas))
	
	fmt.Println("\nPlazas:")
	for _, p := range Plazas {
		estado := "Libre"
		if p.Ocupada {
			estado = "Ocupada - " + p.Matricula
		}
		fmt.Printf("Plaza %d (Mecánico: %s): %s\n", p.Numero, p.MecanicoID, estado)
	}
	
	pausa()
}

func AsignarVehiculoAPlaza() {
	fmt.Print("Matrícula del vehículo: ")
	var matricula string
	fmt.Scan(&matricula)

	_, existe := Vehiculos[matricula]
	if !existe {
		fmt.Println("Vehículo no encontrado")
		pausa()
		return
	}

	// Buscar plaza libre
	for i := range Plazas {
		if !Plazas[i].Ocupada {
			Plazas[i].Ocupada = true
			Plazas[i].Matricula = matricula
			fmt.Printf("Vehículo asignado a la plaza %d (Mecánico: %s)\n", 
				Plazas[i].Numero, Plazas[i].MecanicoID)
			pausa()
			return
		}
	}

	fmt.Println("No hay plazas disponibles")
	pausa()
}

// GetMecanicos devuelve todos los mecánicos
func GetMecanicos() ([]Mecanico, error) {
    var mecanicos []Mecanico
    for _, m := range Mecanicos {
        mecanicos = append(mecanicos, m)
    }
    return mecanicos, nil
}

// GetVehiculos devuelve todos los vehículos
func GetVehiculos() ([]Vehiculo, error) {
    var vehiculos []Vehiculo
    for _, v := range Vehiculos {
        vehiculos = append(vehiculos, v)
    }
    return vehiculos, nil
}

// GetIncidencias devuelve todas las incidencias
func GetIncidencias() ([]Incidencia, error) {
    var incidencias []Incidencia
    for _, i := range Incidencias {
        incidencias = append(incidencias, i)
    }
    return incidencias, nil
}