package types

import (
	"testing"
	"time"

	"go.mau.fi/whatsmeow"
	waTypes "go.mau.fi/whatsmeow/types"
)

func TestNewSendResponseFromWhatsmeow(t *testing.T) {
	tests := []struct {
		name      string
		resp      *whatsmeow.SendResponse
		requestID string
		expected  SendResponse
	}{
		{
			name: "ServerID is zero - should return empty string",
			resp: &whatsmeow.SendResponse{
				Timestamp: time.Unix(1640995200, 0),
				ID:        "3EB0F488858A97C62C5E9B",
				ServerID:  0, // This is the problematic case
				Sender:    waTypes.NewJID("559984059035", "s.whatsapp.net"),
			},
			requestID: "",
			expected: SendResponse{
				Timestamp: 1640995200,
				ID:        "3EB0F488858A97C62C5E9B",
				ServerID:  "", // Should be empty string, not "\u0000"
				Sender:    "559984059035@s.whatsapp.net",
				Success:   true,
				MessageID: "3EB0F488858A97C62C5E9B",
			},
		},
		{
			name: "ServerID is non-zero - should return string representation",
			resp: &whatsmeow.SendResponse{
				Timestamp: time.Unix(1640995200, 0),
				ID:        "3EB0F488858A97C62C5E9B",
				ServerID:  12345,
				Sender:    waTypes.NewJID("559984059035", "s.whatsapp.net"),
			},
			requestID: "",
			expected: SendResponse{
				Timestamp: 1640995200,
				ID:        "3EB0F488858A97C62C5E9B",
				ServerID:  "12345",
				Sender:    "559984059035@s.whatsapp.net",
				Success:   true,
				MessageID: "3EB0F488858A97C62C5E9B",
			},
		},
		{
			name: "With custom request ID",
			resp: &whatsmeow.SendResponse{
				Timestamp: time.Unix(1640995200, 0),
				ID:        "3EB0F488858A97C62C5E9B",
				ServerID:  0,
				Sender:    waTypes.NewJID("559984059035", "s.whatsapp.net"),
			},
			requestID: "custom-message-id",
			expected: SendResponse{
				Timestamp: 1640995200,
				ID:        "3EB0F488858A97C62C5E9B",
				ServerID:  "",
				Sender:    "559984059035@s.whatsapp.net",
				Success:   true,
				MessageID: "custom-message-id",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewSendResponseFromWhatsmeow(tt.resp, tt.requestID)

			if result.Timestamp != tt.expected.Timestamp {
				t.Errorf("Timestamp = %v, want %v", result.Timestamp, tt.expected.Timestamp)
			}
			if result.ID != tt.expected.ID {
				t.Errorf("ID = %v, want %v", result.ID, tt.expected.ID)
			}
			if result.ServerID != tt.expected.ServerID {
				t.Errorf("ServerID = %q, want %q", result.ServerID, tt.expected.ServerID)
			}
			if result.Sender != tt.expected.Sender {
				t.Errorf("Sender = %v, want %v", result.Sender, tt.expected.Sender)
			}
			if result.Success != tt.expected.Success {
				t.Errorf("Success = %v, want %v", result.Success, tt.expected.Success)
			}
			if result.MessageID != tt.expected.MessageID {
				t.Errorf("MessageID = %v, want %v", result.MessageID, tt.expected.MessageID)
			}
		})
	}
}

func TestServerIDNullCharacterIssue(t *testing.T) {

	resp := &whatsmeow.SendResponse{
		Timestamp: time.Unix(1640995200, 0),
		ID:        "3EB0F488858A97C62C5E9B",
		ServerID:  0, // Zero value that was causing the issue
		Sender:    waTypes.NewJID("559984059035", "s.whatsapp.net"),
	}

	result := NewSendResponseFromWhatsmeow(resp, "")


	if result.ServerID == "\u0000" {
		t.Error("ServerID should not be null character \\u0000")
	}


	if result.ServerID != "" {
		t.Errorf("ServerID should be empty string when input is 0, got %q", result.ServerID)
	}
}
