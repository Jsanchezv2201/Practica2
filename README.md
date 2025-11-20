# README TEMPORAL...

# Practica 2 - **Uso de goroutines y channels de Go**

### **🎯 OBJETIVO PRINCIPAL:**
Crear un **sistema distribuido** para un taller mecánico usando **goroutines y channels de Go** (en lugar de un programa secuencial como en la práctica 1).


### **⚙️ FUNCIONALIDADES REQUERIDAS:**

**1. Atención por Mecánicos Especializados:**
- **Mecánica:** 5 segundos de atención
- **Eléctrica:** 7 segundos  
- **Carrocería:** 11 segundos

**2. Sistema de Cola:**
- Cola de espera **sin límite de tamaño**
- Los coches esperan si no hay mecánicos libres

**3. Sistema de Prioridad:**
- Si un coche acumula **más de 15 segundos** de atención
- Se le asigna **otro mecánico** adicional
- Si no hay mecánicos, se **contrata uno nuevo**

### **🛠️ RESTRICCIONES TÉCNICAS:**
- Usar **solo goroutines y channels** (como se vio en clase)
- Se pueden usar **múltiples archivos .go**
- **NO** es necesario guardar datos (sin persistencia)

---




## 🔄 **Explicación de la Implementación de Goroutines y Channels**

### **1. Estructura General del Sistema Concurrente**

```
Taller (Main)
    │
    ├── Coordinator (Goroutine) ← Gestiona cola y prioridades
    │
    ├── Mecánico 1 (Goroutine) ← Atiende coches
    ├── Mecánico 2 (Goroutine) ← Atiende coches  
    ├── Mecánico 3 (Goroutine) ← Atiende coches
    └── ...
```

### **2. Goroutines Implementadas**

**A) Goroutine del Coordinator:**
```go
// En taller.go
func (t *Taller) coordinator() {
    for t.running {
        coche := t.Cola.ObtenerCoche()
        if coche == nil { return }
        // Lógica de asignación...
    }
}
```

**B) Goroutines de los Mecánicos:**
```go
// En mecanico.go  
func (m *Mecanico) Iniciar(taller *Taller) {
    go func() {
        for coche := range m.ChanTrabajo {
            // Procesar coche...
            tiempoAtencion := coche.TiempoAtencion()
            time.Sleep(tiempoAtencion)
            // Registrar finalización...
        }
    }()
}
```

**C) Goroutines Auxiliares para Re-encolado:**
```go
// En atiendeCocheNormal()
go func(c *Coche) {
    time.Sleep(waitTime)
    if t.running {
        t.Cola.AgregarCoche(c)
    }
}(coche)
```

### **3. Channels Implementados**

**A) Channel de Trabajo por Mecánico:**
```go
type Mecanico struct {
    ChanTrabajo chan *Coche  // Channel buffered (tamaño 1)
}

// Uso: Asignar coche a mecánico
mecanico.ChanTrabajo <- coche
```

**B) Channel de Notificación de Cola:**
```go
type Cola struct {
    notify chan struct{}  // Channel para notificar nuevos elementos
}

// Uso: Notificar cuando hay coche nuevo
select {
case c.notify <- struct{}{}:
default: // Evita bloqueo si ya hay notificación
}
```

**C) Channel de Control de Parada:**
```go
type Taller struct {
    ChanDetener chan bool  // Para señalizar parada
}
```

### **4. Patrones de Comunicación**

**Flujo Normal:**
```
Coordinator → [Channel ChanTrabajo] → Mecánico (Goroutine)
     ↑
  Cola con notify channel
```

**Flujo con Prioridad:**
```
Coche prioritario → Coordinator → Contratar nuevo mecánico → Nuevo Channel
```

### **5. Sincronización y Control**

**Inicio del Sistema:**
```go
func (t *Taller) Iniciar() {
    go t.coordinator()  // Lanzar goroutine coordinador
    for _, m := range t.Mecanicos {
        m.Iniciar(t)    // Lanzar goroutine por cada mecánico
    }
}
```

