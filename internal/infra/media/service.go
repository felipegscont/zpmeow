package media

import (
	"context"
	"fmt"
	"io"
)

// Service provides media handling capabilities
type Service interface {
	// Upload media file and return URL
	UploadMedia(ctx context.Context, sessionID string, mediaType string, data io.Reader, filename string) (*UploadResult, error)

	// Download media from URL
	DownloadMedia(ctx context.Context, url string) (*DownloadResult, error)

	// Get media info
	GetMediaInfo(ctx context.Context, mediaID string) (*MediaInfo, error)

	// Delete media
	DeleteMedia(ctx context.Context, mediaID string) error
}

// UploadResult represents the result of a media upload
type UploadResult struct {
	MediaID  string `json:"mediaId"`
	URL      string `json:"url"`
	MimeType string `json:"mimeType"`
	Size     int64  `json:"size"`
	Filename string `json:"filename"`
}

// DownloadResult represents the result of a media download
type DownloadResult struct {
	Data     io.ReadCloser `json:"-"`
	MimeType string        `json:"mimeType"`
	Size     int64         `json:"size"`
	Filename string        `json:"filename"`
}

// MediaInfo represents media information
type MediaInfo struct {
	ID         string `json:"id"`
	URL        string `json:"url"`
	MimeType   string `json:"mimeType"`
	Size       int64  `json:"size"`
	Filename   string `json:"filename"`
	UploadedAt int64  `json:"uploadedAt"`
}

// ServiceStub is a stub implementation of the media service
type ServiceStub struct{}

// NewService creates a new media service stub
func NewService() Service {
	return &ServiceStub{}
}

// UploadMedia stub implementation
func (s *ServiceStub) UploadMedia(ctx context.Context, sessionID string, mediaType string, data io.Reader, filename string) (*UploadResult, error) {
	return nil, fmt.Errorf("UploadMedia not implemented yet")
}

// DownloadMedia stub implementation
func (s *ServiceStub) DownloadMedia(ctx context.Context, url string) (*DownloadResult, error) {
	return nil, fmt.Errorf("DownloadMedia not implemented yet")
}

// GetMediaInfo stub implementation
func (s *ServiceStub) GetMediaInfo(ctx context.Context, mediaID string) (*MediaInfo, error) {
	return nil, fmt.Errorf("GetMediaInfo not implemented yet")
}

// DeleteMedia stub implementation
func (s *ServiceStub) DeleteMedia(ctx context.Context, mediaID string) error {
	return fmt.Errorf("DeleteMedia not implemented yet")
}
