package middleware

import (
	"zpmeow/internal/infra/logger"

	"github.com/gin-gonic/gin"
)

// Logger returns a gin.HandlerFunc (middleware) that logs requests using our logger
func Logger() gin.HandlerFunc {
	log := logger.GetLogger().Sub("http")

	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Log using our structured logger
		log.With().
			Str("method", param.Method).
			Str("path", param.Path).
			Int("status", param.StatusCode).
			Dur("latency", param.Latency).
			Str("client_ip", param.ClientIP).
			Time("timestamp", param.TimeStamp).
			Logger().
			Info("HTTP Request")

		// Return empty string since we're handling logging ourselves
		return ""
	})
}
