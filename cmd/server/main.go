package main

import (
	"context"
	"fmt"

	_ "zpmeow/docs" // Import for swagger docs
	"zpmeow/internal/config"
	"zpmeow/internal/domain/session"
	"zpmeow/internal/infra/database"
	"zpmeow/internal/infra/http/handler"
	"zpmeow/internal/infra/http/router"
	"zpmeow/internal/infra/logger"
	"zpmeow/internal/infra/meow"

	"github.com/gin-gonic/gin"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		// Use basic logging before logger is initialized
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	// Initialize logger
	loggerConfig := logger.NewConfigAdapter(
		cfg.LogLevel,
		cfg.LogFormat,
		cfg.LogFilePath,
		cfg.LogFileFormat,
		cfg.LogConsoleColor,
		cfg.LogFileEnabled,
		cfg.LogFileCompress,
		cfg.LogFileMaxSize,
		cfg.LogFileMaxBackups,
		cfg.LogFileMaxAge,
	)
	log := logger.Initialize(loggerConfig)
	logger.SetLogger(log)
	log.Info("Starting zpmeow server")

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create WhatsApp store
	dbLog := logger.GetWALogger("Database")
	container, err := sqlstore.New(context.Background(), "postgres", cfg.DBUrl, dbLog)
	if err != nil {
		log.Fatalf("Failed to create whatsmeow container: %v", err)
	}

	// Initialize repositories
	sessionRepo := database.NewPostgresSessionRepository(db)

	// Initialize services
	waLogger := logger.GetWALogger("MeowService")
	whatsappService := meow.NewMeowService(db, container, waLogger)
	sessionService := session.NewSessionService(sessionRepo, whatsappService)

	// Connect previously connected sessions on startup (like wuzapi)
	ctx := context.Background()
	if err := sessionService.ConnectOnStartup(ctx); err != nil {
		log.Warnf("Failed to connect sessions on startup: %v", err)
	}

	// Initialize handlers
	sessionHandler := handler.NewSessionHandler(sessionService)
	healthHandler := handler.NewHealthHandler()
	sendHandler := handler.NewSendHandler(sessionService, whatsappService.(*meow.MeowServiceImpl))
	chatHandler := handler.NewChatHandler(sessionService, whatsappService.(*meow.MeowServiceImpl))
	groupHandler := handler.NewGroupHandler(sessionService, whatsappService.(*meow.MeowServiceImpl))

	gin.SetMode(cfg.GinMode)

	ginRouter := gin.Default()
	router.SetupRoutes(ginRouter, sessionHandler, healthHandler, sendHandler, chatHandler, groupHandler)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Infof("Server listening on %s", addr)
	if err := ginRouter.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
