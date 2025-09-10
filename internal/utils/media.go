package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/vincent-petithory/dataurl"
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
	// Find "base64" in any position after the MIME type
	hasBase64 := false
	for i := 1; i < len(parts); i++ {
		if strings.TrimSpace(parts[i]) == "base64" {
			hasBase64 = true
			break
		}
	}
	if !hasBase64 {
		return nil, "", fmt.Errorf("only base64 encoding is supported")
	}

	// Decode base64 data
	data, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode base64 data: %w", err)
	}

	return data, mimeType, nil
}

// DecodeAudioDataURL decodes audio data URL using the same method as wuzapi
// This function specifically handles audio/ogg data URLs
func DecodeAudioDataURL(dataURI string) ([]byte, error) {
	// Validate that it starts with "data:audio/ogg" like wuzapi does
	if len(dataURI) < 14 || dataURI[0:14] != "data:audio/ogg" {
		return nil, fmt.Errorf("audio data should start with \"data:audio/ogg\"")
	}

	// Use dataurl library to decode like wuzapi
	dataURL, err := dataurl.DecodeString(dataURI)
	if err != nil {
		return nil, fmt.Errorf("could not decode base64 encoded data from payload: %w", err)
	}

	return dataURL.Data, nil
}

// NormalizeMimeType normalizes MIME types to be compatible with WhatsApp
// This function ensures that any MIME type is converted to a format that whatsmeow accepts
func NormalizeMimeType(mimeType, mediaType string) string {
	// Remove any parameters from MIME type for comparison
	baseMimeType := strings.Split(mimeType, ";")[0]
	baseMimeType = strings.TrimSpace(strings.ToLower(baseMimeType))

	switch mediaType {
	case "audio":
		// For audio, always use OGG with Opus codec for best WhatsApp compatibility
		return "audio/ogg; codecs=opus"

	case "image":
		// Map various image formats to WhatsApp-supported ones
		switch baseMimeType {
		case "image/jpg", "image/jpeg":
			return "image/jpeg"
		case "image/png":
			return "image/png"
		case "image/gif":
			return "image/gif"
		case "image/webp":
			return "image/webp"
		case "image/bmp", "image/x-ms-bmp":
			return "image/jpeg" // Convert BMP to JPEG for better compatibility
		default:
			return "image/jpeg" // Default fallback for unknown image types
		}

	case "video":
		// Map various video formats to WhatsApp-supported ones
		switch baseMimeType {
		case "video/mp4":
			return "video/mp4"
		case "video/avi", "video/x-msvideo":
			return "video/mp4" // Convert AVI to MP4 for better compatibility
		case "video/mov", "video/quicktime":
			return "video/mp4" // Convert MOV to MP4 for better compatibility
		case "video/webm":
			return "video/webm"
		case "video/3gp", "video/3gpp":
			return "video/3gp"
		default:
			return "video/mp4" // Default fallback for unknown video types
		}

	case "document":
		// Keep original MIME type for documents as WhatsApp supports many formats
		// Just clean it up by removing parameters
		return baseMimeType

	case "sticker":
		// Stickers should be WebP for best compatibility
		if baseMimeType == "image/png" {
			return "image/png"
		}
		return "image/webp" // Default to WebP for stickers

	default:
		// For unknown media types, return the cleaned base MIME type
		return baseMimeType
	}
}

