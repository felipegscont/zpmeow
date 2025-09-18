package middleware

import (
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Metrics holds application metrics
type Metrics struct {
	mu                sync.RWMutex
	RequestCount      map[string]int64  `json:"request_count"`
	ResponseTime      map[string]int64  `json:"avg_response_time_ms"`
	StatusCodeCount   map[int]int64     `json:"status_code_count"`
	ActiveConnections int64             `json:"active_connections"`
	TotalRequests     int64             `json:"total_requests"`
	ErrorCount        int64             `json:"error_count"`
	LastReset         time.Time         `json:"last_reset"`
}

var globalMetrics = &Metrics{
	RequestCount:    make(map[string]int64),
	ResponseTime:    make(map[string]int64),
	StatusCodeCount: make(map[int]int64),
	LastReset:       time.Now(),
}

// GetMetrics returns a copy of current metrics
func GetMetrics() *Metrics {
	globalMetrics.mu.RLock()
	defer globalMetrics.mu.RUnlock()

	// Create a copy to avoid race conditions
	metrics := &Metrics{
		RequestCount:      make(map[string]int64),
		ResponseTime:      make(map[string]int64),
		StatusCodeCount:   make(map[int]int64),
		ActiveConnections: globalMetrics.ActiveConnections,
		TotalRequests:     globalMetrics.TotalRequests,
		ErrorCount:        globalMetrics.ErrorCount,
		LastReset:         globalMetrics.LastReset,
	}

	for k, v := range globalMetrics.RequestCount {
		metrics.RequestCount[k] = v
	}
	for k, v := range globalMetrics.ResponseTime {
		metrics.ResponseTime[k] = v
	}
	for k, v := range globalMetrics.StatusCodeCount {
		metrics.StatusCodeCount[k] = v
	}

	return metrics
}

// ResetMetrics resets all metrics
func ResetMetrics() {
	globalMetrics.mu.Lock()
	defer globalMetrics.mu.Unlock()

	globalMetrics.RequestCount = make(map[string]int64)
	globalMetrics.ResponseTime = make(map[string]int64)
	globalMetrics.StatusCodeCount = make(map[int]int64)
	globalMetrics.ActiveConnections = 0
	globalMetrics.TotalRequests = 0
	globalMetrics.ErrorCount = 0
	globalMetrics.LastReset = time.Now()
}

// MetricsMiddleware collects request metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Increment active connections
		globalMetrics.mu.Lock()
		globalMetrics.ActiveConnections++
		globalMetrics.TotalRequests++
		globalMetrics.mu.Unlock()

		// Process request
		c.Next()

		// Calculate metrics
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		endpoint := method + " " + path

		globalMetrics.mu.Lock()
		defer globalMetrics.mu.Unlock()

		// Decrement active connections
		globalMetrics.ActiveConnections--

		// Update request count
		globalMetrics.RequestCount[endpoint]++

		// Update response time (moving average)
		currentAvg := globalMetrics.ResponseTime[endpoint]
		requestCount := globalMetrics.RequestCount[endpoint]
		newAvg := (currentAvg*(requestCount-1) + duration.Milliseconds()) / requestCount
		globalMetrics.ResponseTime[endpoint] = newAvg

		// Update status code count
		globalMetrics.StatusCodeCount[statusCode]++

		// Update error count
		if statusCode >= 400 {
			globalMetrics.ErrorCount++
		}

		// Add metrics headers
		c.Header("X-Response-Time", duration.String())
		c.Header("X-Request-ID", getRequestID(c))
	}
}

// getRequestID gets request ID from context
func getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// PerformanceMiddleware adds performance monitoring
func PerformanceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Add performance headers
		c.Header("X-Server-Timing", "start;desc=\"Request Start\"")

		c.Next()

		// Calculate and add timing information
		duration := time.Since(start)
		c.Header("X-Response-Time", duration.String())
		c.Header("X-Response-Time-Ms", strconv.FormatInt(duration.Milliseconds(), 10))

		// Add server timing header
		serverTiming := "total;dur=" + strconv.FormatFloat(float64(duration.Nanoseconds())/1000000, 'f', 2, 64)
		c.Header("Server-Timing", serverTiming)
	}
}

// HealthMetricsMiddleware provides health check metrics
func HealthMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip metrics collection for health endpoints to avoid noise
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/ping" {
			c.Next()
			return
		}

		// Continue with normal metrics collection
		MetricsMiddleware()(c)
	}
}

// CircuitBreakerMiddleware implements basic circuit breaker pattern
func CircuitBreakerMiddleware(errorThreshold int, timeWindow time.Duration) gin.HandlerFunc {
	var (
		errorCount    int64
		lastErrorTime time.Time
		isOpen        bool
		mu            sync.RWMutex
	)

	return func(c *gin.Context) {
		mu.RLock()
		if isOpen && time.Since(lastErrorTime) < timeWindow {
			mu.RUnlock()
			c.AbortWithStatusJSON(503, gin.H{
				"error":   "Service temporarily unavailable",
				"message": "Circuit breaker is open due to high error rate",
				"code":    503,
			})
			return
		}
		mu.RUnlock()

		c.Next()

		// Check response status
		if c.Writer.Status() >= 500 {
			mu.Lock()
			errorCount++
			lastErrorTime = time.Now()

			if errorCount >= int64(errorThreshold) {
				isOpen = true
			}
			mu.Unlock()
		} else {
			// Reset on successful request
			mu.Lock()
			if time.Since(lastErrorTime) > timeWindow {
				errorCount = 0
				isOpen = false
			}
			mu.Unlock()
		}
	}
}

// RequestTrackingMiddleware tracks request patterns
func RequestTrackingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add request tracking headers
		c.Header("X-Request-Start", strconv.FormatInt(time.Now().UnixNano(), 10))
		c.Header("X-Request-Method", c.Request.Method)
		c.Header("X-Request-Path", c.Request.URL.Path)

		c.Next()

		// Add response tracking headers
		c.Header("X-Request-End", strconv.FormatInt(time.Now().UnixNano(), 10))
		c.Header("X-Response-Status", strconv.Itoa(c.Writer.Status()))
	}
}

// SlowRequestMiddleware logs slow requests
func SlowRequestMiddleware(threshold time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		if duration > threshold {
			// Log slow request (you can integrate with your logging system)
			c.Header("X-Slow-Request", "true")
			c.Header("X-Slow-Request-Duration", duration.String())
		}
	}
}
