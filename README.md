# Service Health Aggregator – N-able

## Overview

**Service Health Aggregator** is a lightweight Go-based service that aggregates the health status of multiple downstream services.  
It exposes a single HTTP endpoint that checks configured service URLs and returns an aggregated health response.

The application is designed to be simple, containerized, and easy to run locally using Docker.

---
## Running Locally (Docker)

### 1. Build the Docker image

```bash
docker build -t svc-health-aggregator:latest .
```
### 2. Run the container

```bash
docker run -p 8000:8000 svc-health-aggregator:latest
```
The service will be available at: `http://localhost:8000/health/aggregate`

---

## API Endpoint

### `GET /health/aggregate`

This endpoint:
- Reads service definitions from a YAML configuration file
- Sends health check requests to each configured service URL
- Aggregates individual service health into an overall status
- Returns a JSON response with:
    - Overall system status
    - Timestamp
    - Per-service health details

### Health Status Rules

- **healthy** → All services return `2xx` responses within timeout
- **degraded** → Some services are down
- **down** → All services are down

---

## Example Response

```json
{
  "status": "degraded",
  "timestamp": "2025-11-26T10:30:00Z",
  "services": [
    {
      "name": "api-gateway",
      "status": "healthy",
      "response_time_ms": 45
    },
    {
      "name": "user-service",
      "status": "down",
      "error": "HTTP 500"
    }
  ]
}
