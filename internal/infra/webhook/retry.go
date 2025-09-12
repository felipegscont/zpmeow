package webhook

import (
	"context"
	"fmt"
	"time"

	"zpmeow/internal/infra/logger"
)

// RetryConfig defines retry configuration
type RetryConfig struct {
	MaxRetries      int
	InitialBackoff  time.Duration
	MaxBackoff      time.Duration
	BackoffMultiplier float64
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:        3,
		InitialBackoff:    1 * time.Second,
		MaxBackoff:        30 * time.Second,
		BackoffMultiplier: 2.0,
	}
}

// RetryStrategy handles retry logic for webhook operations
type RetryStrategy struct {
	config *RetryConfig
	logger logger.Logger
}

// NewRetryStrategy creates a new retry strategy
func NewRetryStrategy(config *RetryConfig) *RetryStrategy {
	if config == nil {
		config = DefaultRetryConfig()
	}
	
	return &RetryStrategy{
		config: config,
		logger: logger.GetLogger().Sub("webhook-retry"),
	}
}

// ExecuteWithRetry executes a function with retry logic
func (r *RetryStrategy) ExecuteWithRetry(ctx context.Context, operation func() error, operationName string) error {
	var lastErr error

	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		if attempt > 0 {
			backoff := r.calculateBackoff(attempt)
			r.logger.Infof("Retrying %s in %v (attempt %d/%d)", operationName, backoff, attempt, r.config.MaxRetries)
			
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		err := operation()
		if err == nil {
			if attempt > 0 {
				r.logger.Infof("%s succeeded after %d retries", operationName, attempt)
			}
			return nil
		}

		lastErr = err
		r.logger.Warnf("%s attempt %d failed: %v", operationName, attempt+1, err)
		
		// Check if error is retryable
		if !r.isRetryableError(err) {
			r.logger.Infof("%s failed with non-retryable error: %v", operationName, err)
			return err
		}
	}

	r.logger.Errorf("%s failed after %d attempts: %v", operationName, r.config.MaxRetries+1, lastErr)
	return fmt.Errorf("%s failed after %d attempts: %w", operationName, r.config.MaxRetries+1, lastErr)
}

// calculateBackoff calculates the backoff duration for a given attempt
func (r *RetryStrategy) calculateBackoff(attempt int) time.Duration {
	backoff := time.Duration(float64(r.config.InitialBackoff) * 
		pow(r.config.BackoffMultiplier, float64(attempt-1)))
	
	if backoff > r.config.MaxBackoff {
		backoff = r.config.MaxBackoff
	}
	
	return backoff
}

// isRetryableError checks if an error should trigger a retry
func (r *RetryStrategy) isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	
	// Check for common retryable error patterns
	errStr := err.Error()
	retryablePatterns := []string{
		"connection refused",
		"timeout",
		"temporary failure",
		"network unreachable",
		"status 5", // 5xx HTTP errors
		"context deadline exceeded",
		"no such host",
		"connection reset",
	}
	
	for _, pattern := range retryablePatterns {
		if contains(errStr, pattern) {
			return true
		}
	}
	
	return false
}

// pow calculates base^exp for float64 (simple implementation)
func pow(base, exp float64) float64 {
	if exp == 0 {
		return 1
	}
	if exp == 1 {
		return base
	}
	
	result := base
	for i := 1; i < int(exp); i++ {
		result *= base
	}
	return result
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return indexOf(s, substr) >= 0
}

// indexOf returns the index of substr in s, or -1 if not found
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
