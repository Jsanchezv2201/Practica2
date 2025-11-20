package crud

import (
	"fmt"
)

func MenuIncidencias() {
	for {
		clearScreen()
		fmt.Println("=== GESTIÓN DE INCIDENCIAS ===")
		fmt.Println("1. Crear incidencia")
		fmt.Println("2. Visualizar incidencia")
		fmt.Println("3. Modificar incidencia")
		fmt.Println("4. Eliminar incidencia")
		fmt.Println("5. Listar todas las incidencias")
		fmt.Println("0. Volver")
		fmt.Print("\nSeleccione opción: ")

		var opcion int
		fmt.Scan(&opcion)

		switch opcion {
		case 1:
			CrearIncidencia()
		case 2:
			VisualizarIncidencia()
		case 3:
			ModificarIncidencia()
		case 4:
			EliminarIncidencia()
		case 5:
			ListarIncidencias()
		case 0:
			return
		default:
			fmt.Println("Opción no válida")
			pausa()
		}
	}
}

func CrearIncidencia() {
	var i Incidencia
	fmt.Print("ID de la incidencia: ")
	fmt.Scan(&i.ID)

	// AÑADIR: Solicitar VehiculoID
	fmt.Print("ID del vehículo: ")
	fmt.Scan(&i.VehiculoID)

	// Selección de tipo con validación
	for {
		fmt.Print("Tipo (mecánica/eléctrica/carrocería): ")
		var tipoStr string
		fmt.Scan(&tipoStr)
		i.Tipo = tipoStr
		if i.Tipo == "mecánica" || i.Tipo == "eléctrica" || i.Tipo == "carrocería" {
			break
		}
		fmt.Println("Tipo no válido. Use: mecánica, eléctrica o carrocería")
	}

	// Selección de prioridad con validación
	for {
		fmt.Print("Prioridad (baja/media/alta): ")
		var prioridadStr string
		fmt.Scan(&prioridadStr)
		i.Prioridad = prioridadStr
		if i.Prioridad == "baja" || i.Prioridad == "media" || i.Prioridad == "alta" {
			break
		}
		fmt.Println("Prioridad no válida. Use: baja, media o alta")
	}

	fmt.Print("Descripción: ")
	i.Descripcion = LeerLinea()
	i.Estado = "abierta"
	i.Mecanicos = []string{}

	Incidencias[i.ID] = i
	fmt.Println("Incidencia creada correctamente")
	pausa()
}

func VisualizarIncidencia() {
	fmt.Print("ID de la incidencia: ")
	var id string
	fmt.Scan(&id)

	i, existe := Incidencias[id]
	if !existe {
		fmt.Println("Incidencia no encontrada")
		pausa()
		return
	}

	fmt.Printf("Tipo: %s\n", i.Tipo)
	fmt.Printf("Prioridad: %s\n", i.Prioridad)
	fmt.Printf("Descripción: %s\n", i.Descripcion)
	fmt.Printf("Estado: %s\n", i.Estado)
	fmt.Printf("Mecánicos asignados: %v\n", i.Mecanicos)
	pausa()
}

func ModificarIncidencia() {
	fmt.Print("ID de la incidencia a modificar: ")
	var id string
	fmt.Scan(&id)

	i, existe := Incidencias[id]
	if !existe {
		fmt.Println("Incidencia no encontrada")
		pausa()
		return
	}

	// Modificar VehiculoID
	fmt.Print("Nuevo ID del vehículo (actual: " + i.VehiculoID + "): ")
	fmt.Scan(&i.VehiculoID)

	// Modificar tipo con validación
	for {
		fmt.Print("Nuevo tipo (actual: " + i.Tipo + "): ")
		var tipoStr string
		fmt.Scan(&tipoStr)
		if tipoStr == "" {
			break // Mantener el valor actual
		}
		if tipoStr == "mecánica" || tipoStr == "eléctrica" || tipoStr == "carrocería" {
			i.Tipo = tipoStr
			break
		}
		fmt.Println("Tipo no válido. Use: mecánica, eléctrica o carrocería")
	}

	// Modificar prioridad con validación
	for {
		fmt.Print("Nueva prioridad (actual: " + i.Prioridad + "): ")
		var prioridadStr string
		fmt.Scan(&prioridadStr)
		if prioridadStr == "" {
			break // Mantener el valor actual
		}
		if prioridadStr == "baja" || prioridadStr == "media" || prioridadStr == "alta" {
			i.Prioridad = prioridadStr
			break
		}
		fmt.Println("Prioridad no válida. Use: baja, media o alta")
	}

	fmt.Print("Nueva descripción (actual: " + i.Descripcion + "): ")
	nuevaDesc := LeerLinea()
	if nuevaDesc != "" {
		i.Descripcion = nuevaDesc
	}

	// Modificar estado con validación
	for {
		fmt.Print("Nuevo estado (actual: " + i.Estado + "): ")
		var estadoStr string
		fmt.Scan(&estadoStr)
		if estadoStr == "" {
			break // Mantener el valor actual
		}
		if estadoStr == "abierta" || estadoStr == "en proceso" || estadoStr == "cerrada" {
			i.Estado = estadoStr
			break
		}
		fmt.Println("Estado no válido. Use: abierta, en proceso o cerrada")
	}

	Incidencias[id] = i
	fmt.Println("Incidencia modificada correctamente")
	pausa()
}

func EliminarIncidencia() {
	fmt.Print("ID de la incidencia a eliminar: ")
	var id string
	fmt.Scan(&id)

	_, existe := Incidencias[id]
	if !existe {
		fmt.Println("Incidencia no encontrada")
		pausa()
		return
	}

	// Verificar si está asignada a algún vehículo
	for _, v := range Vehiculos {
		if v.IncidenciaID == id {
			fmt.Println("No se puede eliminar: está asignada a un vehículo")
			pausa()
			return
		}
	}

	delete(Incidencias, id)
	fmt.Println("Incidencia eliminada correctamente")
	pausa()
}

func ListarIncidencias() {
	fmt.Println("=== LISTA DE INCIDENCIAS ===")
	for _, i := range Incidencias {
		fmt.Printf("ID: %s, Tipo: %s, Prioridad: %s, Estado: %s\n", i.ID, i.Tipo, i.Prioridad, i.Estado)
	}
	pausa()
}