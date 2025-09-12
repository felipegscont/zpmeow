

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

//	@schemes	http https
package main

import (
	"context"
	"fmt"

	// _ "zpmeow/docs" // Import for swagger docs - temporarily disabled
	"zpmeow/internal/application"
	"zpmeow/internal/config"
	"zpmeow/internal/domain"
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


	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}


	dbLog := logger.GetWALogger("Database")
	container, err := sqlstore.New(context.Background(), "postgres", cfg.DBUrl, dbLog)
	if err != nil {
		log.Fatalf("Failed to create whatsmeow container: %v", err)
	}


	sessionRepo := infra.NewPostgresSessionRepository()

	waLogger := logger.GetWALogger("MeowService")

	// Create session service first (without whatsapp service)
	sessionService := application.NewSessionService(sessionRepo, nil)

	// Create whatsapp service with session service
	whatsappService := infra.NewMeowService(db, container, waLogger, sessionService)

	// Update session service with whatsapp service
	sessionService = application.NewSessionService(sessionRepo, whatsappService)


	ctx := context.Background()
	if err := sessionService.ConnectOnStartup(ctx); err != nil {
		log.Warnf("Failed to connect sessions on startup: %v", err)
	}


	sessionHandler := handler.NewSessionHandler(sessionService)
	healthHandler := handler.NewHealthHandler()
	sendHandler := handler.NewSendHandler(sessionService, whatsappService.(*service.MeowServiceImpl))
	chatHandler := handler.NewChatHandler(sessionService, whatsappService.(*service.MeowServiceImpl))
	groupHandler := handler.NewGroupHandler(sessionService, whatsappService.(*service.MeowServiceImpl))
	webhookHandler := handler.NewWebhookHandler(sessionService)
	userHandler := handler.NewUserHandler(sessionService, whatsappService.(*service.MeowServiceImpl))
	newsletterHandler := handler.NewNewsletterHandler(sessionService, whatsappService.(*service.MeowServiceImpl))

	gin.SetMode(cfg.GinMode)


	ginRouter := gin.New()
	router.SetupRoutes(ginRouter, sessionHandler, healthHandler, sendHandler, chatHandler, groupHandler, webhookHandler, userHandler, newsletterHandler)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Infof("Server listening on %s", addr)
	if err := ginRouter.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
