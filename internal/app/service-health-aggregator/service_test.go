package service_health_aggregator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDetermineOverallStatus(t *testing.T) {
	tests := []struct {
		name     string
		services []ServiceHealth
		expected string
	}{
		{
			name: "all healthy",
			services: []ServiceHealth{
				{Status: ServiceStatusHealthy},
				{Status: ServiceStatusHealthy},
			},
			expected: ServiceStatusHealthy,
		},
		{
			name: "all down",
			services: []ServiceHealth{
				{Status: ServiceStatusDown},
				{Status: ServiceStatusDown},
			},
			expected: ServiceStatusDown,
		},
		{
			name: "mixed results",
			services: []ServiceHealth{
				{Status: ServiceStatusHealthy},
				{Status: ServiceStatusDown},
			},
			expected: ServiceStatusDegraded,
		},
		{
			name:     "no services",
			services: []ServiceHealth{},
			expected: ServiceStatusDown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetermineOverallStatus(tt.services)
			if result != tt.expected {
				t.Fatalf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid https url", "https://example.com", false},
		{"valid http url", "http://localhost:8080", false},
		{"missing scheme", "example.com", true},
		{"empty url", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := isValidURL(tt.url)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("did not expect error, got %v", err)
			}
		})
	}
}

func TestCheckHealth_AllHealthy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := &Config{
		Services: []Service{
			{
				Name:    "service-a",
				Url:     server.URL,
				Timeout: 1000,
			},
			{
				Name:    "service-b",
				Url:     server.URL,
				Timeout: 1000,
			},
		},
	}

	svc := NewHealthService(cfg)

	resp, err := svc.CheckHealth(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Status != ServiceStatusHealthy {
		t.Fatalf("expected status %s, got %s", ServiceStatusHealthy, resp.Status)
	}

	if len(resp.Services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(resp.Services))
	}
}

func TestCheckHealth_Degraded(t *testing.T) {
	okServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer okServer.Close()

	failServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer failServer.Close()

	cfg := &Config{
		Services: []Service{
			{
				Name:    "healthy-service",
				Url:     okServer.URL,
				Timeout: 1000,
			},
			{
				Name:    "down-service",
				Url:     failServer.URL,
				Timeout: 1000,
			},
		},
	}

	svc := NewHealthService(cfg)

	resp, err := svc.CheckHealth(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Status != ServiceStatusDegraded {
		t.Fatalf("expected %s, got %s", ServiceStatusDegraded, resp.Status)
	}
}

func TestCheckHealth_InvalidURL(t *testing.T) {
	cfg := &Config{
		Services: []Service{
			{
				Name:    "bad-service",
				Url:     "not-a-valid-url",
				Timeout: 1000,
			},
		},
	}

	svc := NewHealthService(cfg)

	_, err := svc.CheckHealth(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
