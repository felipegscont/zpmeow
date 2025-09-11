package meow

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	waTypes "go.mau.fi/whatsmeow/types"
)


type Validator struct {
	phoneRegex   *regexp.Regexp
	sessionRegex *regexp.Regexp
	jidRegex     *regexp.Regexp
}


func NewValidator() *Validator {
	return &Validator{
		phoneRegex:   regexp.MustCompile(PhoneNumberPattern),
		sessionRegex: regexp.MustCompile(SessionIDPattern),
		jidRegex:     regexp.MustCompile(JIDPattern),
	}
}


func (v *Validator) ValidateSessionID(sessionID string) error {
	if sessionID == "" {
		return errors.New(ErrEmptySessionID)
	}
	
	if len(sessionID) > 100 {
		return fmt.Errorf("session ID too long (max 100 characters)")
	}
	
	if !v.sessionRegex.MatchString(sessionID) {
		return fmt.Errorf("session ID contains invalid characters (only alphanumeric, underscore, and hyphen allowed)")
	}
	
	return nil
}


func (v *Validator) ValidateDeviceJID(deviceJID string) error {
	if deviceJID == "" {
		return errors.New(ErrEmptyDeviceJID)
	}
	
	if len(deviceJID) > 200 {
		return fmt.Errorf("device JID too long (max 200 characters)")
	}
	

	_, err := waTypes.ParseJID(deviceJID)
	if err != nil {
		return fmt.Errorf("invalid device JID format: %w", err)
	}
	
	return nil
}


func (v *Validator) ValidatePhoneNumber(phoneNumber string) error {
	if phoneNumber == "" {
		return fmt.Errorf("phone number cannot be empty")
	}
	

	cleaned := strings.ReplaceAll(phoneNumber, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "-", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")
	
	if !v.phoneRegex.MatchString(cleaned) {
		return errors.New(ErrInvalidPhoneNumber)
	}
	
	return nil
}


func (v *Validator) ValidateJID(jidStr string) error {
	if jidStr == "" {
		return fmt.Errorf("JID cannot be empty")
	}
	
	_, err := waTypes.ParseJID(jidStr)
	if err != nil {
		return fmt.Errorf(ErrInvalidJID + ": %w", err)
	}
	
	return nil
}


func (v *Validator) ValidateClientConnection(client *MeowClient) error {
	if client == nil {
		return errors.New(ErrClientNotFound)
	}






	return nil
}


func (v *Validator) ValidateMessageContent(content string) error {
	if content == "" {
		return fmt.Errorf("message content cannot be empty")
	}
	
	if len(content) > 65536 { // 64KB limit
		return fmt.Errorf(ErrMessageTooLarge + " (max 64KB)")
	}
	
	return nil
}


func (v *Validator) ValidateMediaSize(data []byte, mediaType string) error {
	if len(data) == 0 {
		return fmt.Errorf("media data cannot be empty")
	}
	
	size := len(data)
	
	switch mediaType {
	case "image":
		if size > MaxImageSize {
			return fmt.Errorf("image too large (max %d bytes)", MaxImageSize)
		}
	case "video":
		if size > MaxVideoSize {
			return fmt.Errorf("video too large (max %d bytes)", MaxVideoSize)
		}
	case "audio":
		if size > MaxAudioSize {
			return fmt.Errorf("audio too large (max %d bytes)", MaxAudioSize)
		}
	case "document":
		if size > MaxDocumentSize {
			return fmt.Errorf("document too large (max %d bytes)", MaxDocumentSize)
		}
	default:
		return fmt.Errorf("%s: %s", ErrUnsupportedMediaType, mediaType)
	}
	
	return nil
}


func (v *Validator) ValidateMimeType(mimeType, expectedType string) error {
	if mimeType == "" {
		return fmt.Errorf("MIME type cannot be empty")
	}
	
	if !strings.HasPrefix(mimeType, expectedType+"/") {
		return fmt.Errorf("invalid MIME type for %s: %s", expectedType, mimeType)
	}
	
	return nil
}


func (v *Validator) ValidateFileName(fileName string) error {
	if fileName == "" {
		return fmt.Errorf("file name cannot be empty")
	}
	
	if len(fileName) > 255 {
		return fmt.Errorf("file name too long (max 255 characters)")
	}
	

	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		if strings.Contains(fileName, char) {
			return fmt.Errorf("file name contains invalid character: %s", char)
		}
	}
	
	return nil
}


func (v *Validator) ValidateCoordinates(latitude, longitude float64) error {
	if latitude < -90 || latitude > 90 {
		return fmt.Errorf("invalid latitude: %f (must be between -90 and 90)", latitude)
	}
	
	if longitude < -180 || longitude > 180 {
		return fmt.Errorf("invalid longitude: %f (must be between -180 and 180)", longitude)
	}
	
	return nil
}


func (v *Validator) ValidatePollOptions(options []string) error {
	if len(options) < 2 {
		return fmt.Errorf("poll must have at least 2 options")
	}
	
	if len(options) > 12 {
		return fmt.Errorf("poll cannot have more than 12 options")
	}
	
	for i, option := range options {
		if option == "" {
			return fmt.Errorf("poll option %d cannot be empty", i+1)
		}
		
		if len(option) > 100 {
			return fmt.Errorf("poll option %d too long (max 100 characters)", i+1)
		}
	}
	
	return nil
}


func (v *Validator) ValidateVCard(vcard string) error {
	if vcard == "" {
		return fmt.Errorf("vCard cannot be empty")
	}
	
	if !strings.HasPrefix(vcard, "BEGIN:VCARD") {
		return fmt.Errorf("invalid vCard format: must start with BEGIN:VCARD")
	}
	
	if !strings.HasSuffix(vcard, "END:VCARD") {
		return fmt.Errorf("invalid vCard format: must end with END:VCARD")
	}
	
	return nil
}


var DefaultValidator = NewValidator()


func ValidateSessionID(sessionID string) error {
	return DefaultValidator.ValidateSessionID(sessionID)
}

func ValidateDeviceJID(deviceJID string) error {
	return DefaultValidator.ValidateDeviceJID(deviceJID)
}

func ValidatePhoneNumber(phoneNumber string) error {
	return DefaultValidator.ValidatePhoneNumber(phoneNumber)
}

func ValidateJID(jidStr string) error {
	return DefaultValidator.ValidateJID(jidStr)
}

func ValidateClientConnection(client *MeowClient) error {
	return DefaultValidator.ValidateClientConnection(client)
}

func ValidateMessageContent(content string) error {
	return DefaultValidator.ValidateMessageContent(content)
}

func ValidateMediaSize(data []byte, mediaType string) error {
	return DefaultValidator.ValidateMediaSize(data, mediaType)
}

func ValidateMimeType(mimeType, expectedType string) error {
	return DefaultValidator.ValidateMimeType(mimeType, expectedType)
}

func ValidateFileName(fileName string) error {
	return DefaultValidator.ValidateFileName(fileName)
}

func ValidateCoordinates(latitude, longitude float64) error {
	return DefaultValidator.ValidateCoordinates(latitude, longitude)
}

func ValidatePollOptions(options []string) error {
	return DefaultValidator.ValidatePollOptions(options)
}

func ValidateVCard(vcard string) error {
	return DefaultValidator.ValidateVCard(vcard)
}
