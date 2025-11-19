package rwWorkera

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/remrafvil/Auriga_API/databases"
	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"

	"go.uber.org/zap"
)

type Repository interface {
	GetEmployees(page int) (*rModels.EmployeeResponse, error)
	GetCustomEndpoint(endpoint string) ([]byte, error)
	GetAllEmployees() ([]rModels.Employee, error)
}

type repository struct {
	client *databases.WorkeraClient
	logger *zap.Logger
}

func New(client *databases.WorkeraClient, logger *zap.Logger) Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}

func (r *repository) GetEmployees(page int) (*rModels.EmployeeResponse, error) {
	r.logger.Info("Fetching employees from Workera API", zap.Int("page", page))

	resp, err := r.client.GetEmployees(page)
	if err != nil {
		r.logger.Error("Failed to fetch employees from Workera API",
			zap.Int("page", page),
			zap.Error(err))
		return nil, fmt.Errorf("failed to call Workera API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		r.logger.Error("Workera API returned non-200 status",
			zap.Int("status", resp.StatusCode),
			zap.Int("page", page))
		return nil, fmt.Errorf("Workera API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.logger.Error("Failed to read response body from Workera API", zap.Error(err))
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var employeeResponse rModels.EmployeeResponse
	if err := json.Unmarshal(body, &employeeResponse); err != nil {
		r.logger.Error("Failed to unmarshal response from Workera API",
			zap.Error(err),
			zap.String("response", string(body)))
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	r.logger.Info("Successfully fetched employees from Workera API",
		zap.Int("page", page),
		zap.Int("employeesCount", len(employeeResponse.Data)),
		zap.Int("totalPages", employeeResponse.TotalPages))

	return &employeeResponse, nil
}

// Método para consultar otros endpoints de la API
func (r *repository) GetCustomEndpoint(endpoint string) ([]byte, error) {
	r.logger.Info("Fetching custom endpoint from Workera API", zap.String("endpoint", endpoint))

	resp, err := r.client.GetEndpoint(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to call endpoint %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("endpoint %s returned status %d", endpoint, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from endpoint %s: %w", endpoint, err)
	}

	return body, nil
}

func (r *repository) GetAllEmployees() ([]rModels.Employee, error) {
	var allEmployees []rModels.Employee

	r.logger.Info("Starting to fetch all employees from Workera API")

	// Get first page to know total pages
	firstPage, err := r.GetEmployees(1)
	if err != nil {
		return nil, fmt.Errorf("failed to get first page: %w", err)
	}

	allEmployees = append(allEmployees, firstPage.Data...)
	totalPages := firstPage.TotalPages

	r.logger.Info("Starting to fetch all pages from Workera API",
		zap.Int("totalPages", totalPages))

	// Fetch remaining pages concurrentemente
	type result struct {
		page      int
		employees []rModels.Employee
		err       error
	}

	results := make(chan result, totalPages-1)

	for page := 2; page <= totalPages; page++ {
		go func(p int) {
			response, err := r.GetEmployees(p)
			if err != nil {
				results <- result{page: p, err: err}
				return
			}
			results <- result{page: p, employees: response.Data}
		}(page)
	}

	// Collect results
	for i := 2; i <= totalPages; i++ {
		res := <-results
		if res.err != nil {
			r.logger.Error("Failed to fetch page from Workera API",
				zap.Int("page", res.page),
				zap.Error(res.err))
			continue
		}
		allEmployees = append(allEmployees, res.employees...)
		r.logger.Debug("Added employees from page",
			zap.Int("page", res.page),
			zap.Int("count", len(res.employees)))
	}

	r.logger.Info("Completed fetching all employees from Workera API",
		zap.Int("totalEmployees", len(allEmployees)))
	return allEmployees, nil
}
