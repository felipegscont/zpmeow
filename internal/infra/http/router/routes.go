package router

import (
	"zpmeow/internal/infra/http/handler"
	"zpmeow/internal/infra/http/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configures the application routes
func SetupRoutes(
	router *gin.Engine,
	sessionHandler *handler.SessionHandler,
	healthHandler *handler.HealthHandler,
	sendHandler *handler.SendHandler,
	chatHandler *handler.ChatHandler,
	groupHandler *handler.GroupHandler,
) {
	// Add middlewares
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(gin.Recovery())

	// Health check route
	router.GET("/ping", healthHandler.Ping)

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

	// WuzAPI-style routes for session-based operations
	sessionAPIGroup := router.Group("/session/:sessionId")
	{
		// Send routes
		sendGroup := sessionAPIGroup.Group("/send")
		{
			sendGroup.POST("/text", sendHandler.SendText)
			sendGroup.POST("/image", sendHandler.SendImage)
			sendGroup.POST("/audio", sendHandler.SendAudio)
			sendGroup.POST("/document", sendHandler.SendDocument)
			sendGroup.POST("/video", sendHandler.SendVideo)
			sendGroup.POST("/sticker", sendHandler.SendSticker)
			sendGroup.POST("/location", sendHandler.SendLocation)
			sendGroup.POST("/contact", sendHandler.SendContact)
			sendGroup.POST("/buttons", sendHandler.SendButtons)
			sendGroup.POST("/list", sendHandler.SendList)
			sendGroup.POST("/poll", sendHandler.SendPoll)
		}

		// Chat routes
		chatGroup := sessionAPIGroup.Group("/chat")
		{
			chatGroup.POST("/presence", chatHandler.SetPresence)
			chatGroup.POST("/markread", chatHandler.MarkRead)
			chatGroup.POST("/react", chatHandler.React)
			chatGroup.POST("/delete", chatHandler.Delete)
			chatGroup.POST("/edit", chatHandler.Edit)
			chatGroup.POST("/download/image", chatHandler.DownloadImage)
			chatGroup.POST("/download/video", chatHandler.DownloadVideo)
			chatGroup.POST("/download/audio", chatHandler.DownloadAudio)
			chatGroup.POST("/download/document", chatHandler.DownloadDocument)
		}

		// Group routes
		groupGroup := sessionAPIGroup.Group("/group")
		{
			groupGroup.POST("/create", groupHandler.CreateGroup)
			groupGroup.GET("/list", groupHandler.ListGroups)
			groupGroup.GET("/info", groupHandler.GetGroupInfo)
			groupGroup.POST("/join", groupHandler.JoinGroup)
			groupGroup.POST("/leave", groupHandler.LeaveGroup)
			groupGroup.GET("/invitelink", groupHandler.GetInviteLink)
			groupGroup.POST("/inviteinfo", groupHandler.GetInviteInfo)
			groupGroup.POST("/participants/update", groupHandler.UpdateParticipants)
			groupGroup.POST("/name/set", groupHandler.SetName)
			groupGroup.POST("/topic/set", groupHandler.SetTopic)
			groupGroup.POST("/photo/set", groupHandler.SetPhoto)
			groupGroup.POST("/photo/remove", groupHandler.RemovePhoto)
			groupGroup.POST("/announce/set", groupHandler.SetAnnounce)
			groupGroup.POST("/locked/set", groupHandler.SetLocked)
			groupGroup.POST("/ephemeral/set", groupHandler.SetEphemeral)
		}
	}

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

