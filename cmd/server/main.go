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
package main

import (
	"context"
	"fmt"

	// _ "zpmeow/docs" // Import for swagger docs - temporarily disabled
	"zpmeow/internal/application/usecase"
	"zpmeow/internal/config"
	"zpmeow/internal/infra"
	"zpmeow/internal/infra/database"
	handlers "zpmeow/internal/infra/http/handlers"
	"zpmeow/internal/infra/http/router"
	"zpmeow/internal/infra/logger"

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
	log := logger.Initialize(loggerConfig)
	logger.SetLogger(log)
	log.Info("Starting zpmeow server")

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(cfg); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	dbLog := logger.GetWALogger("Database")
	container, err := sqlstore.New(context.Background(), "postgres", cfg.DBUrl, dbLog)
	if err != nil {
		log.Fatalf("Failed to create whatsmeow container: %v", err)
	}

	sessionRepo := infra.NewPostgresSessionRepository(db)

	waLogger := logger.GetWALogger("MeowService")

	// Create session service first (without whatsapp service)
	sessionService := usecase.NewSessionService(sessionRepo, nil)

	// Create whatsapp service with session service
	whatsappService := infra.NewMeowService(db, container, waLogger, sessionService)

	// Update session service with whatsapp service
	sessionService = usecase.NewSessionService(sessionRepo, whatsappService)

	ctx := context.Background()
	if err := sessionService.ConnectOnStartup(ctx); err != nil {
		log.Warnf("Failed to connect sessions on startup: %v", err)
	}

	// Use application service in handlers
	sessionHandler := handlers.NewSessionHandler(sessionService)
	healthHandler := handlers.NewHealthHandler()
	sendHandler := handlers.NewSendHandler(sessionService, whatsappService.(*infra.MeowServiceImpl))
	chatHandler := handlers.NewChatHandler(sessionService, whatsappService.(*infra.MeowServiceImpl))
	groupHandler := handlers.NewGroupHandler(sessionService, whatsappService.(*infra.MeowServiceImpl))
	webhookHandler := handlers.NewWebhookHandler(sessionService)
	userHandler := handlers.NewUserHandler(sessionService, whatsappService.(*infra.MeowServiceImpl))
	newsletterHandler := handlers.NewNewsletterHandler(sessionService, whatsappService.(*infra.MeowServiceImpl))

	gin.SetMode(cfg.GinMode)

	ginRouter := gin.New()
	router.SetupRoutes(ginRouter, sessionHandler, healthHandler, sendHandler, chatHandler, groupHandler, webhookHandler, userHandler, newsletterHandler)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Infof("Server listening on %s", addr)
	if err := ginRouter.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
