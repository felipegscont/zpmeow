package session

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// WhatsAppService manages the lifecycle of whatsmeow clients
type WhatsAppService struct {
	clients   map[string]*whatsmeow.Client
	mu        sync.RWMutex
	db        *sqlx.DB
	container *sqlstore.Container
}

// NewWhatsAppService creates a new WhatsAppService
func NewWhatsAppService(db *sqlx.DB, container *sqlstore.Container) *WhatsAppService {
	return &WhatsAppService{
		clients:   make(map[string]*whatsmeow.Client),
		db:        db,
		container: container,
	}
}

// AddClient adds a new client to the service
func (s *WhatsAppService) AddClient(sessionID string, client *whatsmeow.Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[sessionID] = client
}

// GetClient retrieves a client from the service
func (s *WhatsAppService) GetClient(sessionID string) (*whatsmeow.Client, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	client, ok := s.clients[sessionID]
	return client, ok
}

// RemoveClient removes a client from the service
func (s *WhatsAppService) RemoveClient(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, sessionID)
}

// StartClient creates a new whatsmeow client and starts the connection process
func (s *WhatsAppService) StartClient(sessionID string) error {
	// Check if client already exists
	if _, ok := s.GetClient(sessionID); ok {
		return fmt.Errorf("client for session %s already exists", sessionID)
	}

	deviceStore, err := s.container.GetFirstDevice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get device: %w", err)
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	s.AddClient(sessionID, client)

	qrChan, _ := client.GetQRChannel(context.Background())

	go func() {
		err := client.Connect()
		if err != nil {
			fmt.Printf("Failed to connect: %v\n", err)
			s.updateSessionStatus(sessionID, "error")
			return
		}
	}()

	go func() {
		for evt := range qrChan {
			if evt.Event == "code" {
				qrCodeImage, err := qrcode.Encode(evt.Code, qrcode.Medium, 256)
				if err != nil {
					fmt.Printf("Failed to encode QR code: %v\n", err)
					continue
				}
				qrCodeBase64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrCodeImage)
				s.updateSessionQRCode(sessionID, qrCodeBase64)
				s.updateSessionStatus(sessionID, "connecting")
			} else if evt.Event == "success" {
				s.updateSessionStatus(sessionID, "connected")
				jid := client.Store.ID.String()
				s.updateSessionJID(sessionID, jid)
			}
		}
	}()

	return nil
}

// PairPhone pairs a phone with the session
func (s *WhatsAppService) PairPhone(sessionID, phoneNumber string) (string, error) {
	client, ok := s.GetClient(sessionID)
	if !ok {
		return "", fmt.Errorf("client for session %s not found", sessionID)
	}

	code, err := client.PairPhone(context.Background(), phoneNumber, true, whatsmeow.PairClientChrome, "ZPMeow")
	if err != nil {
		return "", fmt.Errorf("failed to pair phone: %w", err)
	}

	return code, nil
}

// DisconnectClient disconnects a client
func (s *WhatsAppService) DisconnectClient(sessionID string) error {
	client, ok := s.GetClient(sessionID)
	if !ok {
		return fmt.Errorf("client for session %s not found", sessionID)
	}
	client.Disconnect()
	s.updateSessionStatus(sessionID, "disconnected")
	return nil
}

// LogoutClient logs out a client
func (s *WhatsAppService) LogoutClient(sessionID string) error {
	client, ok := s.GetClient(sessionID)
	if !ok {
		return fmt.Errorf("client for session %s not found", sessionID)
	}
	_ = client.Logout(context.Background())
	s.RemoveClient(sessionID)
	s.updateSessionStatus(sessionID, "disconnected")
	return nil
}

// GetClientStatus returns the connection status of a client
func (s *WhatsAppService) GetClientStatus(sessionID string) (bool, bool) {
	client, ok := s.GetClient(sessionID)
	if !ok {
		return false, false
	}
	return client.IsConnected(), client.IsLoggedIn()
}

func (s *WhatsAppService) updateSessionStatus(sessionID, status string) {
	query := `UPDATE sessions SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := s.db.Exec(query, status, sessionID)
	if err != nil {
		fmt.Printf("Failed to update session status: %v\n", err)
	}
}

func (s *WhatsAppService) updateSessionQRCode(sessionID, qrCode string) {
	query := `UPDATE sessions SET qr_code = $1, updated_at = NOW() WHERE id = $2`
	_, err := s.db.Exec(query, qrCode, sessionID)
	if err != nil {
		fmt.Printf("Failed to update session QR code: %v\n", err)
	}
}

func (s *WhatsAppService) updateSessionJID(sessionID, jid string) {
	query := `UPDATE sessions SET whatsapp_jid = $1, updated_at = NOW() WHERE id = $2`
	_, err := s.db.Exec(query, jid, sessionID)
	if err != nil {
		fmt.Printf("Failed to update session JID: %v\n", err)
	}
}
