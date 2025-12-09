package service_health_aggregator

import "fmt"

type URLValidationError struct {
	URL    string
	Detail string
}

func (e URLValidationError) Error() string {
	return fmt.Sprintf("the URL %s is invalid. Details: %s", e.URL, e.Detail)
}
