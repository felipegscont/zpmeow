package application

// Re-export application layer types to avoid aliases

// DTOs
import (
	sessionDTO "zpmeow/internal/application/dto/session"
	whatsappDTO "zpmeow/internal/application/dto/whatsapp"
	webhookDTO "zpmeow/internal/application/dto/webhook"
	"zpmeow/internal/application/services"
	sessionUseCase "zpmeow/internal/application/usecase/session"
)

// Session DTOs
type CreateSessionRequest = sessionDTO.CreateSessionRequest
type PairSessionRequest = sessionDTO.PairSessionRequest
type ProxyRequest = sessionDTO.ProxyRequest
type BaseSessionInfo = sessionDTO.BaseSessionInfo
type ExtendedSessionInfo = sessionDTO.ExtendedSessionInfo
type CreateSessionResponse = sessionDTO.CreateSessionResponse
type SessionInfoResponse = sessionDTO.SessionInfoResponse
type SessionListResponse = sessionDTO.SessionListResponse
type PairSessionResponse = sessionDTO.PairSessionResponse
type QRCodeResponse = sessionDTO.QRCodeResponse
type ProxyResponse = sessionDTO.ProxyResponse
type SuccessResponse = sessionDTO.SuccessResponse
type ErrorResponse = sessionDTO.ErrorResponse
type MessageResponse = sessionDTO.MessageResponse
type PingResponse = sessionDTO.PingResponse

// WhatsApp DTOs
type ContextInfo = whatsappDTO.ContextInfo
type SendTextRequest = whatsappDTO.SendTextRequest
type SendImageRequest = whatsappDTO.SendImageRequest
type SendAudioRequest = whatsappDTO.SendAudioRequest
type SendMediaRequest = whatsappDTO.SendMediaRequest
type SendLocationRequest = whatsappDTO.SendLocationRequest
type SendContactRequest = whatsappDTO.SendContactRequest
type SendPollRequest = whatsappDTO.SendPollRequest
type UserPresenceRequest = whatsappDTO.UserPresenceRequest
type CheckUserRequest = whatsappDTO.CheckUserRequest
type GetUserInfoRequest = whatsappDTO.GetUserInfoRequest
type ChatPresenceRequest = whatsappDTO.ChatPresenceRequest
type ChatMarkReadRequest = whatsappDTO.ChatMarkReadRequest
type ChatReactRequest = whatsappDTO.ChatReactRequest
type ChatDeleteRequest = whatsappDTO.ChatDeleteRequest
type ChatEditRequest = whatsappDTO.ChatEditRequest
type ChatDownloadRequest = whatsappDTO.ChatDownloadRequest

// Webhook DTOs
type SetWebhookRequest = webhookDTO.SetWebhookRequest
type UpdateWebhookRequest = webhookDTO.UpdateWebhookRequest
type WebhookResponse = webhookDTO.WebhookResponse
type WebhookPayload = webhookDTO.WebhookPayload

// Webhook functions
var (
	NewWebhookPayload = webhookDTO.NewWebhookPayload
)

// Service types
type SessionConverter = services.SessionConverter

// Service instances
var (
	SessionToDTOConverter = services.SessionToDTOConverter
)

// Use case constructors
var (
	NewSessionService = sessionUseCase.NewSessionService
)
