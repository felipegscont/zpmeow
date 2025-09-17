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

	_ "meow/docs" // Import for swagger docs
	"meow/internal/application"
	"meow/internal/config"
	"meow/internal/domain/session"
	"meow/internal/infrastructure/database"
	"meow/internal/infrastructure/database/repository"
	"meow/internal/infrastructure/logging"
	"meow/internal/infrastructure/middleware"
	"meow/internal/infrastructure/webhooks"
	"meow/internal/infrastructure/wmeow"
	"meow/internal/interfaces/http/handlers"
	"meow/internal/interfaces/http/routes"
	"meow/internal/shared/validation"

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
	container, err := sqlstore.New(context.Background(), "postgres", cfg.GetDatabaseURL(), dbLog)
	if err != nil {
		log.Fatalf("Failed to create whatsmeow container: %v", err)
	}

	waLogger := logging.GetWALogger("MeowService")

	// Create session repository
	sessionRepo := repository.NewPostgresRepo(db)

	// Create webhook service
	webhookService := webhooks.NewService()

	// Create WhatsApp service (infrastructure implementation)
	wmeowService := wmeow.NewService(container, waLogger, sessionRepo, cfg.GetMeow(), webhookService)

	// Create session service with proper dependencies
	domainSessionService := session.NewService()
	validator := validation.NewValidator()
	appSessionService := application.NewSessionApp(sessionRepo, domainSessionService, validator)

	// Create webhook application service
	webhookAppService := application.NewWebhookApp(sessionRepo, webhookService)

	// Connect active sessions on startup
	log.Info("Connecting active sessions on startup...")
	if err := wmeowService.ConnectOnStartup(context.Background()); err != nil {
		log.Errorf("Failed to connect active sessions on startup: %v", err)
	} else {
		log.Info("Active sessions connected successfully")
	}

	log.Info("Session service initialized")

	// Create authentication middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg, sessionRepo, log)

	// Use application service in handlers
	sessionHandler := handlers.NewSessionHandler(appSessionService, wmeowService)
	healthHandler := handlers.NewHealthHandler()
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

	// Add recovery middleware
	ginRouter.Use(gin.Recovery())

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

	routes.SetupRoutes(ginRouter, sessionHandler, healthHandler, messageHandler, chatHandler, groupHandler, communityHandler, webhookHandler, contactHandler, newsletterHandler, privacyHandler, authMiddleware)

	addr := fmt.Sprintf(":%s", cfg.GetServer().GetPort())
	log.Infof("Server listening on %s", addr)
	if err := ginRouter.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
