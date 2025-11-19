package databases

import (
	"fmt"
	"net/http"
	"time"

	"github.com/remrafvil/Auriga_API/config"
	"go.uber.org/zap"
)

type WorkeraClient struct {
	BaseURL    string
	ApiUser    string
	ApiKey     string
	httpClient *http.Client
	logger     *zap.Logger
}

func NewWorkeraClient(settings *config.Settings, logger *zap.Logger) *WorkeraClient {
	return &WorkeraClient{
		BaseURL: settings.Workera.Url,
		ApiUser: settings.Workera.API_User,
		ApiKey:  settings.Workera.API_Key,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// Get realiza una petición GET con autenticación
func (c *WorkeraClient) Get(endpoint string) (*http.Response, error) {
	url := c.BaseURL + endpoint
	c.logger.Debug("Making authenticated HTTP request to Workera API",
		zap.String("url", url),
		zap.String("endpoint", endpoint))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Agregar headers de autenticación
	req.Header.Set("API_USER", c.ApiUser)
	req.Header.Set("API_KEY", c.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return c.httpClient.Do(req)
}

// GetEmployees construye la URL específica para empleados
func (c *WorkeraClient) GetEmployees(page int) (*http.Response, error) {
	endpoint := fmt.Sprintf("employee?page=%d", page)
	return c.Get(endpoint)
}

// Método genérico para otros endpoints
func (c *WorkeraClient) GetEndpoint(endpoint string) (*http.Response, error) {
	return c.Get(endpoint)
}
