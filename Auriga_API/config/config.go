package config

import (
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

//go:embed settings.yml
var settingsFile []byte

type App struct {
	EchoPort string `yaml:"port"`
	Zap      string `yaml:"zap"`
	Name     string `yaml:"name"`
	Version  string `yaml:"version"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Timezone string `yaml:"timezone"`
}
type GrafanaConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type InfluxDBConfig struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Token  string `yaml:"token"`
	Org    string `yaml:"org"`
	Bucket string `yaml:"bucket"`
}

type SapDBConfig struct {
	Url  string `yaml:"url"`
	Auth string `yaml:"auth"`
}

type AuthConfig struct {
	Issuer        string        `yaml:"issuer"`
	ClientID      string        `yaml:"client_id"`
	ClientSecret  string        `yaml:"client_secret"`
	RedirectURL   string        `yaml:"redirect_url"`
	JWKSEndpoint  string        `yaml:"jwks_endpoint"`
	JWKSCacheTTL  time.Duration `yaml:"jwks_cache_ttl"`
	TokenCacheTTL time.Duration `yaml:"token_cache_ttl"`
}
type WorkeraConfig struct {
	Url      string `yaml:"url"`
	API_User string `yaml:"API_USER"`
	API_Key  string `yaml:"API_KEY"`
}

type Settings struct {
	App        App                       `yaml:"app"`
	DB         DatabaseConfig            `yaml:"database"`
	Grafana1   GrafanaConfig             `yaml:"grafana1"`
	Grafana2   GrafanaConfig             `yaml:"grafana2"`
	InfluxDBs  map[string]InfluxDBConfig `yaml:"influxdbs"` // Múltiples conexiones
	Sap        SapDBConfig               `yaml:"sapdb"`
	AuthConfig AuthConfig                `yaml:"auth"`
	Workera    WorkeraConfig             `yaml:"workera"`
}

func New(logger *zap.Logger) (*Settings, error) {
	var s Settings

	err := yaml.Unmarshal(settingsFile, &s)
	if err != nil {
		logger.Error("Error al deserializar el archivo de configuración", zap.Error(err))
		return nil, err
	}

	// Log de configuración cargada
	logger.Info("Configuración cargada correctamente",
		zap.String("SAP URL", s.Sap.Url),
		zap.String("SAP Auth", s.Sap.Auth),
		zap.String("App Name", s.App.Name),
		zap.String("App Port", s.App.EchoPort),
	)

	// Log de información de InfluxDBs disponibles
	for name, influx := range s.InfluxDBs {
		logger.Debug("Configuración de InfluxDB",
			zap.String("nombre", name),
			zap.String("host", influx.Host),
			zap.Int("port", influx.Port),
			zap.String("org", influx.Org),
		)
	}

	// Authentik AuthConfig log
	logger.Info("Authentik AuthConfig cargada",
		zap.String("Issuer", s.AuthConfig.Issuer),
		zap.String("ClientID", s.AuthConfig.ClientID),
		zap.String("RedirectURL", s.AuthConfig.RedirectURL),
		zap.String("JWKSEndpoint", s.AuthConfig.JWKSEndpoint),
		zap.Duration("JWKSCacheTTL", s.AuthConfig.JWKSCacheTTL),
		zap.Duration("TokenCacheTTL", s.AuthConfig.TokenCacheTTL),
	)

	// Workera Config log
	logger.Info("Workera Config cargada",
		zap.String("URL", s.Workera.Url),
		zap.String("API_USER", s.Workera.API_User),
		zap.String("API_KEY", s.Workera.API_Key),
	)

	return &s, nil
}
func NewEnv() (*Settings, error) {
	var s_env Settings

	// Configuración de la aplicación
	s_env.App.EchoPort = os.Getenv("AURIGA_API_PORT")
	fmt.Printf("El puerto configurado es: %s\n", s_env.App.EchoPort)
	s_env.App.Zap = os.Getenv("AURIGA_ZAP")
	fmt.Printf("El logger está configurado en modo: %s\n", s_env.App.Zap)
	s_env.App.Name = os.Getenv("AURIGA_NAME")
	s_env.App.Version = os.Getenv("AURIGA_VERSION")
	fmt.Printf("El nombre de la app es: %s y la versión es: %s\n", s_env.App.Name, s_env.App.Version)

	// Configuración de base de datos
	s_env.DB.Host = os.Getenv("DATABASE_HOST")
	s_env.DB.Port, _ = strconv.Atoi(os.Getenv("DATABASE_PORT"))
	s_env.DB.User = os.Getenv("DATABASE_USER")
	s_env.DB.Password = os.Getenv("DATABASE_PASSWORD")
	s_env.DB.Name = os.Getenv("DATABASE_NAME")
	s_env.DB.Timezone = os.Getenv("DATABASE_TIMEZONE")

	// Configuración de Grafana
	s_env.Grafana1.Host = os.Getenv("GRAFANA1_HOST")
	fmt.Printf("El Grafana1 Host es: %s\n", s_env.Grafana1.Host)
	s_env.Grafana1.Port, _ = strconv.Atoi(os.Getenv("GRAFANA1_PORT"))
	fmt.Printf("El Grafana1 Port es: %d\n", s_env.Grafana1.Port)
	s_env.Grafana2.Host = os.Getenv("GRAFANA2_HOST")
	fmt.Printf("El Grafana2 Host es: %s\n", s_env.Grafana2.Host)
	s_env.Grafana2.Port, _ = strconv.Atoi(os.Getenv("GRAFANA2_PORT"))
	fmt.Printf("El Grafana2 Port es: %d\n", s_env.Grafana2.Port)

	// Configuración de SAP
	s_env.Sap.Url = os.Getenv("SAPDB_URL")
	s_env.Sap.Auth = os.Getenv("SAPDB_AUTH")
	fmt.Printf("El Token SAP es: %s\n", s_env.Sap.Auth)

	// Cargar múltiples configuraciones de InfluxDB
	s_env.InfluxDBs = loadMultipleInfluxDBConfigs()

	// Configuración de Authentik
	s_env.AuthConfig.Issuer = os.Getenv("AUTHENTIK_ISSUER")
	s_env.AuthConfig.ClientID = os.Getenv("AUTHENTIK_CLIENT_ID")
	s_env.AuthConfig.ClientSecret = os.Getenv("AUTHENTIK_CLIENT_SECRET")
	s_env.AuthConfig.RedirectURL = os.Getenv("AUTHENTIK_REDIRECT_URL")
	s_env.AuthConfig.JWKSEndpoint = os.Getenv("AUTHENTIK_JWKS_ENDPOINT")
	s_env.AuthConfig.JWKSCacheTTL = time.Hour        // Valor temporal, se sobrescribirá más abajo
	s_env.AuthConfig.TokenCacheTTL = 5 * time.Minute // Valor temporal, se sobrescribirá más abajo

	fmt.Printf("Authentik Issuer: %s\n", s_env.AuthConfig.Issuer)
	fmt.Printf("Authentik ClientID: %s\n", s_env.AuthConfig.ClientID)
	fmt.Printf("Authentik ClientSecret: %s\n", s_env.AuthConfig.ClientSecret)
	fmt.Printf("Authentik RedirectURL: %s\n", s_env.AuthConfig.RedirectURL)
	fmt.Printf("Authentik JWKSEndpoint: %s\n", s_env.AuthConfig.JWKSEndpoint)
	fmt.Printf("JWKS Cache TTL: %s\n", s_env.AuthConfig.JWKSCacheTTL)
	fmt.Printf("Token Cache TTL: %s\n", s_env.AuthConfig.TokenCacheTTL)

	// Configuración de Workera
	s_env.Workera.Url = os.Getenv("WORKERA_URL")
	s_env.Workera.API_User = os.Getenv("WORKERA_API_USER")
	s_env.Workera.API_Key = os.Getenv("WORKERA_API_KEY")
	fmt.Printf("Workera URL: %s\n", s_env.Workera.Url)
	fmt.Printf("Workera API_USER: %s\n", s_env.Workera.API_User)
	// No imprimir API_KEY por seguridad

	return &s_env, nil
}

// loadMultipleInfluxDBConfigs carga todas las configuraciones de InfluxDB desde variables de entorno
func loadMultipleInfluxDBConfigs() map[string]InfluxDBConfig {
	configs := make(map[string]InfluxDBConfig)

	connections := []string{
		"CXB", "CXC", "CXD", "CXE", "CXF", "CXM",
		"EXT", "FPC", "FPB", "FSP", "MNT",
	}

	for _, connName := range connections {
		host := os.Getenv(fmt.Sprintf("%s_INFLUXDB_HOST", connName))
		if host == "" {
			continue
		}

		port, _ := strconv.Atoi(os.Getenv(fmt.Sprintf("%s_INFLUXDB_PORT", connName)))
		token := os.Getenv(fmt.Sprintf("%s_INFLUXDB_TOKEN", connName))
		org := os.Getenv(fmt.Sprintf("%s_INFLUXDB_ORG", connName))
		bucket := os.Getenv(fmt.Sprintf("%s_INFLUXDB_BUCKET", connName))

		// Mantener el nombre original en mayúsculas
		configs[connName] = InfluxDBConfig{
			Host:   host,
			Port:   port,
			Token:  token,
			Org:    org,
			Bucket: bucket,
		}

		fmt.Printf("Configurada conexión InfluxDB: %s - Host: %s:%d, Org: %s\n",
			connName, host, port, org)
	}

	fmt.Printf("Total de configuraciones InfluxDB cargadas: %d\n", len(configs))
	return configs
}

// Métodos de utilidad (mantener igual)
func (c *Settings) IsProduction() bool {
	return strings.ToLower(c.App.Zap) == "production"
}

func (c *Settings) IsDevelopment() bool {
	env := strings.ToLower(c.App.Zap)
	return env == "development" || env == "dev" || env == ""
}

func (c *Settings) IsStaging() bool {
	return strings.ToLower(c.App.Zap) == "staging"
}
