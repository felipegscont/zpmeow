package router

import (
	handler "zpmeow/internal/infra/http/handlers"
	"zpmeow/internal/infra/http/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(
	router *gin.Engine,
	sessionHandler *handler.SessionHandler,
	healthHandler *handler.HealthHandler,
	sendHandler *handler.SendHandler,
	chatHandler *handler.ChatHandler,
	groupHandler *handler.GroupHandler,
	webhookHandler *handler.WebhookHandler,
	userHandler *handler.UserHandler,
	newsletterHandler *handler.NewsletterHandler,
) {

	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(gin.Recovery())

	router.GET("/ping", healthHandler.Ping)

	sessionGroup := router.Group("/sessions")
	{
		sessionGroup.POST("/create", sessionHandler.CreateSession)
		sessionGroup.GET("/list", sessionHandler.GetAllSessions)
		sessionGroup.GET("/:id/info", sessionHandler.GetSession)
		sessionGroup.DELETE("/:id/delete", sessionHandler.DeleteSession)
		// TODO: Implement other session endpoints
	}

	sessionAPIGroup := router.Group("/session/:sessionId")
	{

		sendGroup := sessionAPIGroup.Group("/send")
		{
			sendGroup.POST("/text", sendHandler.SendText)
			sendGroup.POST("/media", sendHandler.SendMedia)
			sendGroup.POST("/image", sendHandler.SendImage)
			sendGroup.POST("/audio", sendHandler.SendAudio)
			sendGroup.POST("/document", sendHandler.SendDocument)
			sendGroup.POST("/video", sendHandler.SendVideo)
			sendGroup.POST("/sticker", sendHandler.SendSticker)
			sendGroup.POST("/location", sendHandler.SendLocation)
			sendGroup.POST("/contact", sendHandler.SendContact)
			sendGroup.POST("/buttons", sendHandler.SendButton)
			sendGroup.POST("/list", sendHandler.SendList)
			sendGroup.POST("/poll", sendHandler.SendPoll)
		}

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

		// User routes
		userGroup := sessionAPIGroup.Group("/user")
		{
			userGroup.POST("/presence", userHandler.SetPresence)
			userGroup.POST("/check", userHandler.CheckUser)
			userGroup.POST("/info", userHandler.GetUserInfo)
			userGroup.POST("/avatar", userHandler.GetAvatar)
			userGroup.GET("/contacts", userHandler.GetContacts)
		}

		// Newsletter routes
		newsletterGroup := sessionAPIGroup.Group("/newsletter")
		{
			newsletterGroup.GET("/list", newsletterHandler.ListNewsletters)
		}

		// Webhook routes (session-specific)
		webhookGroup := sessionAPIGroup.Group("/webhook")
		{
			webhookGroup.POST("", webhookHandler.SetWebhook)
			webhookGroup.GET("", webhookHandler.GetWebhook)
			webhookGroup.PUT("", webhookHandler.UpdateWebhook)
			webhookGroup.DELETE("", webhookHandler.DeleteWebhook)
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
