//	@title			meow meow API
//	@version		1.0
//	@description	A meow API server built with Go, inspired by meow
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	meow API Support
//	@contact.url	https://github.com/your-username/meow
//	@contact.email	support@meow.com

//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT

//	@host		localhost:8080
//	@BasePath	/

//	@schemes	http https

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				API Key authentication. Simply provide your API key directly: "YOUR_API_KEY". The system automatically detects if it's a Global API Key (can access all sessions and session management) or a Session-specific API Key (can only access the specific session it belongs to).
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "meow/docs" // Import for swagger docs
	"meow/internal/application"
	"meow/internal/config"
	"meow/internal/domain/session"
	"meow/internal/infra/database"
	"meow/internal/infra/database/repository"
	"meow/internal/infra/logging"
	"meow/internal/infra/middleware"
	"meow/internal/infra/webhooks"
	"meow/internal/infra/wmeow"
	"meow/internal/interfaces/handlers"
	"meow/internal/interfaces/routes"

	"github.com/gin-gonic/gin"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {

		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	log := logging.Initialize(cfg.GetLogging())
	logging.SetLogger(log)
	log.Info("Starting meow server")

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(cfg); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	dbLog := logging.GetWALogger("Database")
	_, err = sqlstore.New(context.Background(), "postgres", cfg.GetDatabaseURL(), dbLog)
	if err != nil {
		log.Fatalf("Failed to create whatsmeow container: %v", err)
	}

	// Create session repository
	sessionRepo := repository.NewPostgresRepo(db)

	// Create webhook service
	_ = webhooks.NewService()

	// Create WhatsApp service (infrastructure implementation)
	wmeowService := wmeow.NewService()

	// Create domain service
	domainService := session.NewService()

	// Create application services with proper dependencies
	appSessionService := application.NewSessionApp(sessionRepo, domainService)
	webhookAppService := application.NewWebhookApp(sessionRepo)

	log.Info("WhatsApp service initialized")

	log.Info("Session service initialized")

	// Create authentication middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg, sessionRepo, log)

	// Use application service in handlers
	sessionHandler := handlers.NewSessionHandler(appSessionService, wmeowService)
	healthHandler := handlers.NewHealthHandler(db)
	messageHandler := handlers.NewMessageHandler(appSessionService, wmeowService)
	chatHandler := handlers.NewChatHandler(appSessionService, wmeowService)
	groupHandler := handlers.NewGroupHandler(appSessionService, wmeowService)
	communityHandler := handlers.NewCommunityHandler(appSessionService, wmeowService)
	webhookHandler := handlers.NewWebhookHandler(appSessionService, webhookAppService)
	contactHandler := handlers.NewContactHandler(appSessionService, wmeowService)
	newsletterHandler := handlers.NewNewsletterHandler(appSessionService, wmeowService)
	privacyHandler := handlers.NewPrivacyHandler(appSessionService, wmeowService)

	gin.SetMode(cfg.GetServer().GetMode())

	// Create Gin router with custom middleware (no default logging)
	ginRouter := gin.New()

	// Add security and performance middlewares
	ginRouter.Use(middleware.SecurityHeadersMiddleware())
	ginRouter.Use(middleware.ContentSecurityMiddleware())
	ginRouter.Use(middleware.RequestIDMiddleware())
	ginRouter.Use(middleware.RecoveryMiddleware())
	ginRouter.Use(middleware.TimeoutMiddleware(25 * time.Second)) // Slightly less than server timeout
	ginRouter.Use(middleware.RateLimitMiddleware(100))            // 100 requests per minute per IP
	ginRouter.Use(middleware.RequestValidationMiddleware())
	ginRouter.Use(middleware.APIVersionMiddleware())
	ginRouter.Use(middleware.MetricsMiddleware())
	ginRouter.Use(middleware.PerformanceMiddleware())
	ginRouter.Use(middleware.RequestSizeMiddleware(10 * 1024 * 1024)) // 10MB limit
	ginRouter.Use(middleware.CacheControlMiddleware())
	ginRouter.Use(middleware.CircuitBreakerMiddleware(10, 5*time.Minute)) // 10 errors in 5 minutes
	ginRouter.Use(middleware.SlowRequestMiddleware(2 * time.Second))      // Log requests > 2s
	ginRouter.Use(middleware.AuditMiddleware())                           // Audit important actions

	// Add custom logging middleware only for errors and important requests
	ginRouter.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/ping", "/health"}, // Skip health check logs
		Formatter: func(param gin.LogFormatterParams) string {
			// Only log non-2xx responses and important endpoints
			if param.StatusCode >= 400 ||
				(param.Path != "/ping" && param.Path != "/health") {
				return fmt.Sprintf("%s - %s %s %d %s\n",
					param.TimeStamp.Format("15:04:05"),
					param.Method,
					param.Path,
					param.StatusCode,
					param.Latency,
				)
			}
			return ""
		},
	}))

	// Add CORS middleware with configuration
	ginRouter.Use(middleware.CORS(cfg.GetCORS()))

	// Setup routes with handler dependencies
	handlerDeps := &routes.HandlerDependencies{
		HealthHandler:     healthHandler,
		SessionHandler:    sessionHandler,
		ContactHandler:    contactHandler,
		ChatHandler:       chatHandler,
		MessageHandler:    messageHandler,
		GroupHandler:      groupHandler,
		CommunityHandler:  communityHandler,
		NewsletterHandler: newsletterHandler,
		WebhookHandler:    webhookHandler,
		PrivacyHandler:    privacyHandler,
	}

	routes.SetupRoutes(ginRouter, handlerDeps, authMiddleware)

	// Create HTTP server with timeouts
	addr := fmt.Sprintf(":%s", cfg.GetServer().GetPort())
	srv := &http.Server{
		Addr:         addr,
		Handler:      ginRouter,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Infof("Server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exited")
}
