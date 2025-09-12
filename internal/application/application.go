package application

// Re-export application layer types to avoid aliases

// DTOs and Services
import (
	"zpmeow/internal/application/dto/request"
	"zpmeow/internal/application/dto/response"
	"zpmeow/internal/application/services"
	"zpmeow/internal/application/usecase"
)

// Session DTOs
type CreateSessionRequest = request.CreateSessionRequest
type PairSessionRequest = request.PairSessionRequest
type ProxyRequest = request.ProxyRequest
type BaseSessionInfo = response.BaseSessionInfo
type ExtendedSessionInfo = response.ExtendedSessionInfo
type CreateSessionResponse = response.CreateSessionResponse
type SessionInfoResponse = response.SessionInfoResponse
type SessionListResponse = response.SessionListResponse
type PairSessionResponse = response.PairSessionResponse
type QRCodeResponse = response.QRCodeResponse
type ProxyResponse = response.ProxyResponse
type SuccessResponse = response.SuccessResponse
type ErrorResponse = response.ErrorResponse
type MessageResponse = response.MessageResponse
type PingResponse = response.PingResponse

// WhatsApp DTOs
type ContextInfo = request.ContextInfo
type SendTextRequest = request.SendTextRequest
type SendImageRequest = request.SendImageRequest
type SendAudioRequest = request.SendAudioRequest
type SendMediaRequest = request.SendMediaRequest
type SendLocationRequest = request.SendLocationRequest
type SendContactRequest = request.SendContactRequest
type SendPollRequest = request.SendPollRequest

// Webhook DTOs
type SetWebhookRequest = request.SetWebhookRequest
type UpdateWebhookRequest = request.UpdateWebhookRequest
type WebhookResponse = response.WebhookResponse
type WebhookPayload = response.WebhookPayload

// Webhook functions
var (
	NewWebhookPayload = response.NewWebhookPayload
)

// Application Service types
type SessionApplicationService = services.SessionApplicationService

// Use case constructors
var (
	NewSessionService = usecase.NewSessionService
)
