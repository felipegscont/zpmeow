package errors

import (
	"strings"
)

// IsRetryableError checks if an error indicates a retryable condition
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}
	
	errStr := err.Error()
	retryableErrors := []string{
		"connection refused",
		"timeout",
		"temporary failure",
		"network unreachable",
	}
	
	for _, retryable := range retryableErrors {
		if strings.Contains(strings.ToLower(errStr), retryable) {
			return true
		}
	}
	
	return false
}
