package session

import (
	"net/url"
	"regexp"
)

var (
	sessionNameMinLength = 3
	sessionNameMaxLength = 50
	sessionNameRegex     = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

type Service interface {
	CanConnect(session *Session) bool
	CanDisconnect(session *Session) bool
	CanDelete(session *Session) bool

	ValidateStatusTransition(current, newStatus Status) error

	ValidateSessionConfiguration(session *Session) error

	CanRegenerateApiKey(session *Session) bool

	CanSetProxy(session *Session) bool

	CanSubscribeToEvents(session *Session) bool

	// Device Management
	ValidateDeviceConnection(session *Session, deviceJID string) error
}

type DomainService struct{}

func NewService() *DomainService {
	return &DomainService{}
}

func (s *DomainService) CanConnect(session *Session) bool {
	return session.IsDisconnected() || session.HasError() || session.IsConnecting()
}

func (s *DomainService) CanDisconnect(session *Session) bool {
	return session.IsConnected() || session.IsConnecting()
}

func (s *DomainService) CanDelete(session *Session) bool {
	return session.IsDisconnected()
}

func (s *DomainService) ValidateStatusTransition(current, newStatus Status) error {
	return ValidateSessionStatus(current, newStatus)
}

func (s *DomainService) ValidateSessionConfiguration(session *Session) error {
	if err := session.Validate(); err != nil {
		return err
	}

	if session.HasProxy() {
		if err := ValidateProxyURL(session.ProxyURL.Value()); err != nil {
			return err
		}
	}

	if session.IsConnected() && !session.IsAuthenticated() {
		return ErrSessionNotConnected
	}

	return nil
}

func (s *DomainService) CanRegenerateApiKey(session *Session) bool {
	return !session.IsConnected()
}

func (s *DomainService) CanSetProxy(session *Session) bool {
	return !session.IsConnected()
}

func (s *DomainService) CanSubscribeToEvents(session *Session) bool {
	return true
}

func (s *DomainService) ValidateDeviceConnection(session *Session, deviceJID string) error {
	// A session can only be connected if it has a device JID
	if session.IsConnected() && deviceJID == "" {
		return ErrSessionCannotBeConnectedWithoutDevice
	}

	// A session without device JID cannot be connected
	if deviceJID == "" && session.IsConnected() {
		return ErrSessionCannotBeConnectedWithoutDevice
	}

	return nil
}

func ValidateSessionName(name string) error {
	if name == "" {
		return ErrInvalidSessionName
	}

	if len(name) < sessionNameMinLength {
		return ErrSessionNameTooShort
	}

	if len(name) > sessionNameMaxLength {
		return ErrSessionNameTooLong
	}

	if !sessionNameRegex.MatchString(name) {
		return ErrInvalidSessionNameChar
	}

	return nil
}

func ValidateSessionID(id string) error {
	if id == "" {
		return ErrInvalidSessionID
	}

	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(id) {
		return ErrInvalidSessionID
	}

	return nil
}

func ValidateProxyURL(proxyURL string) error {
	if proxyURL == "" {
		return nil // Empty proxy URL is valid (means no proxy)
	}

	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return ErrInvalidProxyURL
	}

	if parsedURL.Scheme == "" {
		return ErrInvalidProxyURL
	}

	if parsedURL.Host == "" {
		return ErrInvalidProxyURL
	}

	return nil
}

func ValidateSessionStatus(currentStatus, newStatus Status) error {
	validTransitions := map[Status][]Status{
		StatusDisconnected: {StatusConnecting},
		StatusConnecting:   {StatusConnected, StatusDisconnected, StatusError},
		StatusConnected:    {StatusDisconnected, StatusError},
		StatusError:        {StatusDisconnected, StatusConnecting},
	}

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return ErrInvalidSessionStatus
	}

	for _, allowed := range allowedStatuses {
		if newStatus == allowed {
			return nil
		}
	}

	return ErrInvalidSessionStatus
}
