package routes

import (
	"meow/docs"
	"meow/internal/infrastructure/middleware"
	"meow/internal/interfaces/http/handlers"

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
	communityHandler *handlers.CommunityHandler,
	webhookHandler *handlers.WebhookHandler,
	contactHandler *handlers.ContactHandler,
	newsletterHandler *handlers.NewsletterHandler,
	privacyHandler *handlers.PrivacyHandler,
	authMiddleware *middleware.AuthMiddleware,
) {

	router.Use(middleware.Logger())
	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
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
		sessionGroup.POST("/:id/pair", sessionHandler.PairPhone)
		sessionGroup.GET("/:id/status", sessionHandler.GetSessionStatus)
		sessionGroup.PUT("/:id/webhook", sessionHandler.UpdateSessionWebhook)
	}

	sessionAPIGroup := router.Group("/session/:sessionId")
	sessionAPIGroup.Use(authMiddleware.AuthenticateSession())
	{
		contacts := sessionAPIGroup.Group("/contacts")
		contacts.POST("/check", contactHandler.CheckUser)  // Check contact
		contacts.POST("/info", contactHandler.GetUserInfo) // Get contact info
		contacts.POST("/avatar", contactHandler.GetAvatar) // Get contact avatar
		contacts.GET("/list", contactHandler.GetContacts)  // List all contacts
		contacts.POST("/sync", contactHandler.GetContacts) // Sync contacts

		presence := sessionAPIGroup.Group("/presences")
		presence.PUT("/set", contactHandler.SetPresence)
		presence.GET("/get", contactHandler.GetUserInfo)
		presence.POST("/contact", contactHandler.GetUserInfo)
		presence.POST("/subscribe", contactHandler.CheckUser)
		presence.POST("/typing", chatHandler.SetPresence)
		presence.POST("/recording", chatHandler.SetPresence)

		privacy := sessionAPIGroup.Group("/privacy")

		privacy.PUT("/set", privacyHandler.SetAllPrivacySettings) // Set multiple privacy settings at once
		privacy.POST("/find", privacyHandler.FindPrivacySettings) // Find specific privacy settings
		privacy.GET("/blocklist", privacyHandler.GetBlocklist)    // Get blocked contacts list
		privacy.PUT("/blocklist", privacyHandler.UpdateBlocklist) // Block/unblock contacts

		message := sessionAPIGroup.Group("/message")
		send := message.Group("/send")
		send.POST("/text", messageHandler.SendText)
		send.POST("/image", messageHandler.SendImage)
		send.POST("/video", messageHandler.SendVideo)
		send.POST("/audio", messageHandler.SendAudio)
		send.POST("/document", messageHandler.SendDocument)
		send.POST("/sticker", messageHandler.SendSticker)
		send.POST("/contact", messageHandler.SendContact)
		send.POST("/location", messageHandler.SendLocation)
		send.POST("/media", messageHandler.SendMedia)
		send.POST("/buttons", messageHandler.SendButton)
		send.POST("/list", messageHandler.SendList)
		send.POST("/poll", messageHandler.SendPoll)

		message.POST("/markread", messageHandler.MarkAsRead)
		message.POST("/react", messageHandler.ReactToMessage)
		message.POST("/edit", messageHandler.EditMessage)
		message.POST("/delete", messageHandler.DeleteMessage)

		chat := sessionAPIGroup.Group("/chat")
		chat.POST("/presence", chatHandler.SetPresence)
		chat.GET("/history", chatHandler.GetChatHistory)

		download := chat.Group("/download")
		download.POST("/image", chatHandler.DownloadImage)
		download.POST("/video", chatHandler.DownloadVideo)
		download.POST("/audio", chatHandler.DownloadAudio)
		download.POST("/document", chatHandler.DownloadDocument)

		chat.POST("/list", chatHandler.ListChats)                          // List all chats (GetJoinedGroups + contacts)
		chat.POST("/info", chatHandler.GetChatInfo)                        // Get specific chat info
		chat.POST("/pin", chatHandler.PinChat)                             // Pin/unpin chat (app state)
		chat.POST("/mute", chatHandler.MuteChat)                           // Mute/unmute chat (app state)
		chat.POST("/archive", chatHandler.ArchiveChat)                     // Archive/unarchive chat (app state)
		chat.POST("/disappearing-timer", chatHandler.SetDisappearingTimer) // Set disappearing timer

		group := sessionAPIGroup.Group("/group")
		group.POST("/create", groupHandler.CreateGroup)
		group.GET("/list", groupHandler.ListGroups)
		group.POST("/info", groupHandler.GetGroupInfo)
		group.POST("/join", groupHandler.JoinGroup)
		group.POST("/join-with-invite", groupHandler.JoinGroupWithInvite)
		group.POST("/leave", groupHandler.LeaveGroup)
		group.POST("/invitelink", groupHandler.GetInviteLink)
		group.POST("/inviteinfo", groupHandler.GetInviteInfo)
		group.POST("/inviteinfo-specific", groupHandler.GetGroupInfoFromInvite)

		participants := group.Group("/participants")
		participants.POST("/update", groupHandler.UpdateParticipants)

		settings := group.Group("/settings")
		settings.POST("/name", groupHandler.SetName)
		settings.POST("/topic", groupHandler.SetTopic)
		settings.POST("/photo/set", groupHandler.SetPhoto)
		settings.POST("/photo/remove", groupHandler.RemovePhoto)
		settings.POST("/announce", groupHandler.SetAnnounce)
		settings.POST("/locked", groupHandler.SetLocked)
		settings.POST("/ephemeral", groupHandler.SetEphemeral)
		settings.POST("/join-approval", groupHandler.SetJoinApproval)
		settings.POST("/member-add-mode", groupHandler.SetMemberAddMode)

		requests := group.Group("/requests")
		requests.POST("/list", groupHandler.GetGroupRequestParticipants)
		requests.POST("/update", groupHandler.UpdateGroupRequestParticipants)

		community := sessionAPIGroup.Group("/community")
		community.POST("/link", communityHandler.LinkGroup)
		community.POST("/unlink", communityHandler.UnlinkGroup)
		community.POST("/subgroups", communityHandler.GetSubGroups)
		community.POST("/participants", communityHandler.GetLinkedGroupsParticipants)

		newsletter := sessionAPIGroup.Group("/newsletter")
		newsletter.POST("", newsletterHandler.CreateNewsletter)                                    // Create newsletter
		newsletter.GET("/list", newsletterHandler.ListNewsletters)                                 // List subscribed newsletters
		newsletter.GET("/:newsletterId", newsletterHandler.GetNewsletter)                          // Get newsletter info
		newsletter.POST("/:newsletterId/subscribe", newsletterHandler.SubscribeToNewsletter)       // Follow newsletter
		newsletter.POST("/:newsletterId/unsubscribe", newsletterHandler.UnsubscribeFromNewsletter) // Unfollow newsletter
		newsletter.POST("/:newsletterId/send", newsletterHandler.SendNewsletterMessage)            // Send message to newsletter

		newsletter.GET("/:newsletterId/messages", newsletterHandler.GetNewsletterMessages)      // Get newsletter messages
		newsletter.GET("/:newsletterId/updates", newsletterHandler.GetNewsletterMessageUpdates) // Get message updates
		newsletter.POST("/:newsletterId/mark-viewed", newsletterHandler.MarkNewsletterViewed)   // Mark messages as viewed

		newsletter.POST("/:newsletterId/reaction", newsletterHandler.SendNewsletterReaction)   // Send reaction
		newsletter.POST("/:newsletterId/mute", newsletterHandler.ToggleNewsletterMute)         // Mute/unmute
		newsletter.POST("/:newsletterId/live-updates", newsletterHandler.SubscribeLiveUpdates) // Subscribe to live updates

		newsletter.POST("/upload", newsletterHandler.UploadNewsletterMedia) // Upload media

		newsletter.GET("/invite/:inviteKey", newsletterHandler.GetNewsletterByInvite) // Get newsletter by invite

		webhook := sessionAPIGroup.Group("/webhook")
		webhook.POST("", webhookHandler.SetWebhook) // SET webhook
		webhook.GET("", webhookHandler.GetWebhook)  // GET webhook

		webhooks := sessionAPIGroup.Group("/webhooks")
		webhooks.GET("/events", webhookHandler.ListEvents) // LIST events
	}

	router.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler, ginswagger.URL("/swagger/doc.json")))
}