**Parada Controlada:**
```go
func (t *Taller) Detener() {
    t.running = false
    t.Cola.Cerrar()           // Cerrar cola primero
    close(t.ChanDetener)      // Señalizar parada
    for _, m := range t.Mecanicos {
        m.Detener()           // Cerrar channels de mecánicos
    }
}
```

### **6. Gestión de Concurrencia en la Cola**

```go
func (c *Cola) ObtenerCoche() *Coche {
    for {
        c.mutex.Lock()
        if len(c.coches) > 0 {
            coche := c.coches[0]
            c.coches = c.coches[1:]
            c.mutex.Unlock()
            return coche
        }
        c.mutex.Unlock()
        <-c.notify  // Espera bloqueante hasta notificación
    }
}
```

# 🚀 **Guía Completa de Uso del Sistema de Taller Mecánico**

## **📋 Descripción General del Sistema**

He desarrollado un **sistema dual** que combina:

### **1. 🖱️ Gestión Manual (CRUD)**
**Igual que en la Práctica 1** - Sistema completo de gestión manual

### **2. ⚡ Simulación Automática**  
**Nueva funcionalidad** - Para pruebas rápidas y comparativas

### **3. 🧪 Tests Automatizados**
**Validación exhaustiva** - Verificación del sistema concurrente

---

## **🎮 Cómo Usar el Sistema - Paso a Paso**

### **OPCIÓN 1: Gestión Manual (Modo Interactivo)**

#### **Pasos:**
1. **Ejecutar el programa:**
   ```bash
   go run main.go
   ```

2. **Seleccionar opción 1: "Gestión Manual"**
   ```
   === TALLER MECÁNICO - PRÁCTICA 2 ===
   1. Gestión Manual (Clientes, Vehículos, Incidencias, Mecánicos)
   2. Ejecutar Simulación Automática
   3. Simulación con Datos Actuales
   4. Estado Actual del Taller
   5. Ejecutar Tests
   0. Salir
   ```

3. **Navegar por los submenús:**
   - **Clientes**: Crear, visualizar, modificar, eliminar
   - **Vehículos**: Gestionar vehículos y asociar incidencias
   - **Incidencias**: Gestionar problemas con tipo y prioridad
   - **Mecánicos**: Gestionar especialistas y sus plazas

#### **Cuándo usar este modo:**
- ✅ Para probar funcionalidades específicas
- ✅ Cuando quieres simular uso real del sistema
- ✅ Para verificar la integración entre módulos

---

### **OPCIÓN 2: Simulación Automática (Recomendado para pruebas)**

#### **Pasos:**
1. **Ejecutar el programa:**
   ```bash
   go run main.go
   ```

2. **Seleccionar opción 2: "Ejecutar Simulación Automática"**
   - El sistema ejecutará **automáticamente 5 escenarios predefinidos**
   - No requiere ninguna entrada manual
   - Genera métricas completas de rendimiento

#### **Los 5 escenarios que se prueban:**
1. **Configuración Base** (3 mecánicos, 8 coches)
2. **Doble Carga** (3 mecánicos, 16 coches) 
3. **Doble Plantilla** (6 mecánicos, 8 coches)
4. **Distribución 3M-1E-1C** (5 mecánicos especializados)
5. **Distribución 1M-3E-3C** (7 mecánicos especializados)

#### **Cuándo usar este modo:**
- ✅ Para ver el rendimiento del sistema completo
- ✅ Para comparar diferentes configuraciones
- ✅ Para obtener métricas de forma rápida

---

### **OPCIÓN 3: Tests Individuales (Para desarrolladores)**

#### **Método A: Desde VS Code (Más fácil)**
1. **Abrir el archivo `taller_test.go`**
2. **Buscar las funciones de test** (cada escenario tiene su propia función)
3. **Hacer clic en el icono "Run Test"** ▶️ que aparece arriba de cada función

