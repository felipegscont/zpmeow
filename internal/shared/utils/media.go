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

func DecodeBase64Media(dataURI string) ([]byte, string, error) {
	if !strings.HasPrefix(dataURI, "data:") {
		return nil, "", fmt.Errorf("invalid data URI format")
	}

	dataURI = dataURI[5:]

	commaIndex := strings.Index(dataURI, ",")
	if commaIndex == -1 {
		return nil, "", fmt.Errorf("invalid data URI format: missing comma")
	}

	metadata := dataURI[:commaIndex]
	encodedData := dataURI[commaIndex+1:]

	var mimeType string
	parts := strings.Split(metadata, ";")
	if len(parts) > 0 {
		mimeType = parts[0]
	}

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

	data, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode base64 data: %w", err)
	}

	return data, mimeType, nil
}

func DecodeAudioDataURL(dataURI string) ([]byte, error) {

	if len(dataURI) < 14 || dataURI[0:14] != "data:audio/ogg" {
		return nil, fmt.Errorf("audio data should start with \"data:audio/ogg\"")
	}

	dataURL, err := dataurl.DecodeString(dataURI)
	if err != nil {
		return nil, fmt.Errorf("could not decode base64 encoded data from payload: %w", err)
	}

	return dataURL.Data, nil
}

func NormalizeMimeType(mimeType, mediaType string) string {

	baseMimeType := strings.Split(mimeType, ";")[0]
	baseMimeType = strings.TrimSpace(strings.ToLower(baseMimeType))

	switch mediaType {
	case "audio":

		return "audio/ogg; codecs=opus"

	case "image":

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

		return baseMimeType

	case "sticker":

		if baseMimeType == "image/png" {
			return "image/png"
		}
		return "image/webp" // Default to WebP for stickers

	default:

		return baseMimeType
	}
}

func ValidateAndNormalizeMimeType(mimeType, mediaType string) (string, error) {

	if mimeType == "" {
		return "", fmt.Errorf("MIME type cannot be empty")
	}

	if !strings.Contains(mimeType, "/") {
		return "", fmt.Errorf("invalid MIME type format: %s", mimeType)
	}

	baseMimeType := strings.Split(mimeType, ";")[0]
	baseMimeType = strings.TrimSpace(strings.ToLower(baseMimeType))

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

		if strings.HasPrefix(baseMimeType, "image/") ||
			strings.HasPrefix(baseMimeType, "audio/") ||
			strings.HasPrefix(baseMimeType, "video/") {
			return "", fmt.Errorf("MIME type %s should not be classified as document", mimeType)
		}
	}

	normalizedMimeType := NormalizeMimeType(mimeType, mediaType)
	return normalizedMimeType, nil
}

func DecodeUniversalMedia(dataURI, mediaType string) ([]byte, string, error) {

	if !strings.HasPrefix(dataURI, "data:") {
		return nil, "", fmt.Errorf("invalid data URI format: must start with 'data:'")
	}

	if mediaType == "audio" {

		if len(dataURI) >= 14 && dataURI[0:14] == "data:audio/ogg" {
			data, err := DecodeAudioDataURL(dataURI)
			if err != nil {
				return nil, "", err
			}
			return data, "audio/ogg; codecs=opus", nil
		} else {

			data, _, err := DecodeDataURLFallback(dataURI)
			if err != nil {
				return nil, "", fmt.Errorf("could not decode audio data URL: %w", err)
			}

			return data, "audio/ogg; codecs=opus", nil
		}
	}

	dataURL, err := dataurl.DecodeString(dataURI)
	if err != nil {

		data, mimeType, err := DecodeDataURLFallback(dataURI)
		if err != nil {
			return nil, "", fmt.Errorf("could not decode data URL: %w", err)
		}

		normalizedMimeType, err := ValidateAndNormalizeMimeType(mimeType, mediaType)
		if err != nil {
			return nil, "", fmt.Errorf("MIME type validation failed: %w", err)
		}

		return data, normalizedMimeType, nil
	}

	normalizedMimeType, err := ValidateAndNormalizeMimeType(dataURL.ContentType(), mediaType)
	if err != nil {
		return nil, "", fmt.Errorf("MIME type validation failed: %w", err)
	}

	return dataURL.Data, normalizedMimeType, nil
}

