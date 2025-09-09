package router

import (
	"zpmeow/internal/health"
	"zpmeow/internal/session"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configures the application routes
func SetupRoutes(router *gin.Engine, waService *session.WhatsAppService, db *sqlx.DB) {
	// Create handlers
	sessionHandler := session.NewSessionHandler(waService, db)

	// Ping route for health checks
	router.GET("/ping", health.Ping)

	// Session routes
	sessionGroup := router.Group("/sessions")
	{
		sessionGroup.POST("/create", sessionHandler.CreateSession)
		sessionGroup.GET("/list", sessionHandler.ListSessions)
		sessionGroup.GET("/:id/info", sessionHandler.GetSessionInfo)
		sessionGroup.DELETE("/:id/delete", sessionHandler.DeleteSession)
		sessionGroup.POST("/:id/connect", sessionHandler.ConnectSession)
		sessionGroup.POST("/:id/logout", sessionHandler.LogoutSession)
		sessionGroup.GET("/:id/qr", sessionHandler.GetSessionQR)
		sessionGroup.POST("/:id/pair", sessionHandler.PairSession)
		sessionGroup.POST("/:id/proxy/set", sessionHandler.SetProxy)
		sessionGroup.GET("/:id/proxy/find", sessionHandler.GetProxy)
	}

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

