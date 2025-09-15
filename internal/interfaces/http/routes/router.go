package routes

import (
	"zpmeow/docs"
	"zpmeow/internal/infrastructure/middleware"
	"zpmeow/internal/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(
	router *gin.Engine,
	sessionHandler *handlers.SessionHandler,
	healthHandler *handlers.HealthHandler,
	messageHandler *handlers.MessageHandler,
	chatHandler *handlers.ChatHandler,
	groupHandler *handlers.GroupHandler,
	webhookHandler *handlers.WebhookHandler,
	userHandler *handlers.UserHandler,
	newsletterHandler *handlers.NewsletterHandler,
	authMiddleware *middleware.AuthMiddleware,
) {

	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(gin.Recovery())

	// Middleware to set dynamic host for Swagger
	router.Use(func(c *gin.Context) {
		// Set the host dynamically based on the request
		if c.Request.URL.Path == "/swagger/doc.json" {
			host := c.Request.Host
			docs.SwaggerInfo.Host = host
		}
		c.Next()
	})

	router.GET("/ping", healthHandler.Ping)
	router.GET("/health", healthHandler.Health)

	sessionGroup := router.Group("/sessions")
	sessionGroup.Use(authMiddleware.AuthenticateGlobal())
	{
		sessionGroup.POST("/create", sessionHandler.CreateSession)
		sessionGroup.GET("/list", sessionHandler.GetSessions) // Updated to use GetSessions
		sessionGroup.GET("/:id/info", sessionHandler.GetSession)
		sessionGroup.DELETE("/:id/delete", sessionHandler.DeleteSession)
		sessionGroup.POST("/:id/connect", sessionHandler.ConnectSession)
		sessionGroup.POST("/:id/disconnect", sessionHandler.DisconnectSession)
		sessionGroup.GET("/:id/qr", sessionHandler.GetQRCode)
		sessionGroup.POST("/:id/pair", sessionHandler.PairPhone)
		sessionGroup.GET("/:id/status", sessionHandler.GetSessionStatus)
		sessionGroup.PUT("/:id/webhook", sessionHandler.UpdateSessionWebhook)
		sessionGroup.POST("/:id/regenerate-key", sessionHandler.RegenerateApiKey)
	}

	sessionAPIGroup := router.Group("/session/:sessionId")
	sessionAPIGroup.Use(authMiddleware.AuthenticateSession())
	{
		sendGroup := sessionAPIGroup.Group("/send")
		{
			sendGroup.POST("/text", messageHandler.SendText)
			sendGroup.POST("/media", messageHandler.SendMedia)
			sendGroup.POST("/location", messageHandler.SendLocation)
			sendGroup.POST("/contact", messageHandler.SendContact)
			sendGroup.POST("/image", messageHandler.SendImage)
			sendGroup.POST("/audio", messageHandler.SendAudio)
			sendGroup.POST("/document", messageHandler.SendDocument)
			sendGroup.POST("/video", messageHandler.SendVideo)
			sendGroup.POST("/sticker", messageHandler.SendSticker)
			sendGroup.POST("/buttons", messageHandler.SendButton)
			sendGroup.POST("/list", messageHandler.SendList)
			sendGroup.POST("/poll", messageHandler.SendPoll)
		}

		chatGroup := sessionAPIGroup.Group("/chat")
		{
			chatGroup.POST("/presence", chatHandler.SetPresence)
			chatGroup.POST("/markread", chatHandler.MarkAsRead)
			chatGroup.POST("/react", chatHandler.ReactToMessage)
			chatGroup.POST("/delete", chatHandler.DeleteMessage)
			chatGroup.POST("/edit", chatHandler.EditMessage)
			chatGroup.GET("/history", chatHandler.GetChatHistory)
			chatGroup.POST("/download/image", chatHandler.DownloadImage)
			chatGroup.POST("/download/video", chatHandler.DownloadVideo)
			chatGroup.POST("/download/audio", chatHandler.DownloadAudio)
			chatGroup.POST("/download/document", chatHandler.DownloadDocument)
		}

		groupGroup := sessionAPIGroup.Group("/group")
		{
			groupGroup.POST("/create", groupHandler.CreateGroup)
			groupGroup.GET("/list", groupHandler.ListGroups)
			groupGroup.POST("/info", groupHandler.GetGroupInfo)
			groupGroup.POST("/join", groupHandler.JoinGroup)
			groupGroup.POST("/join-with-invite", groupHandler.JoinGroupWithInvite)
			groupGroup.POST("/leave", groupHandler.LeaveGroup)
			groupGroup.POST("/invitelink", groupHandler.GetInviteLink)
			groupGroup.POST("/inviteinfo", groupHandler.GetInviteInfo)
			groupGroup.POST("/inviteinfo-specific", groupHandler.GetGroupInfoFromInvite)
			groupGroup.POST("/participants/update", groupHandler.UpdateParticipants)
			groupGroup.POST("/name/set", groupHandler.SetName)
			groupGroup.POST("/topic/set", groupHandler.SetTopic)
			groupGroup.POST("/photo/set", groupHandler.SetPhoto)
			groupGroup.POST("/photo/remove", groupHandler.RemovePhoto)
			groupGroup.POST("/announce/set", groupHandler.SetAnnounce)
			groupGroup.POST("/locked/set", groupHandler.SetLocked)
			groupGroup.POST("/ephemeral/set", groupHandler.SetEphemeral)
			groupGroup.POST("/join-approval/set", groupHandler.SetJoinApproval)
			groupGroup.POST("/member-add-mode/set", groupHandler.SetMemberAddMode)
			groupGroup.POST("/requests/list", groupHandler.GetGroupRequestParticipants)
			groupGroup.POST("/requests/update", groupHandler.UpdateGroupRequestParticipants)
			groupGroup.POST("/community/link", groupHandler.LinkGroup)
			groupGroup.POST("/community/unlink", groupHandler.UnlinkGroup)
			groupGroup.POST("/community/subgroups", groupHandler.GetSubGroups)
			groupGroup.POST("/community/participants", groupHandler.GetLinkedGroupsParticipants)
		}

		userGroup := sessionAPIGroup.Group("/user")
		{
			userGroup.POST("/presence", userHandler.SetPresence)
			userGroup.POST("/check", userHandler.CheckUser)
			userGroup.POST("/info", userHandler.GetUserInfo)
			userGroup.POST("/avatar", userHandler.GetAvatar)
			userGroup.GET("/contacts", userHandler.GetContacts)
		}

		newsletterGroup := sessionAPIGroup.Group("/newsletter")
		{
			newsletterGroup.GET("/list", newsletterHandler.ListNewsletters)
		}

		webhookGroup := sessionAPIGroup.Group("/webhook")
		{
			webhookGroup.POST("", webhookHandler.SetWebhook)
			webhookGroup.GET("", webhookHandler.GetWebhook)
			webhookGroup.PUT("", webhookHandler.UpdateWebhook)
			webhookGroup.DELETE("", webhookHandler.DeleteWebhook)
		}
	}

	// Configure Swagger to use dynamic host
	router.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler, ginswagger.URL("/swagger/doc.json")))
}
