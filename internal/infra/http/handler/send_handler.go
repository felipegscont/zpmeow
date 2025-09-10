package handler

import (
	"fmt"
	"net/http"
	"strings"

	"zpmeow/internal/domain/session"
	"zpmeow/internal/infra/meow"
	"zpmeow/internal/infra/logger"
	"zpmeow/internal/types"
	"zpmeow/internal/utils"

	"github.com/gin-gonic/gin"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
)

// SendHandler handles HTTP requests for message sending operations
type SendHandler struct {
	sessionService session.SessionService
	meowService    *meow.MeowServiceImpl
	logger         logger.Logger
}

// handleSendResponse is a helper function to standardize response handling for send operations
func (h *SendHandler) handleSendResponse(c *gin.Context, resp *whatsmeow.SendResponse, requestID string, err error, operation string) {
	if err != nil {
		h.logger.Errorf("Failed to %s: %v", operation, err)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to "+operation, err.Error())
		return
	}

	// Convert whatsmeow response to our API response format
	response := types.NewSendResponseFromWhatsmeow(resp, requestID)
	utils.RespondWithJSON(c, http.StatusOK, response)
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
	h.logger.Infof("DEBUG: SendText called with sessionID: %s", sessionID)

	if sessionID == "" {
		h.logger.Errorf("DEBUG: Session ID is empty")
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.SendTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("DEBUG: Failed to bind JSON: %v", err)
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	h.logger.Infof("DEBUG: Request parsed successfully - Phone: %s, Body: %s, ID: %s", req.Phone, req.Body, req.ID)

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		h.logger.Errorf("DEBUG: Invalid phone number format: %s", req.Phone)
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	h.logger.Infof("DEBUG: Phone number validation passed: %s", req.Phone)

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

	whatsmeowResp, err := h.meowService.SendTextMessage(c.Request.Context(), sessionID, req.Phone, req.Body, contextInfo)
	h.handleSendResponse(c, whatsmeowResp, req.ID, err, "send text message")
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

	// Try to parse as JSON first
	if err := c.ShouldBindJSON(&req); err != nil {
		// If JSON parsing fails, try form-data
		if err := c.ShouldBind(&req); err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
			return
		}
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Process media using unified processor (supports base64, URL, and form-data)
	var imageData []byte
	var mimeType string
	var err error

	// Check if there's a file upload
	file, _ := c.FormFile("image")
	if file == nil {
		file, _ = c.FormFile("media") // Also check for 'media' field
	}

	if file != nil {
		// Process form-data upload
		imageData, mimeType, err = utils.ProcessUnifiedMedia(c.Request.Context(), "", file, "image")
	} else {
		// Process base64 or URL
		media := req.Image
		if media == "" {
			media = c.PostForm("media") // Also check for 'media' field
		}
		imageData, mimeType, err = utils.ProcessUnifiedMedia(c.Request.Context(), media, nil, "image")
	}

	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid image data", err.Error())
		return
	}

	// Allow override of MIME type if provided
	if req.MimeType != "" {
		normalizedMimeType, err := utils.ValidateAndNormalizeMimeType(req.MimeType, "image")
		if err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid MIME type override", err.Error())
			return
		}
		mimeType = normalizedMimeType
	}

	// Validate size
	if err := utils.ValidateMediaSize(imageData, "image"); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid image size", err.Error())
		return
	}

	// Send image message through Meow service
	h.logger.Infof("Sending image message to %s from session %s", req.Phone, sessionID)

	whatsmeowResp, err := h.meowService.SendImageMessage(c.Request.Context(), sessionID, req.Phone, imageData, req.Caption, mimeType)
	h.handleSendResponse(c, whatsmeowResp, req.ID, err, "send image message")
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

	// Try to parse as JSON first
	if err := c.ShouldBindJSON(&req); err != nil {
		// If JSON parsing fails, try form-data
		if err := c.ShouldBind(&req); err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
			return
		}
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Process media using unified processor (supports base64, URL, and form-data)
	var audioData []byte
	var mimeType string
	var err error

	// Check if there's a file upload
	file, _ := c.FormFile("audio")
	if file == nil {
		file, _ = c.FormFile("media") // Also check for 'media' field
	}

	if file != nil {
		// Process form-data upload
		audioData, mimeType, err = utils.ProcessUnifiedMedia(c.Request.Context(), "", file, "audio")
	} else {
		// Process base64 or URL
		media := req.Audio
		if media == "" {
			media = c.PostForm("media") // Also check for 'media' field
		}
		audioData, mimeType, err = utils.ProcessUnifiedMedia(c.Request.Context(), media, nil, "audio")
	}

	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid audio data", err.Error())
		return
	}

	// Validate size
	if err := utils.ValidateMediaSize(audioData, "audio"); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid audio size", err.Error())
		return
	}

	// Send audio message through Meow service
	h.logger.Infof("Sending audio message to %s from session %s", req.Phone, sessionID)

	whatsmeowResp, err := h.meowService.SendAudioMessage(c.Request.Context(), sessionID, req.Phone, audioData, mimeType)
	h.handleSendResponse(c, whatsmeowResp, req.ID, err, "send audio message")
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

	// Try to parse as JSON first
	if err := c.ShouldBindJSON(&req); err != nil {
		// If JSON parsing fails, try form-data
		if err := c.ShouldBind(&req); err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
			return
		}
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Process media using unified processor (supports base64, URL, and form-data)
	var documentData []byte
	var mimeType string
	var err error

	// Check if there's a file upload
	file, _ := c.FormFile("document")
	if file == nil {
		file, _ = c.FormFile("media") // Also check for 'media' field
	}

	if file != nil {
		// Process form-data upload
		documentData, mimeType, err = utils.ProcessUnifiedMedia(c.Request.Context(), "", file, "document")
	} else {
		// Process base64 or URL
		media := req.Document
		if media == "" {
			media = c.PostForm("media") // Also check for 'media' field
		}
		documentData, mimeType, err = utils.ProcessUnifiedMedia(c.Request.Context(), media, nil, "document")
	}

	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid document data", err.Error())
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

	whatsmeowResp, err := h.meowService.SendDocumentMessage(c.Request.Context(), sessionID, req.Phone, documentData, filename, req.Caption, mimeType)
	h.handleSendResponse(c, whatsmeowResp, req.ID, err, "send document message")
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

	// Try to parse as JSON first
	if err := c.ShouldBindJSON(&req); err != nil {
		// If JSON parsing fails, try form-data
		if err := c.ShouldBind(&req); err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
			return
		}
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Process media using unified processor (supports base64, URL, and form-data)
	var videoData []byte
	var mimeType string
	var err error

	// Check if there's a file upload
	file, _ := c.FormFile("video")
	if file == nil {
		file, _ = c.FormFile("media") // Also check for 'media' field
	}

	if file != nil {
		// Process form-data upload
		videoData, mimeType, err = utils.ProcessUnifiedMedia(c.Request.Context(), "", file, "video")
	} else {
		// Process base64 or URL
		media := req.Video
		if media == "" {
			media = c.PostForm("media") // Also check for 'media' field
		}
		videoData, mimeType, err = utils.ProcessUnifiedMedia(c.Request.Context(), media, nil, "video")
	}

	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid video data", err.Error())
		return
	}

	// Validate size
	if err := utils.ValidateMediaSize(videoData, "video"); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid video size", err.Error())
		return
	}

	// Send video message through Meow service
	h.logger.Infof("Sending video message to %s from session %s", req.Phone, sessionID)

	whatsmeowResp, err := h.meowService.SendVideoMessage(c.Request.Context(), sessionID, req.Phone, videoData, req.Caption, mimeType)
	h.handleSendResponse(c, whatsmeowResp, req.ID, err, "send video message")
}