// ValidateAndNormalizeMimeType validates basic MIME type format and normalizes it for WhatsApp compatibility
// This function is permissive and tries to make any MIME type work with WhatsApp
func ValidateAndNormalizeMimeType(mimeType, mediaType string) (string, error) {
	// Basic validation - ensure it's not empty and has a reasonable format
	if mimeType == "" {
		return "", fmt.Errorf("MIME type cannot be empty")
	}

	// Check if it has the basic format (type/subtype)
	if !strings.Contains(mimeType, "/") {
		return "", fmt.Errorf("invalid MIME type format: %s", mimeType)
	}

	// Get the base MIME type (before any parameters)
	baseMimeType := strings.Split(mimeType, ";")[0]
	baseMimeType = strings.TrimSpace(strings.ToLower(baseMimeType))

	// Basic category validation
	switch mediaType {
	case "image":
		if !strings.HasPrefix(baseMimeType, "image/") {
			return "", fmt.Errorf("MIME type must start with 'image/' for images, got: %s", mimeType)
		}
	case "audio":
		if !strings.HasPrefix(baseMimeType, "audio/") {
			return "", fmt.Errorf("MIME type must start with 'audio/' for audio, got: %s", mimeType)
		}
	case "video":
		if !strings.HasPrefix(baseMimeType, "video/") {
			return "", fmt.Errorf("MIME type must start with 'video/' for video, got: %s", mimeType)
		}
	case "sticker":
		if !strings.HasPrefix(baseMimeType, "image/") {
			return "", fmt.Errorf("MIME type must start with 'image/' for stickers, got: %s", mimeType)
		}
	case "document":
		// Documents can have various MIME types, so we're more permissive
		// Just ensure it's not an image, audio, or video type being misclassified
		if strings.HasPrefix(baseMimeType, "image/") ||
		   strings.HasPrefix(baseMimeType, "audio/") ||
		   strings.HasPrefix(baseMimeType, "video/") {
			return "", fmt.Errorf("MIME type %s should not be classified as document", mimeType)
		}
	}

	// Normalize the MIME type for WhatsApp compatibility
	normalizedMimeType := NormalizeMimeType(mimeType, mediaType)
	return normalizedMimeType, nil
}

// DecodeUniversalMedia decodes any media data URL and returns the data with normalized MIME type
// This function is designed to be as permissive as possible while ensuring WhatsApp compatibility
func DecodeUniversalMedia(dataURI, mediaType string) ([]byte, string, error) {
	// Basic validation
	if !strings.HasPrefix(dataURI, "data:") {
		return nil, "", fmt.Errorf("invalid data URI format: must start with 'data:'")
	}

	// For audio, use the specific wuzapi-compatible decoder
	if mediaType == "audio" {
		// Check if it's audio/ogg format (required by wuzapi compatibility)
		if len(dataURI) >= 14 && dataURI[0:14] == "data:audio/ogg" {
			data, err := DecodeAudioDataURL(dataURI)
			if err != nil {
				return nil, "", err
			}
			return data, "audio/ogg; codecs=opus", nil
		} else {
			// For non-OGG audio, try to decode with fallback method
			data, _, err := DecodeDataURLFallback(dataURI)
			if err != nil {
				return nil, "", fmt.Errorf("could not decode audio data URL: %w", err)
			}
			// Always return OGG format for audio compatibility
			return data, "audio/ogg; codecs=opus", nil
		}
	}

	// For all other media types, try the dataurl library first, then fallback
	dataURL, err := dataurl.DecodeString(dataURI)
	if err != nil {
		// If dataurl library fails, try our fallback method
		data, mimeType, err := DecodeDataURLFallback(dataURI)
		if err != nil {
			return nil, "", fmt.Errorf("could not decode data URL: %w", err)
		}

		// Validate and normalize the MIME type
		normalizedMimeType, err := ValidateAndNormalizeMimeType(mimeType, mediaType)
		if err != nil {
			return nil, "", fmt.Errorf("MIME type validation failed: %w", err)
		}

		return data, normalizedMimeType, nil
	}

	// Validate and normalize the MIME type
	normalizedMimeType, err := ValidateAndNormalizeMimeType(dataURL.ContentType(), mediaType)
	if err != nil {
		return nil, "", fmt.Errorf("MIME type validation failed: %w", err)
	}

	return dataURL.Data, normalizedMimeType, nil
}

