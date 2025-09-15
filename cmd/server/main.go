//	@title			zpmeow WhatsApp API
//	@version		1.0
//	@description	A WhatsApp API server built with Go, inspired by wuzapi
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	zpmeow API Support
//	@contact.url	https://github.com/your-username/zpmeow
//	@contact.email	support@zpmeow.com

//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT

//	@host		localhost:8080
//	@BasePath	/

// @schemes	http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description API Key authentication. Simply provide your API key directly: "YOUR_API_KEY". The system automatically detects if it's a Global API Key (can access all sessions and session management) or a Session-specific API Key (can only access the specific session it belongs to).
package main

import (
	"context"
	"fmt"

	_ "zpmeow/docs" // Import for swagger docs
	"zpmeow/internal/application"
	"zpmeow/internal/domain/session"
	"zpmeow/internal/infrastructure/config"
	"zpmeow/internal/infrastructure/database"
	"zpmeow/internal/infrastructure/database/repository"
	"zpmeow/internal/infrastructure/logging"
	"zpmeow/internal/infrastructure/middleware"
	"zpmeow/internal/infrastructure/wameow"
	"zpmeow/internal/interfaces/http/handlers"
	"zpmeow/internal/interfaces/http/routes"
	"zpmeow/internal/shared/validation"

	"github.com/gin-gonic/gin"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {

		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	loggerConfig := cfg.GetLoggerConfig()
	log := logging.Initialize(loggerConfig)
	logging.SetLogger(log)
	log.Info("Starting zpmeow server")

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(cfg); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	dbLog := logging.GetWALogger("Database")
	container, err := sqlstore.New(context.Background(), "postgres", cfg.DBUrl, dbLog)
	if err != nil {
		log.Fatalf("Failed to create whatsmeow container: %v", err)
	}

	waLogger := logging.GetWALogger("MeowService")

	// Create session repository
	sessionRepo := repository.NewSessionRepo(db)

	// Create whatsapp service
	whatsappService := wameow.NewMeowService(container, waLogger, sessionRepo)

	// Create session service with proper dependencies
	domainSessionService := session.NewSessionService()
	validator := validation.NewValidator()
	appSessionService := application.NewSessionService(sessionRepo, domainSessionService, validator)

	// Connect active sessions on startup
	log.Info("Connecting active sessions on startup...")
	if err := whatsappService.ConnectOnStartup(context.Background()); err != nil {
		log.Errorf("Failed to connect active sessions on startup: %v", err)
	} else {
		log.Info("Active sessions connected successfully")
	}

	log.Info("Session service initialized")

	// Create authentication middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg, sessionRepo, log)

	// Use application service in handlers
	sessionHandler := handlers.NewSessionHandler(appSessionService, whatsappService)
	healthHandler := handlers.NewHealthHandler()
	meowServiceImpl := whatsappService.(*wameow.MeowService)
	messageHandler := handlers.NewMessageHandler(appSessionService, meowServiceImpl)
	chatHandler := handlers.NewChatHandler(nil, meowServiceImpl)
	groupHandler := handlers.NewGroupHandler(appSessionService, meowServiceImpl)
	webhookHandler := handlers.NewWebhookHandler(nil)
	userHandler := handlers.NewUserHandler(nil, meowServiceImpl)
	newsletterHandler := handlers.NewNewsletterHandler(nil, meowServiceImpl)

	gin.SetMode(cfg.GinMode)

	ginRouter := gin.New()
	routes.SetupRoutes(ginRouter, sessionHandler, healthHandler, messageHandler, chatHandler, groupHandler, webhookHandler, userHandler, newsletterHandler, authMiddleware)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Infof("Server listening on %s", addr)
	if err := ginRouter.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
