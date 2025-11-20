package crud

import "fmt"

func MenuClientes() {
	for {
		clearScreen()
		fmt.Println("=== GESTIÓN DE CLIENTES ===")
		fmt.Println("1. Crear cliente")
		fmt.Println("2. Visualizar cliente")
		fmt.Println("3. Modificar cliente")
		fmt.Println("4. Eliminar cliente")
		fmt.Println("5. Listar todos los clientes")
		fmt.Println("0. Volver")
		fmt.Print("\nSeleccione opción: ")

		var opcion int
		fmt.Scan(&opcion)

		switch opcion {
		case 1:
			CrearCliente()
		case 2:
			VisualizarCliente()
		case 3:
			ModificarCliente()
		case 4:
			EliminarCliente()
		case 5:
			ListarClientes()
		case 0:
			return
		default:
			fmt.Println("Opción no válida")
			pausa()
		}
	}
}

func CrearCliente() {
	var c Cliente
	fmt.Print("ID del cliente: ")
	fmt.Scan(&c.ID)
	fmt.Print("Nombre: ")
	c.Nombre = LeerLinea()
	fmt.Print("Teléfono: ")
	fmt.Scan(&c.Telefono)
	fmt.Print("Email: ")
	fmt.Scan(&c.Email)
	c.Vehiculos = []string{}

	Clientes[c.ID] = c
	fmt.Println("Cliente creado correctamente")
	pausa()
}

func VisualizarCliente() {
	fmt.Print("ID del cliente: ")
	var id string
	fmt.Scan(&id)

	c, existe := Clientes[id]
	if !existe {
		fmt.Println("Cliente no encontrado")
		pausa()
		return
	}

	fmt.Printf("Nombre: %s\n", c.Nombre)
	fmt.Printf("Teléfono: %s\n", c.Telefono)
	fmt.Printf("Email: %s\n", c.Email)
	fmt.Printf("Vehículos: %v\n", c.Vehiculos)
	pausa()
}

func ModificarCliente() {
	fmt.Print("ID del cliente a modificar: ")
	var id string
	fmt.Scan(&id)

	c, existe := Clientes[id]
	if !existe {
		fmt.Println("Cliente no encontrado")
		pausa()
		return
	}

	fmt.Print("Nuevo nombre (actual: " + c.Nombre + "): ")
	c.Nombre = LeerLinea()
	fmt.Print("Nuevo teléfono (actual: " + c.Telefono + "): ")
	fmt.Scan(&c.Telefono)
	fmt.Print("Nuevo email (actual: " + c.Email + "): ")
	fmt.Scan(&c.Email)

	Clientes[id] = c
	fmt.Println("Cliente modificado correctamente")
	pausa()
}

func EliminarCliente() {
	fmt.Print("ID del cliente a eliminar: ")
	var id string
	fmt.Scan(&id)

	_, existe := Clientes[id]
	if !existe {
		fmt.Println("Cliente no encontrado")
		pausa()
		return
	}

	// Verificar si tiene vehículos en el taller
	c := Clientes[id]
	for _, matricula := range c.Vehiculos {
		if v, existe := Vehiculos[matricula]; existe {
			fmt.Printf("No se puede eliminar: tiene vehículo %s en el taller\n", v.Matricula)
			pausa()
			return
		}
	}

	delete(Clientes, id)
	fmt.Println("Cliente eliminado correctamente")
	pausa()
}

func ListarClientes() {
	fmt.Println("=== LISTA DE CLIENTES ===")
	for _, c := range Clientes {
		fmt.Printf("ID: %s, Nombre: %s, Tel: %s, Email: %s\n", c.ID, c.Nombre, c.Telefono, c.Email)
	}
	pausa()
}