// DecodeDataURLFallback provides a robust fallback method for decoding data URLs
// This function is more permissive than the standard dataurl library and handles non-standard formats
func DecodeDataURLFallback(dataURI string) ([]byte, string, error) {
	// Remove "data:" prefix
	if !strings.HasPrefix(dataURI, "data:") {
		return nil, "", fmt.Errorf("invalid data URI: must start with 'data:'")
	}

	content := dataURI[5:] // Remove "data:"

	// Find the comma that separates metadata from data
	commaIndex := strings.Index(content, ",")
	if commaIndex == -1 {
		return nil, "", fmt.Errorf("invalid data URI: missing comma separator")
	}

	metadata := content[:commaIndex]
	encodedData := content[commaIndex+1:]

	// Parse metadata to extract MIME type and encoding
	var mimeType string
	var encoding string

	// Split metadata by semicolon
	parts := strings.Split(metadata, ";")
	if len(parts) > 0 && parts[0] != "" {
		rawMimeType := strings.TrimSpace(parts[0])
		mimeType = NormalizeMimeTypeFromRaw(rawMimeType)
	} else {
		mimeType = "text/plain" // Default MIME type
	}

	// Look for encoding in the parts
	encoding = "base64" // Default to base64
	for _, part := range parts[1:] {
		part = strings.TrimSpace(part)
		if part == "base64" {
			encoding = "base64"
			break
		} else if strings.Contains(part, "charset=") {
			// Handle charset parameters but keep base64 encoding
			continue
		}
	}

	// Decode the data
	var data []byte
	var err error

	if encoding == "base64" {
		// Clean up the base64 string - remove any whitespace
		cleanedData := strings.ReplaceAll(encodedData, " ", "")
		cleanedData = strings.ReplaceAll(cleanedData, "\n", "")
		cleanedData = strings.ReplaceAll(cleanedData, "\r", "")
		cleanedData = strings.ReplaceAll(cleanedData, "\t", "")

		data, err = base64.StdEncoding.DecodeString(cleanedData)
		if err != nil {
			return nil, "", fmt.Errorf("failed to decode base64 data: %w", err)
		}
	} else {
		// For non-base64 encoding, treat as URL-encoded
		data = []byte(encodedData)
	}

	return data, mimeType, nil
}