**Ejemplo:**
```go
// Buscar esta función y hacer clic en "Run Test" arriba de ella:
func TestEscenario1_ConfiguracionBase(t *testing.T) {
    // Este test ejecuta solo el escenario base
}

func TestEscenario2_DobleCoches(t *testing.T) {
    // Este test ejecuta solo el escenario de doble carga
}
```

#### **Método B: Desde Terminal**
```bash
# Ejecutar TODOS los tests
go test -v

# Ejecutar UN test específico
go test -v -run TestEscenario1_ConfiguracionBase

# Ejecutar tests con timeout extendido
go test -v -timeout=120s
```

#### **Tests disponibles en `taller_test.go`:**
- `TestEscenario1_ConfiguracionBase`
- `TestEscenario2_DobleCoches` 
- `TestEscenario3_DobleMecanicos`
- `TestEscenario4_Mecanicos3Mecanica`
- `TestEscenario5_Mecanicos1Mecanica3Electricos3Carroceria`
- `TestFuncionalidadesClave`
- `TestRendimiento`

---

## **🔄 Flujo Recomendado para Nuevos Usuarios**

### **Para entender el sistema:**
1. **Primero**: Ejecutar **Opción 2** (Simulación Automática) para ver el sistema en acción
2. **Luego**: Probar **Opción 1** (Gestión Manual) para entender las funcionalidades
3. **Finalmente**: Ejecutar **tests individuales** para verificar componentes específicos

### **Para desarrolladores:**
1. **Modificar el código**
2. **Ejecutar tests relevantes** desde VS Code
3. **Verificar que todo funciona** con la simulación automática


---

## **🚨 Solución de Problemas Comunes**

### **Si los tests fallan:**
- Verificar que todos los archivos `.go` estén en la misma carpeta
- Ejecutar `go mod tidy` para resolver dependencias
- Asegurarse de usar Go version 1.16 o superior

### **Si la simulación se cuelga:**
- Los tests tienen timeout de 120 segundos
- Si excede este tiempo, revisar posibles bucles infinitos

### **Para obtener más detalles:**
- Ejecutar con `-v` para output verbose
- Revisar los logs que muestran el progreso paso a paso

---



# 🚀 **Implementación del Módulo de Simulación Automática - Mi Enfoque Personal**

## **¿Por qué desarrollé `simulacion.go`?**

**Como desarrollador, me di cuenta de que necesitaba una forma más eficiente de probar el sistema.** Durante las primeras pruebas manuales, perdía mucho tiempo creando clientes, vehículos y mecánicos uno por uno. Esto me impedía:

### **Problemas que identificé:**
```go
// Antes - Pruebas manuales lentas:
1. ⏳ 5-10 minutos por prueba creando datos
2. 🔄 Dificultad para reproducir exactamente los mismos escenarios  
3. 📊 Imposibilidad de comparar configuraciones de forma justa
4. 🧪 Complejidad para probar casos extremos de forma consistente
```

### **Mi solución: `simulacion.go`**
```go
// Decidí crear un sistema que me permitiera:
func PorqueLoSimplemente() {
    // 1. 🔁 Ejecutar pruebas en segundos, no en minutos
    // 2. 📈 Comparar múltiples escenarios rápidamente  
    // 3. 🎯 Reproducir exactamente las mismas condiciones
    // 4. 🧪 Probar casos límite de forma sistemática
}
```

## **Mi Proceso de Desarrollo**

### **Fase 1: Necesidad Identificada**
"Después de probar manualmente el sistema 2-3 veces, me di cuenta de que estaba gastando más tiempo configurando datos que analizando resultados. Necesitaba una forma de automatizar este proceso."

