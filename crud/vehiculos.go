package crud

import (
	"fmt"
)

func MenuVehiculos() {
	for {
		clearScreen()
		fmt.Println("=== GESTIÓN DE VEHÍCULOS ===")
		fmt.Println("1. Crear vehículo")
		fmt.Println("2. Visualizar vehículo")
		fmt.Println("3. Modificar vehículo")
		fmt.Println("4. Eliminar vehículo")
		fmt.Println("5. Listar todos los vehículos")
		fmt.Println("0. Volver")
		fmt.Print("\nSeleccione opción: ")

		var opcion int
		fmt.Scan(&opcion)

		switch opcion {
		case 1:
			CrearVehiculo()
		case 2:
			VisualizarVehiculo()
		case 3:
			ModificarVehiculo()
		case 4:
			EliminarVehiculo()
		case 5:
			ListarVehiculos()
		case 0:
			return
		default:
			fmt.Println("Opción no válida")
			pausa()
		}
	}
}

func CrearVehiculo() {
	var v Vehiculo
	fmt.Print("Matrícula: ")
	fmt.Scan(&v.Matricula)
	fmt.Print("Marca: ")
	v.Marca = LeerLinea()
	fmt.Print("Modelo: ")
	v.Modelo = LeerLinea()
	fmt.Print("Fecha de entrada: ")
	fmt.Scan(&v.FechaEntrada)
	fmt.Print("Fecha estimada de salida: ")
	fmt.Scan(&v.FechaSalida)
	
	// Crear incidencia automáticamente para este vehículo
	incidenciaID := "I_" + v.Matricula
	fmt.Print("Tipo de incidencia (mecánica/eléctrica/carrocería): ")
	var tipo string
	fmt.Scan(&tipo)
	
	incidencia := Incidencia{
		ID:          incidenciaID,
		VehiculoID:  v.Matricula, 
		Tipo:        tipo,
		Prioridad:   "media",
		Descripcion: "Incidencia para " + v.Matricula,
		Estado:      "abierta",
		Mecanicos:   []string{},
	}
	
	Incidencias[incidenciaID] = incidencia
	v.IncidenciaID = incidenciaID

	Vehiculos[v.Matricula] = v

	// Asociar vehículo a cliente
	fmt.Print("ID del cliente propietario: ")
	var clienteID string
	fmt.Scan(&clienteID)

	if c, existe := Clientes[clienteID]; existe {
		c.Vehiculos = append(c.Vehiculos, v.Matricula)
		Clientes[clienteID] = c
	}

	fmt.Println("Vehículo e incidencia creados correctamente")
	pausa()
}

func VisualizarVehiculo() {
	fmt.Print("Matrícula del vehículo: ")
	var matricula string
	fmt.Scan(&matricula)

	v, existe := Vehiculos[matricula]
	if !existe {
		fmt.Println("Vehículo no encontrado")
		pausa()
		return
	}

	fmt.Printf("Matrícula: %s\n", v.Matricula)
	fmt.Printf("Marca: %s\n", v.Marca)
	fmt.Printf("Modelo: %s\n", v.Modelo)
	fmt.Printf("Fecha entrada: %s\n", v.FechaEntrada)
	fmt.Printf("Fecha salida: %s\n", v.FechaSalida)
	fmt.Printf("Incidencia ID: %s\n", v.IncidenciaID)
	pausa()
}

func ModificarVehiculo() {
	fmt.Print("Matrícula del vehículo a modificar: ")
	var matricula string
	fmt.Scan(&matricula)

	v, existe := Vehiculos[matricula]
	if !existe {
		fmt.Println("Vehículo no encontrado")
		pausa()
		return
	}

	fmt.Print("Nueva marca (actual: " + v.Marca + "): ")
	nuevaMarca := LeerLinea()
	if nuevaMarca != "" {
		v.Marca = nuevaMarca
	}

	fmt.Print("Nuevo modelo (actual: " + v.Modelo + "): ")
	nuevoModelo := LeerLinea()
	if nuevoModelo != "" {
		v.Modelo = nuevoModelo
	}

	fmt.Print("Nueva fecha entrada (actual: " + v.FechaEntrada + "): ")
	fmt.Scan(&v.FechaEntrada)

	fmt.Print("Nueva fecha salida (actual: " + v.FechaSalida + "): ")
	fmt.Scan(&v.FechaSalida)

	fmt.Print("Nuevo ID incidencia (actual: " + v.IncidenciaID + "): ")
	fmt.Scan(&v.IncidenciaID)

	Vehiculos[matricula] = v
	fmt.Println("Vehículo modificado correctamente")
	pausa()
}

func EliminarVehiculo() {
	fmt.Print("Matrícula del vehículo a eliminar: ")
	var matricula string
	fmt.Scan(&matricula)

	_, existe := Vehiculos[matricula]
	if !existe {
		fmt.Println("Vehículo no encontrado")
		pausa()
		return
	}

	// Liberar plaza si estaba asignada
	for i := range Plazas {
		if Plazas[i].Matricula == matricula {
			Plazas[i].Ocupada = false
			Plazas[i].Matricula = ""
		}
	}

	// Eliminar de cliente
	for id, c := range Clientes {
		for i, v := range c.Vehiculos {
			if v == matricula {
				c.Vehiculos = append(c.Vehiculos[:i], c.Vehiculos[i+1:]...)
				Clientes[id] = c
				break
			}
		}
	}

	delete(Vehiculos, matricula)
	fmt.Println("Vehículo eliminado correctamente")
	pausa()
}

func ListarVehiculos() {
	fmt.Println("=== LISTA DE VEHÍCULOS ===")
	for _, v := range Vehiculos {
		fmt.Printf("Matrícula: %s, Marca: %s, Modelo: %s\n", v.Matricula, v.Marca, v.Modelo)
	}
	pausa()
}