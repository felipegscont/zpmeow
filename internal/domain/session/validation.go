package session

import (
	"net/url"
	"regexp"
	"strings"
)


type ValidationService struct {
	phoneRegex    *regexp.Regexp
	sanitizeRegex *regexp.Regexp
}


func NewValidationService() *ValidationService {
	return &ValidationService{
		phoneRegex:    regexp.MustCompile(`\D`), // Non-digit characters
		sanitizeRegex: regexp.MustCompile(`[\x00-\x1f\x7f]`), // Control characters
	}
}




func (v *ValidationService) ValidatePhoneNumber(phone string) error {
	if phone == "" {
		return NewDomainError("phone number cannot be empty")
	}


	cleaned := v.phoneRegex.ReplaceAllString(phone, "")
	

	if len(cleaned) < 10 || len(cleaned) > 15 {
		return NewDomainError("phone number must be between 10 and 15 digits")
	}
	
	return nil
}




func (v *ValidationService) ValidateSessionName(name string) error {
	if name == "" {
		return ErrInvalidSessionName
	}


	if len(name) < 3 {
		return ErrSessionNameTooShort
	}
	if len(name) > 50 {
		return ErrSessionNameTooLong
	}




	for i, r := range name {
		if !v.isValidSessionNameChar(r) {
			return ErrInvalidSessionNameChar
		}


		if i == 0 || i == len(name)-1 {
			if r == '-' || r == '_' {
				return ErrInvalidSessionNameFormat
			}
		}
	}


	if v.isReservedSessionName(name) {
		return ErrReservedSessionName
	}

	return nil
}


func (v *ValidationService) isValidSessionNameChar(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '-' || r == '_'
}


func (v *ValidationService) isReservedSessionName(name string) bool {
	reserved := map[string]bool{
		"api": true, "admin": true, "root": true, "system": true, "config": true, "settings": true,
		"health": true, "ping": true, "status": true, "info": true, "debug": true, "test": true,
		"sessions": true, "session": true, "create": true, "list": true, "delete": true,
		"connect": true, "disconnect": true, "logout": true, "qr": true, "pair": true,
		"send": true, "chat": true, "group": true, "contact": true, "media": true, "file": true,
		"swagger": true, "docs": true, "documentation": true,
	}

	return reserved[strings.ToLower(name)]
}




func (v *ValidationService) ValidateProxyURL(proxyURL string) error {
	if proxyURL == "" {
		return nil // Empty is valid (no proxy)
	}


	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return NewDomainError("invalid proxy URL format")
	}


	supportedSchemes := map[string]bool{
		"http":   true,
		"https":  true,
		"socks5": true,
	}

	if !supportedSchemes[parsedURL.Scheme] {
		return NewDomainError("proxy URL must use http, https, or socks5 scheme")
	}


	if parsedURL.Host == "" {
		return NewDomainError("proxy URL must include a host")
	}

	return nil
}




func (v *ValidationService) ValidateSessionID(id string) error {
	if strings.TrimSpace(id) == "" {
		return ErrInvalidSessionID
	}
	return nil
}


func (v *ValidationService) ValidateSessionIDOrName(idOrName string) error {
	if strings.TrimSpace(idOrName) == "" {
		return ErrInvalidSessionID // Use ID error as it's more generic
	}
	return nil
}




func (v *ValidationService) SanitizeString(input string) string {

	cleaned := v.sanitizeRegex.ReplaceAllString(input, "")
	return strings.TrimSpace(cleaned)
}


var DefaultValidationService = NewValidationService()




func ValidatePhoneNumber(phone string) error {
	return DefaultValidationService.ValidatePhoneNumber(phone)
}


func ValidateSessionName(name string) error {
	return DefaultValidationService.ValidateSessionName(name)
}


func ValidateProxyURL(proxyURL string) error {
	return DefaultValidationService.ValidateProxyURL(proxyURL)
}


func ValidateSessionID(id string) error {
	return DefaultValidationService.ValidateSessionID(id)
}


func ValidateSessionIDOrName(idOrName string) error {
	return DefaultValidationService.ValidateSessionIDOrName(idOrName)
}


func SanitizeString(input string) string {
	return DefaultValidationService.SanitizeString(input)
}




func IsValidPhoneNumber(phone string) bool {
	return ValidatePhoneNumber(phone) == nil
}


func IsValidProxyURL(proxyURL string) bool {
	return ValidateProxyURL(proxyURL) == nil
}
