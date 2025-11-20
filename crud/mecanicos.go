package crud

import (
	"fmt"
)

func MenuMecanicos() {
	for {
		clearScreen()
		fmt.Println("=== GESTIÓN DE MECÁNICOS ===")
		fmt.Println("1. Crear mecánico")
		fmt.Println("2. Visualizar mecánico")
		fmt.Println("3. Modificar mecánico")
		fmt.Println("4. Eliminar mecánico")
		fmt.Println("5. Listar todos los mecánicos")
		fmt.Println("6. Dar de alta/baja")
		fmt.Println("0. Volver")
		fmt.Print("\nSeleccione opción: ")

		var opcion int
		fmt.Scan(&opcion)

		switch opcion {
		case 1:
			CrearMecanico()
		case 2:
			VisualizarMecanico()
		case 3:
			ModificarMecanico()
		case 4:
			EliminarMecanico()
		case 5:
			ListarMecanicos()
		case 6:
			CambiarEstadoMecanico()
		case 0:
			return
		default:
			fmt.Println("Opción no válida")
			pausa()
		}
	}
}

func CrearMecanico() {
	var m Mecanico
	fmt.Print("ID del mecánico: ")
	fmt.Scan(&m.ID)
	fmt.Print("Nombre: ")
	m.Nombre = LeerLinea()

	// Selección de especialidad con validación
	for {
		fmt.Print("Especialidad (mecánica/eléctrica/carrocería): ")
		var especialidadStr string
		fmt.Scan(&especialidadStr)
		m.Especialidad = especialidadStr
		if m.Especialidad == "mecánica" || m.Especialidad == "eléctrica" || m.Especialidad == "carrocería" {
			break
		}
		fmt.Println("Especialidad no válida. Use: mecánica, eléctrica o carrocería")
	}

	fmt.Print("Años de experiencia: ")
	fmt.Scan(&m.AniosExp)
	m.Activo = true

	Mecanicos[m.ID] = m
	
	// Añadir 2 plazas para este mecánico
	for i := 0; i < 2; i++ {
		Plazas = append(Plazas, Plaza{
			Numero:     ProxPlaza,
			Ocupada:    false,
			Matricula:  "",
			MecanicoID: m.ID,
		})
		ProxPlaza++
	}
	
	fmt.Println("Mecánico creado correctamente y 2 plazas asignadas")
	pausa()
}

func VisualizarMecanico() {
	fmt.Print("ID del mecánico: ")
	var id string
	fmt.Scan(&id)

	m, existe := Mecanicos[id]
	if !existe {
		fmt.Println("Mecánico no encontrado")
		pausa()
		return
	}

	fmt.Printf("Nombre: %s\n", m.Nombre)
	fmt.Printf("Especialidad: %s\n", m.Especialidad)
	fmt.Printf("Años experiencia: %d\n", m.AniosExp)
	fmt.Printf("Activo: %t\n", m.Activo)
	
	// Mostrar plazas asignadas
	fmt.Println("Plazas asignadas:")
	for _, p := range Plazas {
		if p.MecanicoID == id {
			estado := "Libre"
			if p.Ocupada {
				estado = "Ocupada - " + p.Matricula
			}
			fmt.Printf("  - Plaza %d: %s\n", p.Numero, estado)
		}
	}
	
	pausa()
}

func ModificarMecanico() {
	fmt.Print("ID del mecánico a modificar: ")
	var id string
	fmt.Scan(&id)

	m, existe := Mecanicos[id]
	if !existe {
		fmt.Println("Mecánico no encontrado")
		pausa()
		return
	}

	fmt.Print("Nuevo nombre (actual: " + m.Nombre + "): ")
	nuevoNombre := LeerLinea()
	if nuevoNombre != "" {
		m.Nombre = nuevoNombre
	}

	// Modificar especialidad con validación
	for {
		fmt.Print("Nueva especialidad (actual: " + m.Especialidad + "): ")
		var especialidadStr string
		fmt.Scan(&especialidadStr)
		if especialidadStr == "" {
			break // Mantener el valor actual
		}
		if especialidadStr == "mecánica" || especialidadStr == "eléctrica" || especialidadStr == "carrocería" {
			m.Especialidad = especialidadStr
			break
		}
		fmt.Println("Especialidad no válida. Use: mecánica, eléctrica o carrocería")
	}

	fmt.Printf("Nuevos años experiencia (actual: %d): ", m.AniosExp)
	fmt.Scan(&m.AniosExp)

	Mecanicos[id] = m
	fmt.Println("Mecánico modificado correctamente")
	pausa()
}

func EliminarMecanico() {
	fmt.Print("ID del mecánico a eliminar: ")
	var id string
	fmt.Scan(&id)

	_, existe := Mecanicos[id]
	if !existe {
		fmt.Println("Mecánico no encontrado")
		pausa()
		return
	}

	// Verificar si está asignado a alguna incidencia
	for _, i := range Incidencias {
		for _, mid := range i.Mecanicos {
			if mid == id {
				fmt.Println("No se puede eliminar: está asignado a una incidencia")
				pausa()
				return
			}
		}
	}
	
	// Verificar si sus plazas están libres
	for _, p := range Plazas {
		if p.MecanicoID == id && p.Ocupada {
			fmt.Println("No se puede eliminar: alguna de sus plazas está ocupada")
			pausa()
			return
		}
	}
	
	// Eliminar plazas de este mecánico
	nuevasPlazas := make([]Plaza, 0)
	for _, p := range Plazas {
		if p.MecanicoID != id {
			nuevasPlazas = append(nuevasPlazas, p)
		}
	}
	Plazas = nuevasPlazas
	
	delete(Mecanicos, id)
	fmt.Println("Mecánico y sus plazas eliminados correctamente")
	pausa()
}

func ListarMecanicos() {
	fmt.Println("=== LISTA DE MECÁNICOS ===")
	for _, m := range Mecanicos {
		estado := "Activo"
		if !m.Activo {
			estado = "Inactivo"
		}
		fmt.Printf("ID: %s, Nombre: %s, Especialidad: %s, %s\n", 
			m.ID, m.Nombre, m.Especialidad, estado)
	}
	pausa()
}

func CambiarEstadoMecanico() {
	fmt.Print("ID del mecánico: ")
	var id string
	fmt.Scan(&id)

	m, existe := Mecanicos[id]
	if !existe {
		fmt.Println("Mecánico no encontrado")
		pausa()
		return
	}

	if m.Activo {
		// Dar de baja: eliminar plazas si están libres
		for _, p := range Plazas {
			if p.MecanicoID == id && p.Ocupada {
				fmt.Println("No se puede dar de baja: alguna de sus plazas está ocupada")
				pausa()
				return
			}
		}
		
		// Eliminar plazas
		nuevasPlazas := make([]Plaza, 0)
		for _, p := range Plazas {
			if p.MecanicoID != id {
				nuevasPlazas = append(nuevasPlazas, p)
			}
		}
		Plazas = nuevasPlazas
		
		m.Activo = false
		fmt.Printf("Mecánico %s dado de baja y plazas eliminadas\n", m.Nombre)
	} else {
		// Dar de alta: añadir 2 plazas
		for i := 0; i < 2; i++ {
			Plazas = append(Plazas, Plaza{
				Numero:     ProxPlaza,
				Ocupada:    false,
				Matricula:  "",
				MecanicoID: id,
			})
			ProxPlaza++
		}
		m.Activo = true
		fmt.Printf("Mecánico %s dado de alta y 2 plazas añadidas\n", m.Nombre)
	}
	
	Mecanicos[id] = m
	pausa()
}