package dto

import (
	"time"
)

// ============================================================================
// CONTACT REQUEST DTOs
// ============================================================================

// CheckContactRequest represents a request to check if contacts are on WhatsApp
type CheckContactRequest struct {
	Phones []string `json:"phones" binding:"required" example:"[\"5511999999999\", \"5511888888888\"]"`
}

// GetContactInfoRequest represents a request to get contact information
type GetContactInfoRequest struct {
	Phones []string `json:"phones" binding:"required" example:"[\"5511999999999\", \"5511888888888\"]"`
}

// GetAvatarRequest represents a request to get contact avatar
type GetAvatarRequest struct {
	Phone string `json:"phone" binding:"required" example:"5511999999999"`
}

// SetContactPresenceRequest represents a request to set global contact presence
type SetContactPresenceRequest struct {
	State string `json:"state" binding:"required" example:"available"`
}

// ============================================================================
// CONTACT DATA STRUCTURES
// ============================================================================

// ContactCheckResult represents the result of checking if a contact is on WhatsApp
type ContactCheckResult struct {
	Query        string `json:"query" example:"5511999999999"`
	IsInWhatsapp bool   `json:"is_in_whatsapp" example:"true"`
	JID          string `json:"jid" example:"5511999999999@s.whatsapp.net"`
	VerifiedName string `json:"verified_name,omitempty" example:"João Silva"`
}

// ContactInfo represents detailed contact information
type ContactInfo struct {
	JID          string `json:"jid" example:"5511999999999@s.whatsapp.net"`
	Name         string `json:"name,omitempty" example:"João Silva"`
	DisplayName  string `json:"display_name,omitempty" example:"João Silva"`
	VerifiedName string `json:"verified_name,omitempty" example:"João Silva Empresa"`
	Avatar       string `json:"avatar,omitempty" example:"https://..."`
	Status       string `json:"status,omitempty" example:"Disponível"`
	PictureID    string `json:"picture_id,omitempty" example:"pic_123"`
	DeviceCount  int    `json:"device_count,omitempty" example:"2"`
	Notify       string `json:"notify,omitempty" example:"João"`
	PushName     string `json:"push_name,omitempty" example:"João"`
	BusinessName string `json:"business_name,omitempty" example:"Empresa João"`
	IsBlocked    bool   `json:"is_blocked" example:"false"`
	IsMuted      bool   `json:"is_muted" example:"false"`
}

// AvatarInfo represents contact avatar information
type AvatarInfo struct {
	Phone     string    `json:"phone" example:"5511999999999"`
	JID       string    `json:"jid" example:"5511999999999@s.whatsapp.net"`
	AvatarURL string    `json:"avatar_url,omitempty" example:"https://pps.whatsapp.net/..."`
	PictureID string    `json:"picture_id,omitempty" example:"pic_123"`
	Timestamp time.Time `json:"timestamp" example:"2023-01-01T12:00:00Z"`
}

// ============================================================================
// CONTACT RESPONSE DTOs
// ============================================================================

// ContactResponse represents the standardized response format for contact operations
type ContactResponse struct {
	Success bool                  `json:"success"`
	Code    int                   `json:"code"`
	Data    ContactData           `json:"data"`
	Error   *ContactErrorResponse `json:"error,omitempty"`
}

// ContactData contains the response data for contact operations
type ContactData struct {
	Action       string               `json:"action" example:"check_contacts"`
	Status       string               `json:"status" example:"success"`
	Timestamp    time.Time            `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	CheckResults []ContactCheckResult `json:"check_results,omitempty"`
	ContactInfos []ContactInfo        `json:"contact_infos,omitempty"`
	Avatar       *AvatarInfo          `json:"avatar,omitempty"`
}

// ContactErrorResponse represents error information for contact operations
type ContactErrorResponse struct {
	Code    string `json:"code" example:"INVALID_PHONE"`
	Message string `json:"message" example:"Invalid phone number format"`
	Details string `json:"details,omitempty" example:"Phone number must include country code"`
}

// ContactsResponse represents the response for contacts operations
type ContactsResponse struct {
	Success bool                  `json:"success"`
	Code    int                   `json:"code"`
	Data    ContactsData          `json:"data"`
	Error   *ContactErrorResponse `json:"error,omitempty"`
}

// ContactsData contains contacts response data
type ContactsData struct {
	Action    string        `json:"action" example:"get_contacts"`
	Status    string        `json:"status" example:"success"`
	Timestamp time.Time     `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Contacts  []ContactInfo `json:"contacts"`
	Count     int           `json:"count" example:"10"`
}

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

// NewContactSuccessResponse creates a successful contact operation response
func NewContactSuccessResponse(action string, checkResults []ContactCheckResult, contactInfos []ContactInfo) *ContactResponse {
	return &ContactResponse{
		Success: true,
		Code:    200,
		Data: ContactData{
			Action:       action,
			Status:       "success",
			Timestamp:    time.Now(),
			CheckResults: checkResults,
			ContactInfos: contactInfos,
		},
	}
}

// NewContactErrorResponse creates an error response for contact operations
func NewContactErrorResponse(code int, errorCode, message, details string) *ContactResponse {
	return &ContactResponse{
		Success: false,
		Code:    code,
		Data: ContactData{
			Status:    "error",
			Timestamp: time.Now(),
		},
		Error: &ContactErrorResponse{
			Code:    errorCode,
			Message: message,
			Details: details,
		},
	}
}

// NewContactAvatarResponse creates a response for avatar operations
func NewContactAvatarResponse(avatar *AvatarInfo) *ContactResponse {
	return &ContactResponse{
		Success: true,
		Code:    200,
		Data: ContactData{
			Action:    "get_avatar",
			Status:    "success",
			Timestamp: time.Now(),
			Avatar:    avatar,
		},
	}
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

// ============================================================================
// BACKWARD COMPATIBILITY ALIASES
// ============================================================================

// Legacy aliases for backward compatibility
type CheckUserRequest = CheckContactRequest
type GetUserInfoRequest = GetContactInfoRequest
type SetUserPresenceRequest = SetContactPresenceRequest
type UserCheckResult = ContactCheckResult
type UserInfo = ContactInfo
type UserResponse = ContactResponse
type UserData = ContactData
type UserErrorResponse = ContactErrorResponse

// Legacy utility functions for backward compatibility
func NewUserSuccessResponse(action string, checkResults []ContactCheckResult, contactInfos []ContactInfo) *ContactResponse {
	return NewContactSuccessResponse(action, checkResults, contactInfos)
}

func NewUserErrorResponse(code int, errorCode, message, details string) *ContactResponse {
	return NewContactErrorResponse(code, errorCode, message, details)
}

func NewUserAvatarResponse(avatar *AvatarInfo) *ContactResponse {
	return NewContactAvatarResponse(avatar)
}