func DecodeDataURLFallback(dataURI string) ([]byte, string, error) {

	if !strings.HasPrefix(dataURI, "data:") {
		return nil, "", fmt.Errorf("invalid data URI: must start with 'data:'")
	}

	content := dataURI[5:] // Remove "data:"

	commaIndex := strings.Index(content, ",")
	if commaIndex == -1 {
		return nil, "", fmt.Errorf("invalid data URI: missing comma separator")
	}

	metadata := content[:commaIndex]
	encodedData := content[commaIndex+1:]

	var mimeType string
	var encoding string

	parts := strings.Split(metadata, ";")
	if len(parts) > 0 && parts[0] != "" {
		rawMimeType := strings.TrimSpace(parts[0])
		mimeType = NormalizeMimeTypeFromRaw(rawMimeType)
	} else {
		mimeType = "text/plain" // Default MIME type
	}

	encoding = "base64" // Default to base64
	for _, part := range parts[1:] {
		part = strings.TrimSpace(part)
		if part == "base64" {
			encoding = "base64"
			break
		} else if strings.Contains(part, "charset=") {

			continue
		}
	}

	var data []byte
	var err error

	if encoding == "base64" {

		cleanedData := strings.ReplaceAll(encodedData, " ", "")
		cleanedData = strings.ReplaceAll(cleanedData, "\n", "")
		cleanedData = strings.ReplaceAll(cleanedData, "\r", "")
		cleanedData = strings.ReplaceAll(cleanedData, "\t", "")

		data, err = base64.StdEncoding.DecodeString(cleanedData)
		if err != nil {
			return nil, "", fmt.Errorf("failed to decode base64 data: %w", err)
		}
	} else {

		data = []byte(encodedData)
	}

	return data, mimeType, nil
}

func NormalizeMimeTypeFromRaw(rawMimeType string) string {

	lower := strings.ToLower(strings.TrimSpace(rawMimeType))

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

	if strings.HasPrefix(lower, "file/") {

		return NormalizeMimeTypeFromRaw("@" + lower)
	}

	if strings.Contains(lower, "/") {
		return lower
	}

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

func GetFileExtension(mimeType string) string {
	extensions, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(extensions) == 0 {

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

func ProcessUnifiedMedia(ctx context.Context, media string, file *multipart.FileHeader, mediaType string) ([]byte, string, error) {

	if file != nil {
		return processFormDataMedia(file, mediaType)
	}

	if media == "" {
		return nil, "", fmt.Errorf("media parameter is required")
	}

	if isValidURL(media) {
		return downloadMediaFromURL(ctx, media, mediaType)
	}

	return DecodeUniversalMedia(media, mediaType)
}

func processFormDataMedia(fileHeader *multipart.FileHeader, mediaType string) ([]byte, string, error) {

	file, err := fileHeader.Open()
	if err != nil {
		return nil, "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read uploaded file: %w", err)
	}

	mimeType := http.DetectContentType(data)

	if mimeType == "application/octet-stream" {
		ext := filepath.Ext(fileHeader.Filename)
		if ext != "" {
			detectedMime := mime.TypeByExtension(ext)
			if detectedMime != "" {
				mimeType = detectedMime
			}
		}
	}

	normalizedMimeType, err := ValidateAndNormalizeMimeType(mimeType, mediaType)
	if err != nil {
		return nil, "", fmt.Errorf("MIME type validation failed: %w", err)
	}

	return data, normalizedMimeType, nil
}

func downloadMediaFromURL(ctx context.Context, mediaURL string, mediaType string) ([]byte, string, error) {

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", mediaURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "ZpMeow/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download media from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download media: HTTP %d", resp.StatusCode)
	}

	if resp.ContentLength > 100*1024*1024 {
		return nil, "", fmt.Errorf("file too large: %d bytes (max 100MB)", resp.ContentLength)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response body: %w", err)
	}

	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = http.DetectContentType(data)
	}

	if strings.Contains(mimeType, ";") {
		mimeType = strings.Split(mimeType, ";")[0]
	}
	mimeType = strings.TrimSpace(mimeType)

	normalizedMimeType, err := ValidateAndNormalizeMimeType(mimeType, mediaType)
	if err != nil {
		return nil, "", fmt.Errorf("MIME type validation failed: %w", err)
	}

	return data, normalizedMimeType, nil
}

func isValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

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
