package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"zpmeow/internal/application"
	"zpmeow/internal/infra/middleware"
	"zpmeow/internal/infra/wmeow"
	"zpmeow/internal/interfaces/handlers"
	"zpmeow/internal/interfaces/routes"
)

// ServerConfig contains configuration for the HTTP server
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DefaultServerConfig returns default server configuration
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:         "8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

// Server represents the HTTP server
type Server struct {
	config     *ServerConfig
	router     *gin.Engine
	httpServer *http.Server
	handlers   *routes.HandlerDependencies
}

// NewServer creates a new HTTP server with all dependencies
func NewServer(
	config *ServerConfig,
	db *sqlx.DB,
	sessionApp *application.SessionApp,
	webhookApp *application.WebhookApp,
	wmeowService wmeow.WameowService,
	authMiddleware *middleware.AuthMiddleware,
) *Server {
	// Create handlers
	handlerDeps := &routes.HandlerDependencies{
		HealthHandler:     handlers.NewHealthHandler(db),
		SessionHandler:    handlers.NewSessionHandler(sessionApp, wmeowService),
		ContactHandler:    handlers.NewContactHandler(sessionApp, wmeowService),
		ChatHandler:       handlers.NewChatHandler(sessionApp, wmeowService),
		MessageHandler:    handlers.NewMessageHandler(sessionApp, wmeowService),
		GroupHandler:      handlers.NewGroupHandler(sessionApp, wmeowService),
		CommunityHandler:  handlers.NewCommunityHandler(sessionApp, wmeowService),
		NewsletterHandler: handlers.NewNewsletterHandler(sessionApp, wmeowService),
		WebhookHandler:    handlers.NewWebhookHandler(sessionApp, webhookApp),
	}

	// Create router
	router := gin.New()

	// Setup routes
	routes.SetupRoutes(router, handlerDeps, authMiddleware)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      router,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	return &Server{
		config:     config,
		router:     router,
		httpServer: httpServer,
		handlers:   handlerDeps,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	fmt.Printf("🚀 Starting HTTP server on port %s\n", s.config.Port)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	fmt.Println("🛑 Stopping HTTP server...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	fmt.Println("✅ HTTP server stopped")
	return nil
}

// GetRouter returns the gin router (useful for testing)
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// GetHandlers returns the handler dependencies (useful for testing)
func (s *Server) GetHandlers() *routes.HandlerDependencies {
	return s.handlers
}
