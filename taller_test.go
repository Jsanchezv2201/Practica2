package main

import (
	"practica2/taller"
	"testing"
)

// TestEscenario1_ConfiguracionBase - Escenario base de referencia
func TestEscenario1_ConfiguracionBase(t *testing.T) {
	t.Log("\n📊 ESCENARIO 1: CONFIGURACIÓN BASE (REFERENCIA)")
	config := taller.CrearConfiguracionAutomatica(1)
	stats, duracion := taller.EjecutarSimulacion(config)

	// Verificaciones principales
	if stats.CochesTotales != 8 {
		t.Errorf("ERROR: Se esperaban 8 coches totales, se obtuvieron %d", stats.CochesTotales)
	}

	t.Logf("\n✅ RESULTADO BASE - Duración: %v, Atendidos: %d/%d, Mecánicos extra: %d, Prioritarios: %d", 
		duracion, stats.CochesAtendidos, stats.CochesTotales, 
		stats.MecanicosContratados, stats.CochesPrioritarios)

	// Verificación de funcionamiento correcto
	if stats.CochesAtendidos == 0 {
		t.Errorf("ERROR: No se atendió ningún coche")
	}
}

// TestEscenario2_DobleCoches - Duplicar cantidad de coches (mismos mecánicos)
func TestEscenario2_DobleCoches(t *testing.T) {
	t.Log("\n🚗 ESCENARIO 2: DUPLICAR CANTIDAD DE COCHES (16 coches, 4 mecánicos)")
	config := taller.CrearConfiguracionAutomatica(2)
	stats, duracion := taller.EjecutarSimulacion(config)

	// Verificaciones para escenario duplicado
	if stats.CochesTotales != 16 {
		t.Errorf("ERROR: Se esperaban 16 coches totales (doble), se obtuvieron %d", stats.CochesTotales)
	}

	t.Logf("\n✅ RESULTADO DOBLE COCHES - Duración: %v, Atendidos: %d/%d, Mecánicos extra: %d, Prioritarios: %d", 
		duracion, stats.CochesAtendidos, stats.CochesTotales, 
		stats.MecanicosContratados, stats.CochesPrioritarios)

	// Análisis comparativo
	if stats.MecanicosContratados > 0 {
		t.Logf("\n💡 OBSERVACIÓN: Con 16 coches se contrataron %d mecánicos extra", stats.MecanicosContratados)
	}
}

// TestEscenario3_DobleMecanicos - Duplicar plantilla de mecánicos
func TestEscenario3_DobleMecanicos(t *testing.T) {
	t.Log("\n👥 ESCENARIO 3: DUPLICAR PLANTILLA DE MECÁNICOS (6 mecánicos, 12 coches)")
	config := taller.CrearConfiguracionAutomatica(4)
	stats, duracion := taller.EjecutarSimulacion(config)

	// Verificaciones para escenario duplicado mecánicos
	if stats.CochesTotales != 8 {
		t.Errorf("ERROR: Se esperaban 8 coches totales, se obtuvieron %d", stats.CochesTotales)
	}

	t.Logf("\n✅ RESULTADO DOBLE MECÁNICOS - Duración: %v, Atendidos: %d/%d, Mecánicos extra: %d, Prioritarios: %d", 
		duracion, stats.CochesAtendidos, stats.CochesTotales, 
		stats.MecanicosContratados, stats.CochesPrioritarios)

	// Con más mecánicos debería haber menos contrataciones extra
	if stats.MecanicosContratados == 0 {
		t.Log("\n✅ OBSERVACIÓN: Con 6 mecánicos base, no se necesitaron contrataciones extra")
	}
}

// TestEscenario4_Mecanicos3Mecanica - 3 mecánicos mecánica / 1 eléctrica / 1 carrocería
func TestEscenario4_Mecanicos3Mecanica(t *testing.T) {
	t.Log("\n🔧 ESCENARIO 4: 3 MECÁNICA / 1 ELÉCTRICA / 1 CARROCERÍA (5 mecánicos, 8 coches)")
	config := taller.CrearConfiguracionAutomatica(3)
	stats, duracion := taller.EjecutarSimulacion(config)

	// Verificaciones para distribución 3-1-1
	if stats.CochesTotales != 8 {
		t.Errorf("ERROR: Se esperaban 8 coches totales, se obtuvieron %d", stats.CochesTotales)
	}

	t.Logf("\n✅ RESULTADO 3-1-1 - Duración: %v, Atendidos: %d/%d, Mecánicos extra: %d, Prioritarios: %d", 
		duracion, stats.CochesAtendidos, stats.CochesTotales, 
		stats.MecanicosContratados, stats.CochesPrioritarios)

	// Análisis de especialización
	if stats.CochesPrioritarios > 0 {
		t.Logf("\n💡 OBSERVACIÓN: En distribución 3-1-1, hubo %d coches prioritarios", stats.CochesPrioritarios)
	}
}

// TestEscenario5_Mecanicos1Mecanica3Electricos3Carroceria - 1 mecánica / 3 eléctrica / 3 carrocería
func TestEscenario5_Mecanicos1Mecanica3Electricos3Carroceria(t *testing.T) {
	t.Log("\n⚖️ ESCENARIO 5: 1 MECÁNICA / 3 ELÉCTRICA / 3 CARROCERÍA (7 mecánicos, 10 coches)")
	config := taller.CrearConfiguracionAutomatica(5)
	stats, duracion := taller.EjecutarSimulacion(config)

	// Verificaciones para distribución 1-3-3
	if stats.CochesTotales != 8 {
		t.Errorf("ERROR: Se esperaban 8 coches totales, se obtuvieron %d", stats.CochesTotales)
	}

	t.Logf("\n✅ RESULTADO 1-3-3 - Duración: %v, Atendidos: %d/%d, Mecánicos extra: %d, Prioritarios: %d", 
		duracion, stats.CochesAtendidos, stats.CochesTotales, 
		stats.MecanicosContratados, stats.CochesPrioritarios)

	// Análisis de distribución extrema
	if stats.MecanicosContratados == 0 && stats.CochesPrioritarios == 0 {
		t.Log("✅ EXCELENTE: Distribución 1-3-3 funcionó perfectamente sin contrataciones extra ni prioridades")
	}
}