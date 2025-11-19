package databases

import (
	"context"
	"fmt"
	"sync"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/remrafvil/Auriga_API/config"
	"go.uber.org/zap"
)

// InfluxClientManager maneja múltiples conexiones a InfluxDB
type InfluxClientManager struct {
	clients      map[string]influxdb2.Client
	configs      map[string]config.InfluxDBConfig
	logger       *zap.Logger
	mu           sync.RWMutex
	healthStatus map[string]bool
}

// ConnectionResult representa el resultado de una conexión
type ConnectionResult struct {
	Name   string
	Client influxdb2.Client
	Error  error
	Config config.InfluxDBConfig
}

// ConnectionStatus representa el estado de una conexión
type ConnectionStatus struct {
	Healthy   bool      `json:"healthy"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Org       string    `json:"org"`
	Bucket    string    `json:"bucket"`
	LastCheck time.Time `json:"last_check"`
}

// NewInfluxClientManager crea un nuevo manager de clientes InfluxDB
func NewInfluxClientManager(settings *config.Settings, logger *zap.Logger) *InfluxClientManager {
	manager := &InfluxClientManager{
		clients:      make(map[string]influxdb2.Client),
		configs:      make(map[string]config.InfluxDBConfig),
		healthStatus: make(map[string]bool),
		logger:       logger,
	}

	// Cargar todas las configuraciones
	for name, cfg := range settings.InfluxDBs {
		manager.configs[name] = cfg
	}

	// Inicializar conexiones automáticamente
	manager.initializeAllConnections(context.Background())

	return manager
}

// initializeAllConnections inicializa todas las conexiones
func (m *InfluxClientManager) initializeAllConnections(ctx context.Context) {
	m.logger.Info("Inicializando conexiones InfluxDB",
		zap.Int("total_configuraciones", len(m.configs)))

	var wg sync.WaitGroup
	results := make(chan ConnectionResult, len(m.configs))

	// Inicializar conexiones concurrentemente
	for name, cfg := range m.configs {
		wg.Add(1)
		go func(name string, cfg config.InfluxDBConfig) {
			defer wg.Done()
			m.initializeConnection(ctx, name, cfg, results)
		}(name, cfg)
	}

	// Cerrar el channel cuando todas las goroutines terminen
	go func() {
		wg.Wait()
		close(results)
	}()

	// Procesar resultados y actualizar el estado
	successfulConnections := 0
	var availableConnections []string

	for result := range results {
		m.mu.Lock()
		if result.Error == nil {
			m.clients[result.Name] = result.Client
			m.healthStatus[result.Name] = true
			successfulConnections++
			availableConnections = append(availableConnections, result.Name)
			m.logger.Info("✅ Conexión InfluxDB establecida",
				zap.String("nombre", result.Name),
				zap.String("host", result.Config.Host),
				zap.String("org", result.Config.Org),
			)
		} else {
			m.healthStatus[result.Name] = false
			m.logger.Warn("❌ No se pudo establecer conexión InfluxDB",
				zap.String("nombre", result.Name),
				zap.String("host", result.Config.Host),
				zap.Error(result.Error),
			)
		}
		m.mu.Unlock()
	}

	// Log del estado final
	total := len(m.configs)
	m.logger.Info("Estado final de conexiones InfluxDB",
		zap.Int("configuraciones_totales", total),
		zap.Int("conexiones_exitosas", successfulConnections),
		zap.Int("conexiones_fallidas", total-successfulConnections),
		zap.Strings("conexiones_disponibles", availableConnections),
	)
}

// initializeConnection intenta establecer una conexión individual
func (m *InfluxClientManager) initializeConnection(ctx context.Context, name string, cfg config.InfluxDBConfig, results chan<- ConnectionResult) {
	influxURL := fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)

	m.logger.Debug("Conectando a InfluxDB",
		zap.String("nombre", name),
		zap.String("url", influxURL),
	)

	client := influxdb2.NewClient(influxURL, cfg.Token)

	// Verificar conexión con timeout
	healthCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	health, err := client.Health(healthCtx)
	if err != nil {
		client.Close()
		results <- ConnectionResult{Name: name, Error: fmt.Errorf("health check failed: %w", err)}
		return
	}

	if health.Status != "pass" {
		client.Close()
		results <- ConnectionResult{Name: name, Error: fmt.Errorf("unhealthy status: %s", health.Status)}
		return
	}

	results <- ConnectionResult{Name: name, Client: client, Config: cfg}
}

// GetClient obtiene un cliente por nombre
func (m *InfluxClientManager) GetClient(name string) (influxdb2.Client, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	client, exists := m.clients[name]
	return client, exists
}

// GetAvailableClients retorna los nombres de los clientes disponibles
func (m *InfluxClientManager) GetAvailableClients() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	available := make([]string, 0, len(m.clients))
	for name := range m.clients {
		available = append(available, name)
	}
	return available
}

// GetConfig retorna la configuración para un cliente específico
func (m *InfluxClientManager) GetConfig(name string) (config.InfluxDBConfig, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cfg, exists := m.configs[name]
	return cfg, exists
}

// GetAllConfigs retorna todas las configuraciones
func (m *InfluxClientManager) GetAllConfigs() map[string]config.InfluxDBConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	configs := make(map[string]config.InfluxDBConfig)
	for k, v := range m.configs {
		configs[k] = v
	}
	return configs
}

// IsHealthy verifica si una conexión específica está saludable EN TIEMPO REAL
func (m *InfluxClientManager) IsHealthy(name string) bool {
	m.mu.RLock()
	client, exists := m.clients[name]
	oldStatus := m.healthStatus[name]
	m.mu.RUnlock()

	if !exists || client == nil {
		m.updateHealthStatus(name, false, oldStatus)
		return false
	}

	// Verificación de salud en tiempo real con timeout corto
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	health, err := client.Health(ctx)
	if err != nil {
		m.updateHealthStatus(name, false, oldStatus)
		return false
	}

	// Si health es nil, considerar no saludable
	if health == nil {
		m.updateHealthStatus(name, false, oldStatus)
		return false
	}

	healthy := health.Status == "pass"
	m.updateHealthStatus(name, healthy, oldStatus)

	return healthy
}

// updateHealthStatus actualiza el estado de salud y detecta cambios
func (m *InfluxClientManager) updateHealthStatus(name string, healthy bool, oldStatus bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.healthStatus[name] = healthy

	// Detectar cambio de estado y loggear
	if oldStatus != healthy {
		if healthy {
			m.logger.Info("🟢 CONEXIÓN RECUPERADA",
				zap.String("connection", name),
				zap.String("host", m.configs[name].Host),
				zap.Int("port", m.configs[name].Port),
			)
		} else {
			m.logger.Error("🔴 CONEXIÓN PERDIDA",
				zap.String("connection", name),
				zap.String("host", m.configs[name].Host),
				zap.Int("port", m.configs[name].Port),
			)
		}
	}
}

// StartConnectionMonitor inicia el monitoreo continuo de conexiones
func (m *InfluxClientManager) StartConnectionMonitor(ctx context.Context) {
	m.logger.Info("🚀 Iniciando monitoreo de conexiones InfluxDB")

	ticker := time.NewTicker(30 * time.Second) // Verificar cada 30 segundos
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			m.logger.Info("🛑 Deteniendo monitoreo de conexiones")
			return
		case <-ticker.C:
			m.monitorAllConnections()
		}
	}
}

// monitorAllConnections verifica el estado de todas las conexiones
func (m *InfluxClientManager) monitorAllConnections() {
	m.mu.RLock()
	connections := make([]string, 0, len(m.configs))
	for name := range m.configs {
		connections = append(connections, name)
	}
	m.mu.RUnlock()

	for _, name := range connections {
		m.IsHealthy(name) // Esto automáticamente actualiza el estado y loggea cambios
	}

	// Log resumen periódico (cada 5 minutos)
	if time.Now().Minute()%5 == 0 { // Cada 5 minutos
		m.logConnectionSummary()
	}
}

// logConnectionSummary muestra un resumen del estado de todas las conexiones
func (m *InfluxClientManager) logConnectionSummary() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	total := len(m.configs)
	healthy := 0
	var healthyConns, unhealthyConns []string

	for name, status := range m.healthStatus {
		if status {
			healthy++
			healthyConns = append(healthyConns, name)
		} else {
			unhealthyConns = append(unhealthyConns, name)
		}
	}

	m.logger.Info("📊 RESUMEN ESTADO CONEXIONES",
		zap.Int("total", total),
		zap.Int("saludables", healthy),
		zap.Int("fallidas", total-healthy),
		zap.Strings("conexiones_activas", healthyConns),
		zap.Strings("conexiones_inactivas", unhealthyConns),
	)
}

// GetConnectionStatus retorna el estado detallado de todas las conexiones
func (m *InfluxClientManager) GetConnectionStatus() map[string]ConnectionStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status := make(map[string]ConnectionStatus)
	for name, healthy := range m.healthStatus {
		cfg := m.configs[name]
		status[name] = ConnectionStatus{
			Healthy:   healthy,
			Host:      cfg.Host,
			Port:      cfg.Port,
			Org:       cfg.Org,
			Bucket:    cfg.Bucket,
			LastCheck: time.Now(),
		}
	}
	return status
}

// CloseAll cierra todas las conexiones InfluxDB
func (m *InfluxClientManager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Info("🔌 Cerrando todas las conexiones InfluxDB")

	for name, client := range m.clients {
		m.logger.Debug("Cerrando conexión", zap.String("connection", name))
		client.Close()
	}

	// Limpiar los mapas
	m.clients = make(map[string]influxdb2.Client)
	m.healthStatus = make(map[string]bool)

	m.logger.Info("✅ Todas las conexiones InfluxDB cerradas")
}
