package service_health_aggregator

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	ServiceStatusHealthy  = "healthy"
	ServiceStatusDown     = "down"
	ServiceStatusDegraded = "degraded"
)

type HealthService struct {
	config *Config
}

type AggregatedHealthResponse struct {
	Status    string          `json:"status"`
	Timestamp string          `json:"timestamp"`
	Services  []ServiceHealth `json:"services"`
}

type ServiceHealth struct {
	Name         string  `json:"name"`
	Status       string  `json:"status"`
	ResponseTime float64 `json:"response_time_ms,omitempty"`
	Error        string  `json:"error,omitempty"`
}

func NewHealthService(cfg *Config) *HealthService {
	return &HealthService{
		config: cfg,
	}
}

func (s *HealthService) CheckHealth(ctx context.Context) (*AggregatedHealthResponse, error) {
	log.Default().Println("Health Check Called")
	var healthResults AggregatedHealthResponse
	healthResults.Timestamp = time.Now().Format(time.RFC3339)

	for _, service := range s.config.Services {
		// Ensure valid URL
		urlErr := isValidURL(service.Url)
		if urlErr != nil {
			return nil, urlErr
		}
		timeout := time.Duration(service.Timeout) * time.Millisecond
		start := time.Now()
		req, err := http.NewRequestWithContext(ctx, "GET", service.Url, nil)
		if err != nil {
			return nil, err
		}

		client := &http.Client{Timeout: timeout}
		resp, respErr := client.Do(req)
		if respErr != nil {
			return nil, respErr
		}
		duration := time.Since(start)

		defer resp.Body.Close()

		var serviceResp ServiceHealth
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			serviceResp.Name = service.Name
			serviceResp.Status = ServiceStatusHealthy
			serviceResp.ResponseTime = float64(duration.Milliseconds())
		} else {
			serviceResp.Name = service.Name
			serviceResp.Status = ServiceStatusDown
			serviceResp.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
		}

		healthResults.Services = append(healthResults.Services, serviceResp)
	}

	healthResults.Status = DetermineOverallStatus(healthResults.Services)

	return &healthResults, nil
}

func isValidURL(u string) error {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return &URLValidationError{URL: u, Detail: err.Error()}
	}
	fmt.Printf("Parsed Host: %s", parsedURL.Host)
	if parsedURL.Scheme == "" {
		return &URLValidationError{URL: u, Detail: fmt.Sprintf("Scheme is empty")}
	}

	if parsedURL.Host == "" {
		return &URLValidationError{URL: u, Detail: fmt.Sprintf("Host is empty")}
	}
	return nil
}

func DetermineOverallStatus(services []ServiceHealth) string {
	if len(services) == 0 {
		return ServiceStatusDown
	}

	up := 0
	down := 0

	for _, s := range services {
		switch s.Status {
		case ServiceStatusHealthy:
			up++
		case ServiceStatusDown:
			down++
		default:
			down++
		}
	}

	switch {
	case up == len(services):
		return ServiceStatusHealthy
	case down == len(services):
		return ServiceStatusDown
	default:
		return ServiceStatusDegraded
	}
}