// NormalizeMimeTypeFromRaw converts non-standard MIME types to standard ones
// This handles cases like @file/png, @file/jpeg, etc. from n8n and other tools
func NormalizeMimeTypeFromRaw(rawMimeType string) string {
	// Convert to lowercase for comparison
	lower := strings.ToLower(strings.TrimSpace(rawMimeType))

	// Handle n8n's @file/ format
	if strings.HasPrefix(lower, "@file/") {
		fileType := lower[6:] // Remove "@file/"
		switch fileType {
		case "png":
			return "image/png"
		case "jpg", "jpeg":
			return "image/jpeg"
		case "gif":
			return "image/gif"
		case "webp":
			return "image/webp"
		case "bmp":
			return "image/bmp"
		case "svg":
			return "image/svg+xml"
		case "pdf":
			return "application/pdf"
		case "doc":
			return "application/msword"
		case "docx":
			return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		case "xls":
			return "application/vnd.ms-excel"
		case "xlsx":
			return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		case "ppt":
			return "application/vnd.ms-powerpoint"
		case "pptx":
			return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
		case "txt":
			return "text/plain"
		case "csv":
			return "text/csv"
		case "zip":
			return "application/zip"
		case "rar":
			return "application/x-rar-compressed"
		case "mp4":
			return "video/mp4"
		case "avi":
			return "video/avi"
		case "mov":
			return "video/mov"
		case "webm":
			return "video/webm"
		case "3gp":
			return "video/3gp"
		case "mp3":
			return "audio/mpeg"
		case "wav":
			return "audio/wav"
		case "ogg":
			return "audio/ogg"
		case "aac":
			return "audio/aac"
		default:
			// If we don't recognize the file type, try to guess based on common patterns
			if strings.Contains(fileType, "image") || strings.Contains(fileType, "img") {
				return "image/jpeg" // Default to JPEG for unknown image types
			} else if strings.Contains(fileType, "video") || strings.Contains(fileType, "vid") {
				return "video/mp4" // Default to MP4 for unknown video types
			} else if strings.Contains(fileType, "audio") || strings.Contains(fileType, "sound") {
				return "audio/mpeg" // Default to MP3 for unknown audio types
			} else {
				return "application/octet-stream" // Generic binary type
			}
		}
	}

	// Handle other non-standard formats
	if strings.HasPrefix(lower, "file/") {
		// Similar to @file/ but without @
		return NormalizeMimeTypeFromRaw("@" + lower)
	}

	// If it's already a standard MIME type, return as-is
	if strings.Contains(lower, "/") {
		return lower
	}

	// If it's just a file extension, try to convert it
	switch lower {
	case "png":
		return "image/png"
	case "jpg", "jpeg":
		return "image/jpeg"
	case "gif":
		return "image/gif"
	case "webp":
		return "image/webp"
	case "pdf":
		return "application/pdf"
	case "mp4":
		return "video/mp4"
	case "mp3":
		return "audio/mpeg"
	case "ogg":
		return "audio/ogg"
	default:
		return "application/octet-stream"
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

// ProcessUnifiedMedia processes media from different sources: base64, URL, or form-data
// Returns the media data, MIME type, and any error
func ProcessUnifiedMedia(ctx context.Context, media string, file *multipart.FileHeader, mediaType string) ([]byte, string, error) {
	// If file is provided (form-data), process it first
	if file != nil {
		return processFormDataMedia(file, mediaType)
	}

	// If media is empty, return error
	if media == "" {
		return nil, "", fmt.Errorf("media parameter is required")
	}

	// Check if it's a URL
	if isValidURL(media) {
		return downloadMediaFromURL(ctx, media, mediaType)
	}

	// Otherwise, treat as base64 data URL
	return DecodeUniversalMedia(media, mediaType)
}

// processFormDataMedia processes media uploaded via multipart/form-data
func processFormDataMedia(fileHeader *multipart.FileHeader, mediaType string) ([]byte, string, error) {
	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// Read file content
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read uploaded file: %w", err)
	}

	// Detect MIME type from file content
	mimeType := http.DetectContentType(data)

	// If MIME type detection fails, try to get from filename
	if mimeType == "application/octet-stream" {
		ext := filepath.Ext(fileHeader.Filename)
		if ext != "" {
			detectedMime := mime.TypeByExtension(ext)
			if detectedMime != "" {
				mimeType = detectedMime
			}
		}
	}

	// Validate and normalize the MIME type
	normalizedMimeType, err := ValidateAndNormalizeMimeType(mimeType, mediaType)
	if err != nil {
		return nil, "", fmt.Errorf("MIME type validation failed: %w", err)
	}

	return data, normalizedMimeType, nil
}

// downloadMediaFromURL downloads media from a public URL
func downloadMediaFromURL(ctx context.Context, mediaURL string, mediaType string) ([]byte, string, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", mediaURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set user agent
	req.Header.Set("User-Agent", "ZpMeow/1.0")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download media from URL: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download media: HTTP %d", resp.StatusCode)
	}

	// Check content length (limit to 100MB)
	if resp.ContentLength > 100*1024*1024 {
		return nil, "", fmt.Errorf("file too large: %d bytes (max 100MB)", resp.ContentLength)
	}

	// Read response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Get MIME type from response header or detect from content
	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = http.DetectContentType(data)
	}

	// Clean up MIME type (remove charset and other parameters)
	if strings.Contains(mimeType, ";") {
		mimeType = strings.Split(mimeType, ";")[0]
	}
	mimeType = strings.TrimSpace(mimeType)

	// Validate and normalize the MIME type
	normalizedMimeType, err := ValidateAndNormalizeMimeType(mimeType, mediaType)
	if err != nil {
		return nil, "", fmt.Errorf("MIME type validation failed: %w", err)
	}

	return data, normalizedMimeType, nil
}

// isValidURL checks if a string is a valid HTTP/HTTPS URL
func isValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

// ValidateMediaType validates that the mediaType is one of the supported types
func ValidateMediaType(mediaType string) error {
	validTypes := map[string]bool{
		"image":    true,
		"audio":    true,
		"document": true,
		"video":    true,
	}

	if !validTypes[mediaType] {
		return fmt.Errorf("invalid mediaType '%s': must be one of: image, audio, document, video", mediaType)
	}

	return nil
}
