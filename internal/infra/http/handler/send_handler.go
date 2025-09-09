package handler

import (
	"net/http"
	"time"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infra/meow"
	"zpmeow/internal/infra/logger"
	"zpmeow/internal/types"
	"zpmeow/internal/utils"

	"github.com/gin-gonic/gin"
	"go.mau.fi/whatsmeow/proto/waE2E"
)

// SendHandler handles HTTP requests for message sending operations
type SendHandler struct {
	sessionService session.SessionService
	meowService    *meow.MeowServiceImpl
	logger         logger.Logger
}

// NewSendHandler creates a new send handler
func NewSendHandler(sessionService session.SessionService, meowService *meow.MeowServiceImpl) *SendHandler {
	return &SendHandler{
		sessionService: sessionService,
		meowService:    meowService,
		logger:         logger.GetLogger().Sub("send-handler"),
	}
}

// SendText godoc
// @Summary Send a text message
// @Description Send a text message to a WhatsApp contact
// @Tags send
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.SendTextRequest true "Text message request"
// @Success 200 {object} types.SendResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/send/text [post]
func (h *SendHandler) SendText(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.SendTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Send message through Meow service
	h.logger.Infof("Sending text message to %s from session %s", req.Phone, sessionID)

	// Convert ContextInfo if provided
	var contextInfo *waE2E.ContextInfo
	if req.ContextInfo.StanzaID != "" {
		contextInfo = &waE2E.ContextInfo{
			StanzaID:    &req.ContextInfo.StanzaID,
			Participant: &req.ContextInfo.Participant,
		}
	}

	err := h.meowService.SendTextMessage(c.Request.Context(), sessionID, req.Phone, req.Body, contextInfo)
	if err != nil {
		h.logger.Errorf("Failed to send text message: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to send message", err.Error())
		return
	}

	response := types.SendResponse{
		Success:   true,
		MessageID: req.ID,
		Timestamp: time.Now().Unix(),
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

// SendImage godoc
// @Summary Send an image message
// @Description Send an image message to a WhatsApp contact
// @Tags send
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.SendImageRequest true "Image message request"
// @Success 200 {object} types.SendResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/send/image [post]
func (h *SendHandler) SendImage(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.SendImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Decode base64 image data
	imageData, mimeType, err := utils.DecodeBase64Media(req.Image)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid image data", err.Error())
		return
	}

	// Validate MIME type
	if req.MimeType != "" {
		mimeType = req.MimeType
	}
	if err := utils.ValidateMimeType(mimeType, "image"); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid image type", err.Error())
		return
	}

	// Validate size
	if err := utils.ValidateMediaSize(imageData, "image"); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid image size", err.Error())
		return
	}

	// Send image message through Meow service
	h.logger.Infof("Sending image message to %s from session %s", req.Phone, sessionID)

	err = h.meowService.SendImageMessage(c.Request.Context(), sessionID, req.Phone, imageData, req.Caption, mimeType)
	if err != nil {
		h.logger.Errorf("Failed to send image message: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to send image", err.Error())
		return
	}

	response := types.SendResponse{
		Success:   true,
		MessageID: req.ID,
		Timestamp: time.Now().Unix(),
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

// SendAudio godoc
// @Summary Send an audio message
// @Description Send an audio message to a WhatsApp contact
// @Tags send
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.SendAudioRequest true "Audio message request"
// @Success 200 {object} types.SendResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/send/audio [post]
func (h *SendHandler) SendAudio(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.SendAudioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Decode base64 audio data
	audioData, mimeType, err := utils.DecodeBase64Media(req.Audio)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid audio data", err.Error())
		return
	}

	// Validate MIME type
	if err := utils.ValidateMimeType(mimeType, "audio"); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid audio type", err.Error())
		return
	}

	// Validate size
	if err := utils.ValidateMediaSize(audioData, "audio"); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid audio size", err.Error())
		return
	}

	// Send audio message through Meow service
	h.logger.Infof("Sending audio message to %s from session %s", req.Phone, sessionID)

	err = h.meowService.SendAudioMessage(c.Request.Context(), sessionID, req.Phone, audioData, mimeType)
	if err != nil {
		h.logger.Errorf("Failed to send audio message: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to send audio", err.Error())
		return
	}

	response := types.SendResponse{
		Success:   true,
		MessageID: req.ID,
		Timestamp: time.Now().Unix(),
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

// SendDocument godoc
// @Summary Send a document message
// @Description Send a document message to a WhatsApp contact
// @Tags send
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.SendDocumentRequest true "Document message request"
// @Success 200 {object} types.SendResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/send/document [post]
func (h *SendHandler) SendDocument(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.SendDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Decode base64 document data
	documentData, mimeType, err := utils.DecodeBase64Media(req.Document)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid document data", err.Error())
		return
	}

	// Validate MIME type
	if err := utils.ValidateMimeType(mimeType, "document"); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid document type", err.Error())
		return
	}

	// Validate size
	if err := utils.ValidateMediaSize(documentData, "document"); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid document size", err.Error())
		return
	}

	// Generate filename if not provided
	filename := req.Filename
	if filename == "" {
		ext := utils.GetFileExtension(mimeType)
		filename = "document" + ext
	}

	// Send document message through Meow service
	h.logger.Infof("Sending document message to %s from session %s", req.Phone, sessionID)

	err = h.meowService.SendDocumentMessage(c.Request.Context(), sessionID, req.Phone, documentData, filename, req.Caption, mimeType)
	if err != nil {
		h.logger.Errorf("Failed to send document message: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to send document", err.Error())
		return
	}

	response := types.SendResponse{
		Success:   true,
		MessageID: req.ID,
		Timestamp: time.Now().Unix(),
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

// SendVideo godoc
// @Summary Send a video message
// @Description Send a video message to a WhatsApp contact
// @Tags send
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.SendVideoRequest true "Video message request"
// @Success 200 {object} types.SendResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/send/video [post]
func (h *SendHandler) SendVideo(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.SendVideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Decode base64 video data
	videoData, mimeType, err := utils.DecodeBase64Media(req.Video)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid video data", err.Error())
		return
	}

	// Validate MIME type
	if err := utils.ValidateMimeType(mimeType, "video"); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid video type", err.Error())
		return
	}

	// Validate size
	if err := utils.ValidateMediaSize(videoData, "video"); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid video size", err.Error())
		return
	}

	// Send video message through Meow service
	h.logger.Infof("Sending video message to %s from session %s", req.Phone, sessionID)

	err = h.meowService.SendVideoMessage(c.Request.Context(), sessionID, req.Phone, videoData, req.Caption, mimeType)
	if err != nil {
		h.logger.Errorf("Failed to send video message: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to send video", err.Error())
		return
	}

	response := types.SendResponse{
		Success:   true,
		MessageID: req.ID,
		Timestamp: time.Now().Unix(),
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

// SendSticker godoc
// @Summary Send a sticker message
// @Description Send a sticker message to a WhatsApp contact
// @Tags send
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.SendStickerRequest true "Sticker message request"
// @Success 200 {object} types.SendResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/send/sticker [post]
func (h *SendHandler) SendSticker(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.SendStickerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Decode base64 sticker data
	stickerData, mimeType, err := utils.DecodeBase64Media(req.Sticker)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid sticker data", err.Error())
		return
	}

	// Validate MIME type
	if err := utils.ValidateMimeType(mimeType, "sticker"); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid sticker type", err.Error())
		return
	}

	// Validate size
	if err := utils.ValidateMediaSize(stickerData, "sticker"); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid sticker size", err.Error())
		return
	}

	// Send sticker message through Meow service
	h.logger.Infof("Sending sticker message to %s from session %s", req.Phone, sessionID)

	err = h.meowService.SendStickerMessage(c.Request.Context(), sessionID, req.Phone, stickerData, mimeType)
	if err != nil {
		h.logger.Errorf("Failed to send sticker message: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to send sticker", err.Error())
		return
	}

	response := types.SendResponse{
		Success:   true,
		MessageID: req.ID,
		Timestamp: time.Now().Unix(),
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

// SendLocation godoc
// @Summary Send a location message
// @Description Send a location message to a WhatsApp contact
// @Tags send
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.SendLocationRequest true "Location message request"
// @Success 200 {object} types.SendResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/send/location [post]
func (h *SendHandler) SendLocation(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.SendLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Validate coordinates
	if req.Latitude < -90 || req.Latitude > 90 {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid latitude")
		return
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid longitude")
		return
	}

	// Send location message through Meow service
	h.logger.Infof("Sending location message to %s from session %s", req.Phone, sessionID)

	err := h.meowService.SendLocationMessage(c.Request.Context(), sessionID, req.Phone, req.Latitude, req.Longitude, req.Name, req.Address)
	if err != nil {
		h.logger.Errorf("Failed to send location message: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to send location", err.Error())
		return
	}

	response := types.SendResponse{
		Success:   true,
		MessageID: req.ID,
		Timestamp: time.Now().Unix(),
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

// SendContact godoc
// @Summary Send a contact message
// @Description Send a contact message to a WhatsApp contact
// @Tags send
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.SendContactRequest true "Contact message request"
// @Success 200 {object} types.SendResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/send/contact [post]
func (h *SendHandler) SendContact(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.SendContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Send contact message through Meow service
	h.logger.Infof("Sending contact message to %s from session %s", req.Phone, sessionID)

	err := h.meowService.SendContactMessage(c.Request.Context(), sessionID, req.Phone, req.Contact.DisplayName, req.Contact.VCard)
	if err != nil {
		h.logger.Errorf("Failed to send contact message: %v", err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to send contact", err.Error())
		return
	}

	response := types.SendResponse{
		Success:   true,
		MessageID: req.ID,
		Timestamp: time.Now().Unix(),
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

// SendButtons godoc
// @Summary Send an interactive buttons message
// @Description Send an interactive buttons message to a WhatsApp contact
// @Tags send
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.SendButtonsRequest true "Buttons message request"
// @Success 200 {object} types.SendResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/send/buttons [post]
func (h *SendHandler) SendButtons(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.SendButtonsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Validate buttons
	if len(req.Buttons) == 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "At least one button is required")
		return
	}
	if len(req.Buttons) > 3 {
		utils.RespondWithError(c, http.StatusBadRequest, "Maximum 3 buttons allowed")
		return
	}

	// TODO: Implement actual buttons sending through Meow service
	h.logger.Infof("Sending buttons message to %s from session %s", req.Phone, sessionID)

	response := types.SendResponse{
		Success:   true,
		MessageID: "mock-buttons-id-" + sessionID,
		Timestamp: 1640995200,
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

// SendList godoc
// @Summary Send an interactive list message
// @Description Send an interactive list message to a WhatsApp contact
// @Tags send
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.SendListRequest true "List message request"
// @Success 200 {object} types.SendResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/send/list [post]
func (h *SendHandler) SendList(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.SendListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Validate sections
	if len(req.Sections) == 0 {
		utils.RespondWithError(c, http.StatusBadRequest, "At least one section is required")
		return
	}
	if len(req.Sections) > 10 {
		utils.RespondWithError(c, http.StatusBadRequest, "Maximum 10 sections allowed")
		return
	}

	// TODO: Implement actual list sending through Meow service
	h.logger.Infof("Sending list message to %s from session %s", req.Phone, sessionID)

	response := types.SendResponse{
		Success:   true,
		MessageID: "mock-list-id-" + sessionID,
		Timestamp: 1640995200,
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

// SendPoll godoc
// @Summary Send a poll message
// @Description Send a poll message to a WhatsApp contact
// @Tags send
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.SendPollRequest true "Poll message request"
// @Success 200 {object} types.SendResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/send/poll [post]
func (h *SendHandler) SendPoll(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.SendPollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Validate poll options
	if len(req.Options) < 2 {
		utils.RespondWithError(c, http.StatusBadRequest, "At least 2 options are required")
		return
	}
	if len(req.Options) > 12 {
		utils.RespondWithError(c, http.StatusBadRequest, "Maximum 12 options allowed")
		return
	}

	// Validate selectable count
	if req.SelectableCount <= 0 {
		req.SelectableCount = 1 // Default to single selection
	}
	if req.SelectableCount > len(req.Options) {
		utils.RespondWithError(c, http.StatusBadRequest, "Selectable count cannot exceed number of options")
		return
	}

	// TODO: Implement actual poll sending through Meow service
	h.logger.Infof("Sending poll message to %s from session %s", req.Phone, sessionID)

	response := types.SendResponse{
		Success:   true,
		MessageID: "mock-poll-id-" + sessionID,
		Timestamp: 1640995200,
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}
