package handlers

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"

	"zpmeow/internal/application"
	"zpmeow/internal/application/dto/request"
	"zpmeow/internal/application/services"
	"zpmeow/internal/domain"
	"zpmeow/internal/infra"
	httpUtils "zpmeow/internal/infra/http/utils"
	"zpmeow/internal/infra/logger"
)

type SendHandler struct {
	sessionService    domain.SessionService
	meowService       *infra.MeowServiceImpl
	mediaFactory      *services.MediaStrategyFactoryImpl
	validationService *services.ApplicationValidationService
	responseService   *services.ResponseService
	logger            logger.Logger
}

func NewSendHandler(sessionService domain.SessionService, meowService *infra.MeowServiceImpl) *SendHandler {
	return &SendHandler{
		sessionService:    sessionService,
		meowService:       meowService,
		mediaFactory:      services.MediaStrategyFactory,
		validationService: services.DefaultValidationService,
		responseService:   services.DefaultResponseService,
		logger:            logger.GetLogger().Sub("send-handler"),
	}
}

func (h *SendHandler) SendText(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		httpUtils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req application.SendTextRequest
	if !httpUtils.ValidateAndBindJSON(c, &req) {
		return
	}

	// Validate phone number
	if err := h.validationService.ValidatePhoneNumber(req.Phone); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate message text
	if err := h.validationService.ValidateTextMessage(req.Body); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Implement actual text message sending
	h.logger.Infof("Sending text message to %s from session %s", req.Phone, sessionID)

	response := h.responseService.CreateSendResponse()
	httpUtils.RespondWithJSON(c, http.StatusOK, response)
}

func (h *SendHandler) SendMedia(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		httpUtils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req application.SendMediaRequest
	if !httpUtils.ValidateAndBindJSON(c, &req) {
		return
	}

	// Validate phone number
	if err := h.validationService.ValidatePhoneNumber(req.Phone); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate media type
	if err := h.mediaFactory.ValidateMediaType(req.MediaType); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Decode base64 media data
	mediaData, err := base64.StdEncoding.DecodeString(req.Media)
	if err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, "Invalid base64 media data")
		return
	}

	// Get strategy for media type
	strategy := h.mediaFactory.CreateStrategy(req.MediaType)
	if strategy == nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, "Unsupported media type: "+req.MediaType)
		return
	}

	// Validate media
	if err := strategy.ValidateMedia(mediaData); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, "Media validation failed: "+err.Error())
		return
	}

	// Process media
	processedData, mimeType, err := strategy.ProcessMedia(c.Request.Context(), mediaData, req.Filename)
	if err != nil {
		httpUtils.RespondWithError(c, http.StatusInternalServerError, "Media processing failed: "+err.Error())
		return
	}

	// TODO: Implement actual media sending using strategy
	h.logger.Infof("Sending %s media to %s from session %s (size: %d bytes, mime: %s)",
		req.MediaType, req.Phone, sessionID, len(processedData), mimeType)

	response := h.responseService.CreateSendResponse()
	httpUtils.RespondWithJSON(c, http.StatusOK, response)
}

func (h *SendHandler) SendImage(c *gin.Context) {
	h.sendMediaByType(c, "image")
}

func (h *SendHandler) SendAudio(c *gin.Context) {
	h.sendMediaByType(c, "audio")
}

func (h *SendHandler) SendDocument(c *gin.Context) {
	h.sendMediaByType(c, "document")
}

func (h *SendHandler) SendVideo(c *gin.Context) {
	h.sendMediaByType(c, "video")
}

func (h *SendHandler) SendSticker(c *gin.Context) {
	h.sendMediaByType(c, "sticker")
}

