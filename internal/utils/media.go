package utils

import (
	"encoding/base64"
	"fmt"
	"mime"
	"strings"
)

// DecodeBase64Media decodes base64 media data from data URI format
// Supports formats like: data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD...
func DecodeBase64Media(dataURI string) ([]byte, string, error) {
	if !strings.HasPrefix(dataURI, "data:") {
		return nil, "", fmt.Errorf("invalid data URI format")
	}

	// Remove "data:" prefix
	dataURI = dataURI[5:]

	// Find the comma that separates metadata from data
	commaIndex := strings.Index(dataURI, ",")
	if commaIndex == -1 {
		return nil, "", fmt.Errorf("invalid data URI format: missing comma")
	}

	// Extract metadata and data parts
	metadata := dataURI[:commaIndex]
	encodedData := dataURI[commaIndex+1:]

	// Parse metadata to extract MIME type
	var mimeType string
	parts := strings.Split(metadata, ";")
	if len(parts) > 0 {
		mimeType = parts[0]
	}

	// Validate that it's base64 encoded
	if len(parts) < 2 || parts[1] != "base64" {
		return nil, "", fmt.Errorf("only base64 encoding is supported")
	}

	// Decode base64 data
	data, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode base64 data: %w", err)
	}

	return data, mimeType, nil
}

// ValidateMimeType validates if the MIME type is supported for the given media type
func ValidateMimeType(mimeType, mediaType string) error {
	switch mediaType {
	case "image":
		if !strings.HasPrefix(mimeType, "image/") {
			return fmt.Errorf("invalid MIME type for image: %s", mimeType)
		}
		// Common image types
		validTypes := []string{"image/jpeg", "image/png", "image/gif", "image/webp", "image/bmp"}
		for _, validType := range validTypes {
			if mimeType == validType {
				return nil
			}
		}
		return fmt.Errorf("unsupported image MIME type: %s", mimeType)

	case "audio":
		if !strings.HasPrefix(mimeType, "audio/") {
			return fmt.Errorf("invalid MIME type for audio: %s", mimeType)
		}
		// Common audio types
		validTypes := []string{"audio/mpeg", "audio/mp4", "audio/ogg", "audio/wav", "audio/aac", "audio/webm"}
		for _, validType := range validTypes {
			if mimeType == validType {
				return nil
			}
		}
		return fmt.Errorf("unsupported audio MIME type: %s", mimeType)

	case "video":
		if !strings.HasPrefix(mimeType, "video/") {
			return fmt.Errorf("invalid MIME type for video: %s", mimeType)
		}
		// Common video types
		validTypes := []string{"video/mp4", "video/avi", "video/mov", "video/webm", "video/3gp"}
		for _, validType := range validTypes {
			if mimeType == validType {
				return nil
			}
		}
		return fmt.Errorf("unsupported video MIME type: %s", mimeType)

	case "document":
		// Documents can have various MIME types
		validTypes := []string{
			"application/pdf",
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"application/vnd.ms-excel",
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			"application/vnd.ms-powerpoint",
			"application/vnd.openxmlformats-officedocument.presentationml.presentation",
			"text/plain",
			"text/csv",
			"application/zip",
			"application/x-rar-compressed",
		}
		for _, validType := range validTypes {
			if mimeType == validType {
				return nil
			}
		}
		return fmt.Errorf("unsupported document MIME type: %s", mimeType)

	case "sticker":
		// Stickers are usually WebP images
		if mimeType != "image/webp" && mimeType != "image/png" {
			return fmt.Errorf("stickers must be WebP or PNG format, got: %s", mimeType)
		}
		return nil

	default:
		return fmt.Errorf("unknown media type: %s", mediaType)
	}
}

// GetFileExtension returns the file extension for a given MIME type
func GetFileExtension(mimeType string) string {
	extensions, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(extensions) == 0 {
		// Fallback for common types
		switch mimeType {
		case "image/jpeg":
			return ".jpg"
		case "image/png":
			return ".png"
		case "image/gif":
			return ".gif"
		case "image/webp":
			return ".webp"
		case "audio/mpeg":
			return ".mp3"
		case "audio/mp4":
			return ".m4a"
		case "audio/ogg":
			return ".ogg"
		case "video/mp4":
			return ".mp4"
		case "video/avi":
			return ".avi"
		case "application/pdf":
			return ".pdf"
		default:
			return ""
		}
	}
	return extensions[0]
}

// ValidateMediaSize validates if the media size is within acceptable limits
func ValidateMediaSize(data []byte, mediaType string) error {
	size := len(data)
	
	switch mediaType {
	case "image":
		if size > 16*1024*1024 { // 16MB
			return fmt.Errorf("image size too large: %d bytes (max 16MB)", size)
		}
	case "audio":
		if size > 16*1024*1024 { // 16MB
			return fmt.Errorf("audio size too large: %d bytes (max 16MB)", size)
		}
	case "video":
		if size > 64*1024*1024 { // 64MB
			return fmt.Errorf("video size too large: %d bytes (max 64MB)", size)
		}
	case "document":
		if size > 100*1024*1024 { // 100MB
			return fmt.Errorf("document size too large: %d bytes (max 100MB)", size)
		}
	case "sticker":
		if size > 1*1024*1024 { // 1MB
			return fmt.Errorf("sticker size too large: %d bytes (max 1MB)", size)
		}
	}
	
	if size == 0 {
		return fmt.Errorf("media data is empty")
	}
	
	return nil
}
