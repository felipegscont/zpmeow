package services

import (
	"context"
	"zpmeow/internal/application/dto/request"
	"zpmeow/internal/application/dto/response"
)

// SessionApplicationService defines the interface for session application services
// This layer handles DTOs and coordinates between handlers and domain services
type SessionApplicationService interface {
	// Session CRUD operations
	CreateSession(ctx context.Context, req *request.CreateSessionRequest) (*response.CreateSessionResponse, error)
	GetSession(ctx context.Context, req *request.GetSessionRequest) (*response.SessionInfoResponse, error)
	GetAllSessions(ctx context.Context) (*response.SessionListResponse, error)
	DeleteSession(ctx context.Context, req *request.DeleteSessionRequest) error

	// Session connection operations
	ConnectSession(ctx context.Context, req *request.ConnectSessionRequest) error
	DisconnectSession(ctx context.Context, req *request.DisconnectSessionRequest) error

	// Session pairing operations
	GetQRCode(ctx context.Context, req *request.GetQRCodeRequest) (*response.QRCodeResponse, error)
	PairWithPhone(ctx context.Context, req *request.PairSessionRequest) (*response.PairSessionResponse, error)

	// Proxy operations
	SetProxy(ctx context.Context, req *request.SetProxyRequest) (*response.ProxyResponse, error)
	ClearProxy(ctx context.Context, req *request.ClearProxyRequest) error

	// Health and status
	GetSessionStatus(ctx context.Context, req *request.GetSessionStatusRequest) (*response.SessionStatusResponse, error)
}
