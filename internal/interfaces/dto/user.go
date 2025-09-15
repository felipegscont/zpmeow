package dto

import (
	"time"
)

// ============================================================================
// USER REQUEST DTOs
// ============================================================================

// CheckUserRequest represents a request to check if users are on WhatsApp
type CheckUserRequest struct {
	Phones []string `json:"phones" binding:"required" example:"[\"5511999999999\", \"5511888888888\"]"`
}

// GetUserInfoRequest represents a request to get user information
type GetUserInfoRequest struct {
	Phones []string `json:"phones" binding:"required" example:"[\"5511999999999\", \"5511888888888\"]"`
}

// GetAvatarRequest represents a request to get user avatar
type GetAvatarRequest struct {
	Phone string `json:"phone" binding:"required" example:"5511999999999"`
}

// SetUserPresenceRequest represents a request to set global user presence
type SetUserPresenceRequest struct {
	State string `json:"state" binding:"required" example:"available"`
}

// ============================================================================
// USER DATA STRUCTURES
// ============================================================================

// UserCheckResult represents the result of checking if a user is on WhatsApp
type UserCheckResult struct {
	Query        string `json:"query" example:"5511999999999"`
	IsInWhatsapp bool   `json:"is_in_whatsapp" example:"true"`
	JID          string `json:"jid" example:"5511999999999@s.whatsapp.net"`
	VerifiedName string `json:"verified_name,omitempty" example:"João Silva"`
}

// UserInfo represents detailed user information
type UserInfo struct {
	JID          string `json:"jid" example:"5511999999999@s.whatsapp.net"`
	DisplayName  string `json:"display_name,omitempty" example:"João Silva"`
	VerifiedName string `json:"verified_name,omitempty" example:"João Silva Empresa"`
	Avatar       string `json:"avatar,omitempty" example:"https://..."`
	Status       string `json:"status,omitempty" example:"Disponível"`
	PictureID    string `json:"picture_id,omitempty" example:"pic_123"`
	DeviceCount  int    `json:"device_count,omitempty" example:"2"`
}

// AvatarInfo represents user avatar information
type AvatarInfo struct {
	Phone     string    `json:"phone" example:"5511999999999"`
	JID       string    `json:"jid" example:"5511999999999@s.whatsapp.net"`
	AvatarURL string    `json:"avatar_url,omitempty" example:"https://pps.whatsapp.net/..."`
	PictureID string    `json:"picture_id,omitempty" example:"pic_123"`
	Timestamp time.Time `json:"timestamp" example:"2023-01-01T12:00:00Z"`
}

// ContactInfo represents contact information
type ContactInfo struct {
	JID          string `json:"jid" example:"5511999999999@s.whatsapp.net"`
	Name         string `json:"name,omitempty" example:"João Silva"`
	Notify       string `json:"notify,omitempty" example:"João"`
	PushName     string `json:"push_name,omitempty" example:"João"`
	BusinessName string `json:"business_name,omitempty" example:"Empresa João"`
	IsBlocked    bool   `json:"is_blocked" example:"false"`
	IsMuted      bool   `json:"is_muted" example:"false"`
}

// ============================================================================
// USER RESPONSE DTOs
// ============================================================================

// UserResponse represents the standardized response format for user operations
type UserResponse struct {
	Success bool                `json:"success"`
	Code    int                 `json:"code"`
	Data    UserData            `json:"data"`
	Error   *UserErrorResponse  `json:"error,omitempty"`
}

// UserData contains the response data for user operations
type UserData struct {
	Action       string             `json:"action" example:"check_users"`
	Status       string             `json:"status" example:"success"`
	Timestamp    time.Time          `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	CheckResults []UserCheckResult  `json:"check_results,omitempty"`
	UserInfos    []UserInfo         `json:"user_infos,omitempty"`
	Avatar       *AvatarInfo        `json:"avatar,omitempty"`
}

// UserErrorResponse represents error information for user operations
type UserErrorResponse struct {
	Code    string `json:"code" example:"INVALID_PHONE"`
	Message string `json:"message" example:"Invalid phone number format"`
	Details string `json:"details,omitempty" example:"Phone number must include country code"`
}

// Note: Group DTOs have been moved to internal/interfaces/dto/group.go

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

// NewUserSuccessResponse creates a successful user operation response
func NewUserSuccessResponse(action string, checkResults []UserCheckResult, userInfos []UserInfo) *UserResponse {
	return &UserResponse{
		Success: true,
		Code:    200,
		Data: UserData{
			Action:       action,
			Status:       "success",
			Timestamp:    time.Now(),
			CheckResults: checkResults,
			UserInfos:    userInfos,
		},
	}
}

// NewUserErrorResponse creates an error response for user operations
func NewUserErrorResponse(code int, errorCode, message, details string) *UserResponse {
	return &UserResponse{
		Success: false,
		Code:    code,
		Data: UserData{
			Status:    "error",
			Timestamp: time.Now(),
		},
		Error: &UserErrorResponse{
			Code:    errorCode,
			Message: message,
			Details: details,
		},
	}
}

// NewUserAvatarResponse creates a response for avatar operations
func NewUserAvatarResponse(avatar *AvatarInfo) *UserResponse {
	return &UserResponse{
		Success: true,
		Code:    200,
		Data: UserData{
			Action:    "get_avatar",
			Status:    "success",
			Timestamp: time.Now(),
			Avatar:    avatar,
		},
	}
}

// ContactsResponse represents the response for contacts operations
type ContactsResponse struct {
	Success bool          `json:"success"`
	Code    int           `json:"code"`
	Data    ContactsData  `json:"data"`
	Error   *UserErrorResponse `json:"error,omitempty"`
}

// ContactsData contains contacts response data
type ContactsData struct {
	Action    string        `json:"action" example:"get_contacts"`
	Status    string        `json:"status" example:"success"`
	Timestamp time.Time     `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Contacts  []ContactInfo `json:"contacts"`
	Count     int           `json:"count" example:"10"`
}

// NewContactsResponse creates a response for contacts operations
func NewContactsResponse(contacts []ContactInfo) *ContactsResponse {
	return &ContactsResponse{
		Success: true,
		Code:    200,
		Data: ContactsData{
			Action:    "get_contacts",
			Status:    "success",
			Timestamp: time.Now(),
			Contacts:  contacts,
			Count:     len(contacts),
		},
	}
}
