package riInfluxdb

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// OrderConsumptionCalculate con verificación antes de consultar
func (m *repository) OrderConsumptionCalculate(factory string, prodline string, components []MriDosingComponent, startedAt time.Time, finishedAt time.Time) ([]MriDosingComponent, error) {
	logger := m.logger.With(
		zap.String("factory", factory),
		zap.String("prodline", prodline),
		zap.Time("startedAt", startedAt),
		zap.Time("finishedAt", finishedAt),
		zap.Int("components_count", len(components)),
	)

	fmt.Println("🚀 ===== INICIANDO OrderConsumptionCalculate =====")
	fmt.Printf("📊 Parámetros: Factory=%s, ProdLine=%s, Componentes=%d\n", factory, prodline, len(components))

	// Obtener el cliente InfluxDB para este factory - CON VERIFICACIÓN EN TIEMPO REAL
	client, connectionName, err := m.getClient(factory)
	if err != nil {
		errorMsg := fmt.Sprintf("❌ ERROR: No se pudo obtener cliente InfluxDB para fábrica %s: %v", factory, err)
		fmt.Println(errorMsg)
		logger.Error("Failed to get InfluxDB client for factory",
			zap.String("factory", factory),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get InfluxDB client for factory %s: %w", factory, err)
	}

	fmt.Printf("✅ Conexión específica obtenida: %s para fábrica %s\n", connectionName, factory)

	// Obtener la configuración
	cfg, exists := m.influxManager.GetConfig(connectionName)
	if !exists {
		errorMsg := fmt.Sprintf("❌ ERROR: Configuración no encontrada para conexión: %s", connectionName)
		fmt.Println(errorMsg)
		return nil, fmt.Errorf("configuración no encontrada para conexión: %s", connectionName)
	}

	fmt.Printf("🔧 Usando: Conexión=%s, Org=%s, Bucket=%s\n", connectionName, cfg.Org, cfg.Bucket)

	// VERIFICACIÓN FINAL antes de ejecutar la consulta
	m.logger.Info("🔍 Verificación final antes de ejecutar consulta",
		zap.String("connection", connectionName),
	)

	if !m.influxManager.IsHealthy(connectionName) {
		errorMsg := fmt.Sprintf("❌ ERROR: Conexión se cayó justo antes de ejecutar consulta para fábrica %s", factory)
		fmt.Println(errorMsg)
		return nil, fmt.Errorf("conexión InfluxDB no disponible para consulta en fábrica: %s", factory)
	}

	// Construir y ejecutar consulta
	query := fmt.Sprintf(
		`from(bucket: "%s")
		|> range(start: %s, stop: %s)
		|> filter(fn: (r) => r["name"] == "%s")
		|> filter(fn: (r) => r["EL_Lv1"] == "2_Dosing")
		|> filter(fn: (r) => r["EL_Lv2"] == "Doser_01" or r["EL_Lv2"] == "Doser_02" or r["EL_Lv2"] == "Doser_03" or r["EL_Lv2"] == "Doser_04" or r["EL_Lv2"] == "Doser_05" or r["EL_Lv2"] == "Doser_06")
		|> filter(fn: (r) => r["EL_Lv3"] == "C1" or r["EL_Lv3"] == "C2" or r["EL_Lv3"] == "C3" or r["EL_Lv3"] == "C4" or r["EL_Lv3"] == "C5" or r["EL_Lv3"] == "C6" or r["EL_Lv3"] == "C7")
		|> filter(fn: (r) => r["_field"] == "Dispensed_gr")
		|> cumulativeSum()
		|> last()
		|> pivot(rowKey: ["EL_Lv2"], columnKey: ["EL_Lv3"], valueColumn: "_value")`,
		cfg.Bucket, startedAt.Format(time.RFC3339), finishedAt.Format(time.RFC3339), prodline)

	fmt.Println("📋 Consulta InfluxDB construida correctamente")
	fmt.Println(query)

	// Ejecutar la consulta con timeout
	queryAPI := client.QueryAPI(cfg.Org)

	// Usar contexto con timeout para evitar bloqueos largos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		errorMsg := fmt.Sprintf("❌ ERROR: Falló la consulta a InfluxDB para fábrica %s: %v", factory, err)
		fmt.Println(errorMsg)
		logger.Error("Failed to execute InfluxDB query for factory",
			zap.String("factory", factory),
			zap.String("connection", connectionName),
			zap.Error(err))
		return nil, fmt.Errorf("failed to query InfluxDB for factory %s (connection: %s): %w", factory, connectionName, err)
	}

	fmt.Println("✅ Consulta ejecutada exitosamente")

	// Procesar resultados (mantener tu lógica existente)
	influxData := make(map[string]map[string]float64)
	var recordCount int

	for result.Next() {
		record := result.Record()
		dosingUnit, ok := record.ValueByKey("EL_Lv2").(string)
		if !ok {
			continue
		}

		influxData[dosingUnit] = make(map[string]float64)

		// Extraer valores de los hoppers
		hoppers := []string{"C1", "C2", "C3", "C4", "C5", "C6", "C7"}
		for _, hopper := range hoppers {
			if value, ok := record.ValueByKey(hopper).(float64); ok {
				influxData[dosingUnit][hopper] = value / 1000 // Convertir a kg
			}
		}
		recordCount++
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("error processing query results for factory %s: %w", factory, result.Err())
	}

	// Actualizar los componentes con los datos de InfluxDB
	updatedCount := 0
	for i, component := range components {
		hopperData, exists := influxData[component.DosingUnit]
		if !exists {
			continue
		}

		value, exists := hopperData[component.DosingHopper]
		if !exists {
			continue
		}

		components[i].Quantity = float32(value)
		updatedCount++
	}

	fmt.Printf("✅ Completado para fábrica %s: %d/%d componentes actualizados usando conexión %s\n",
		factory, updatedCount, len(components), connectionName)

	logger.Info("Update completed for factory",
		zap.String("factory", factory),
		zap.String("connection_used", connectionName),
		zap.Int("updatedCount", updatedCount),
		zap.Int("totalComponents", len(components)))

	return components, nil
}