// SendMedia godoc
// @Summary Send a media message (unified endpoint)
// @Description Send any type of media message (image, audio, document, video) using a unified endpoint
// @Tags send
// @Accept json,multipart/form-data
// @Produce json
// @Param sessionId path string true "Session ID"
// @Param request body types.SendMediaRequest true "Media message request"
// @Success 200 {object} types.SendResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /session/{sessionId}/send/media [post]
func (h *SendHandler) SendMedia(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Session ID is required")
		return
	}

	var req types.SendMediaRequest

	// Check content type to determine parsing method
	contentType := c.GetHeader("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") {
		// Handle form-data
		req.Phone = c.PostForm("phone")
		req.MediaType = c.PostForm("mediaType")
		req.Caption = c.PostForm("caption")
		req.Filename = c.PostForm("filename")
		req.ID = c.PostForm("id")
		req.MimeType = c.PostForm("mimeType")
	} else {
		// Handle JSON
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid request", err.Error())
			return
		}
	}

	// Validate phone number format
	if !utils.IsValidPhoneNumber(req.Phone) {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	// Validate mediaType
	if err := utils.ValidateMediaType(req.MediaType); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid media type", err.Error())
		return
	}

	// Process media using unified processor (supports base64, URL, and form-data)
	var mediaData []byte
	var mimeType string
	var err error

	// Check if there's a file upload
	file, _ := c.FormFile("media")

	if file != nil {
		// Process form-data upload
		mediaData, mimeType, err = utils.ProcessUnifiedMedia(c.Request.Context(), "", file, req.MediaType)
	} else {
		// Process base64 or URL
		mediaData, mimeType, err = utils.ProcessUnifiedMedia(c.Request.Context(), req.Media, nil, req.MediaType)
	}

	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid media data", err.Error())
		return
	}

	// Validate size
	if err := utils.ValidateMediaSize(mediaData, req.MediaType); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid media size", err.Error())
		return
	}

	// Send media message through appropriate service method based on mediaType
	h.logger.Infof("Sending %s message to %s from session %s", req.MediaType, req.Phone, sessionID)

	var whatsmeowResp *whatsmeow.SendResponse

	switch req.MediaType {
	case "image":
		whatsmeowResp, err = h.meowService.SendImageMessage(c.Request.Context(), sessionID, req.Phone, mediaData, req.Caption, mimeType)
	case "audio":
		whatsmeowResp, err = h.meowService.SendAudioMessage(c.Request.Context(), sessionID, req.Phone, mediaData, mimeType)
	case "document":
		filename := req.Filename
		if filename == "" {
			ext := utils.GetFileExtension(mimeType)
			filename = "document" + ext
		}
		whatsmeowResp, err = h.meowService.SendDocumentMessage(c.Request.Context(), sessionID, req.Phone, mediaData, filename, req.Caption, mimeType)
	case "video":
		whatsmeowResp, err = h.meowService.SendVideoMessage(c.Request.Context(), sessionID, req.Phone, mediaData, req.Caption, mimeType)
	default:
		utils.RespondWithError(c, http.StatusBadRequest, "Unsupported media type")
		return
	}

	h.handleSendResponse(c, whatsmeowResp, req.ID, err, fmt.Sprintf("send %s message", req.MediaType))
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

	// Decode sticker data using universal decoder (supports all image formats for stickers)
	stickerData, mimeType, err := utils.DecodeUniversalMedia(req.Sticker, "sticker")
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid sticker data", err.Error())
		return
	}

	// Validate size
	if err := utils.ValidateMediaSize(stickerData, "sticker"); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid sticker size", err.Error())
		return
	}

	// Send sticker message through Meow service
	h.logger.Infof("Sending sticker message to %s from session %s", req.Phone, sessionID)

	whatsmeowResp, err := h.meowService.SendStickerMessage(c.Request.Context(), sessionID, req.Phone, stickerData, mimeType)
	h.handleSendResponse(c, whatsmeowResp, req.ID, err, "send sticker message")
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

	whatsmeowResp, err := h.meowService.SendLocationMessage(c.Request.Context(), sessionID, req.Phone, req.Latitude, req.Longitude, req.Name, req.Address)
	h.handleSendResponse(c, whatsmeowResp, req.ID, err, "send location message")
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

	whatsmeowResp, err := h.meowService.SendContactMessage(c.Request.Context(), sessionID, req.Phone, req.Contact.DisplayName, req.Contact.VCard)
	h.handleSendResponse(c, whatsmeowResp, req.ID, err, "send contact message")
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

	// Send buttons message through Meow service
	h.logger.Infof("Sending buttons message to %s from session %s", req.Phone, sessionID)

	whatsmeowResp, err := h.meowService.SendButtonsMessage(c.Request.Context(), sessionID, req.Phone, req.Text, req.Buttons, req.Footer)
	h.handleSendResponse(c, whatsmeowResp, req.ID, err, "send buttons message")
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

	// Send list message through Meow service
	h.logger.Infof("Sending list message to %s from session %s", req.Phone, sessionID)

	whatsmeowResp, err := h.meowService.SendListMessage(c.Request.Context(), sessionID, req.Phone, req.Text, req.ButtonText, req.Sections, req.Footer)
	h.handleSendResponse(c, whatsmeowResp, req.ID, err, "send list message")
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

	// Send poll message through Meow service
	h.logger.Infof("Sending poll message to %s from session %s", req.Phone, sessionID)

	whatsmeowResp, err := h.meowService.SendPollMessage(c.Request.Context(), sessionID, req.Phone, req.Name, req.Options, req.SelectableCount)
	h.handleSendResponse(c, whatsmeowResp, req.ID, err, "send poll message")
}