// sendMediaByType is a helper method that uses strategy pattern for specific media types
func (h *SendHandler) sendMediaByType(c *gin.Context, mediaType string) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		httpUtils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	// Parse request based on media type
	var phone, media, caption, filename string
	var err error

	switch mediaType {
	case "image":
		var req request.SendImageRequest
		if !httpUtils.ValidateAndBindJSON(c, &req) {
			return
		}
		phone, media, caption, filename = req.Phone, req.Image, req.Caption, ""
	case "audio":
		var req request.SendAudioRequest
		if !httpUtils.ValidateAndBindJSON(c, &req) {
			return
		}
		phone, media, caption, filename = req.Phone, req.Audio, "", ""
	case "video":
		var req request.SendVideoRequest
		if !httpUtils.ValidateAndBindJSON(c, &req) {
			return
		}
		phone, media, caption, filename = req.Phone, req.Video, req.Caption, ""
	case "document":
		var req request.SendDocumentRequest
		if !httpUtils.ValidateAndBindJSON(c, &req) {
			return
		}
		phone, media, caption, filename = req.Phone, req.Document, req.Caption, req.Filename
	case "sticker":
		var req request.SendStickerRequest
		if !httpUtils.ValidateAndBindJSON(c, &req) {
			return
		}
		phone, media, caption, filename = req.Phone, req.Sticker, "", ""
	default:
		httpUtils.RespondWithError(c, http.StatusBadRequest, "Unsupported media type: "+mediaType)
		return
	}

	// Validate phone number
	if err := h.validationService.ValidatePhoneNumber(phone); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate caption if provided
	if err := h.validationService.ValidateCaption(caption); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate filename if provided
	if err := h.validationService.ValidateFilename(filename); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Decode base64 media data
	mediaData, err := base64.StdEncoding.DecodeString(media)
	if err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, "Invalid base64 media data")
		return
	}

	// Get strategy for media type
	strategy := h.mediaFactory.CreateStrategy(mediaType)
	if strategy == nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, "Unsupported media type: "+mediaType)
		return
	}

	// Validate and process media using strategy
	if err := strategy.ValidateMedia(mediaData); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, "Media validation failed: "+err.Error())
		return
	}

	processedData, mimeType, err := strategy.ProcessMedia(c.Request.Context(), mediaData, filename)
	if err != nil {
		httpUtils.RespondWithError(c, http.StatusInternalServerError, "Media processing failed: "+err.Error())
		return
	}

	// TODO: Implement actual media sending using strategy
	h.logger.Infof("Sending %s to %s from session %s (size: %d bytes, mime: %s)",
		mediaType, phone, sessionID, len(processedData), mimeType)

	response := h.responseService.CreateSendResponse()
	httpUtils.RespondWithJSON(c, http.StatusOK, response)
}

func (h *SendHandler) SendLocation(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		httpUtils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req application.SendLocationRequest
	if !httpUtils.ValidateAndBindJSON(c, &req) {
		return
	}

	// Validate phone number
	if err := h.validationService.ValidatePhoneNumber(req.Phone); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Implement actual location message sending
	h.logger.Infof("Sending location (lat: %f, lng: %f) to %s from session %s",
		req.Latitude, req.Longitude, req.Phone, sessionID)

	response := h.responseService.CreateSendResponse()
	httpUtils.RespondWithJSON(c, http.StatusOK, response)
}

func (h *SendHandler) SendContact(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		httpUtils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req request.SendContactRequest
	if !httpUtils.ValidateAndBindJSON(c, &req) {
		return
	}

	// Validate phone number
	if err := h.validationService.ValidatePhoneNumber(req.Phone); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate vCard
	if err := h.validationService.ValidateVCard(req.Contact.VCard); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Implement actual contact message sending
	h.logger.Infof("Sending contact %s to %s from session %s",
		req.Contact.DisplayName, req.Phone, sessionID)

	response := h.responseService.CreateSendResponse()
	httpUtils.RespondWithJSON(c, http.StatusOK, response)
}

func (h *SendHandler) SendPoll(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		httpUtils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req request.SendPollRequest
	if !httpUtils.ValidateAndBindJSON(c, &req) {
		return
	}

	// Validate phone number
	if err := h.validationService.ValidatePhoneNumber(req.Phone); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate poll options
	if err := h.validationService.ValidatePollOptions(req.Options); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Implement actual poll message sending
	h.logger.Infof("Sending poll '%s' with %d options to %s from session %s",
		req.Name, len(req.Options), req.Phone, sessionID)

	response := h.responseService.CreateSendResponse()
	httpUtils.RespondWithJSON(c, http.StatusOK, response)
}

func (h *SendHandler) SendButton(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		httpUtils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req request.SendButtonRequest
	if !httpUtils.ValidateAndBindJSON(c, &req) {
		return
	}

	// Validate phone number
	if err := h.validationService.ValidatePhoneNumber(req.Phone); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate text message
	if err := h.validationService.ValidateTextMessage(req.Text); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Implement actual button message sending
	h.logger.Infof("Sending button message with %d buttons to %s from session %s",
		len(req.Buttons), req.Phone, sessionID)

	response := h.responseService.CreateSendResponse()
	httpUtils.RespondWithJSON(c, http.StatusOK, response)
}

func (h *SendHandler) SendList(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		httpUtils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req request.SendListRequest
	if !httpUtils.ValidateAndBindJSON(c, &req) {
		return
	}

	// Validate phone number
	if err := h.validationService.ValidatePhoneNumber(req.Phone); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate text message
	if err := h.validationService.ValidateTextMessage(req.Text); err != nil {
		httpUtils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// TODO: Implement actual list message sending
	h.logger.Infof("Sending list message with %d sections to %s from session %s",
		len(req.Sections), req.Phone, sessionID)

	response := h.responseService.CreateSendResponse()
	httpUtils.RespondWithJSON(c, http.StatusOK, response)
}