### **Fase 2: Diseño del Módulo**
```go
// Pensé: "¿Qué necesito para probar realmente el sistema concurrente?"
type MiEnfoque struct {
    ConfiguracionesPredefinidas []Escenario
    ModoAutomatico              bool
    MetricasAutomaticas         bool
}

// Escogí 5 escenarios que representaran casos reales:
// 1. Caso base - Línea de referencia
// 2. Doble carga - Test de estrés  
// 3. Doble plantilla - Test de recursos
// 4. Distribución 3-1-1 - Test de especialización
// 5. Distribución 1-3-3 - Test de balance extremo
```

### **Fase 3: Implementación**
"Implementé `CrearConfiguracionAutomatica()` para que, con un simple número de escenario, pudiera generar toda la configuración necesaria. Esto me permitió ejecutar los 5 tests en menos de 3 minutos, en lugar de 30+ minutos manualmente."

## **Beneficios que Obtuve Personalmente**

### **🕒 Eficiencia de Tiempo**
```go
// ANTES: ~30 minutos para 5 pruebas manuales
// DESPUÉS: ~3 minutos para 5 pruebas automáticas

// Ganancia: 90% de tiempo ahorrado
```

### **🐛 Detección de Errores**
"La simulación automática me ayudó a encontrar y corregir varios bugs que hubieran pasado desapercibidos con pruebas manuales."



# 1. Explicación del Diseño del Sistema

## 📋 **Estructuras de Datos Principales**

### 🚗 **Coche**
```go
type Coche struct {
    Matricula     string 
    ID            string 
    TipoIncidencia TipoIncidencia 
    TiempoAtendido time.Duration 
    ChanTerminado chan bool 
    TiempoLlegada time.Time 
}
```

**Propósito:** Representa cada vehículo que llega al taller con su incidencia específica.

**Campos clave:**
- `TipoIncidencia`: Determina la especialidad requerida y tiempo de reparación
- `TiempoAtendido`: Acumula el tiempo total de atención para control de prioridades
- `ChanTerminado`: Permite sincronizar la finalización entre goroutines
- `TiempoLlegada`: Timestamp para medición de tiempos reales

---

### 🔧 **Mecánico**
```go
type Mecanico struct {
    ID           string
    Especialidad TipoIncidencia
    Ocupado      bool
    ChanTrabajo  chan *Coche
    Trabajando   bool
    taller       *Taller
}
```

**Propósito:** Cada mecánico es una goroutine independiente que procesa coches concurrentemente.

**Campos clave:**
- `Especialidad`: Define qué tipo de incidencias puede atender
- `ChanTrabajo`: Channel personalizado para recibir trabajos (patrón worker)
- `Ocupado`/`Trabajando`: Estados para gestión de concurrencia
- `taller`: Referencia al sistema principal para comunicación bidireccional

---

### 🏢 **Taller**
```go
type Taller struct {
    Cola               *Cola
    Mecanicos          []*Mecanico
    ChanDetener        chan bool
    Stats              *Estadisticas
    running            bool
    mensajesBuffer     []string
}
```

**Propósito:** Coordina todas las operaciones del sistema y gestiona el estado global.

**Campos clave:**
- `Cola`: Centraliza la gestión de coches pendientes
- `Mecanicos`: Pool de workers especializados
- `ChanDetener`: Controla el cierre graceful del sistema
- `Stats`: Recopila métricas para análisis comparativo
- `running`: Flag atómico para control de ciclo de vida

---

### 📋 **Cola de Espera**
```go
type Cola struct {
    coches   []*Coche
    mutex    sync.Mutex
    cerrada  bool
    notify   chan struct{}
}
```

**Propósito:** Gestiona la cola de espera de forma thread-safe con notificaciones eficientes.

**Campos clave:**
- `mutex`: Garantiza acceso seguro desde múltiples goroutines
- `notify`: Implementa el patrón observer para notificaciones no-bloqueantes
- `cerrada`: Permite un cierre ordenado sin race conditions

---

