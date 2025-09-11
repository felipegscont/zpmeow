package session

import (
	"context"
	"strings"
	"testing"
	"zpmeow/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) Create(ctx context.Context, session *Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) GetByID(ctx context.Context, id string) (*Session, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Session), args.Error(1)
}

func (m *MockSessionRepository) GetByName(ctx context.Context, name string) (*Session, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Session), args.Error(1)
}

func (m *MockSessionRepository) GetAll(ctx context.Context) ([]*Session, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*Session), args.Error(1)
}

func (m *MockSessionRepository) Update(ctx context.Context, session *Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSessionRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockSessionRepository) GetByStatus(ctx context.Context, status types.Status) ([]*Session, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]*Session), args.Error(1)
}


type MockWhatsAppService struct {
	mock.Mock
}

func (m *MockWhatsAppService) StartClient(sessionID string) error {
	args := m.Called(sessionID)
	return args.Error(0)
}

func (m *MockWhatsAppService) StopClient(sessionID string) error {
	args := m.Called(sessionID)
	return args.Error(0)
}

func (m *MockWhatsAppService) LogoutClient(sessionID string) error {
	args := m.Called(sessionID)
	return args.Error(0)
}

func (m *MockWhatsAppService) GetQRCode(sessionID string) (string, error) {
	args := m.Called(sessionID)
	return args.String(0), args.Error(1)
}

func (m *MockWhatsAppService) PairPhone(sessionID, phoneNumber string) (string, error) {
	args := m.Called(sessionID, phoneNumber)
	return args.String(0), args.Error(1)
}

func (m *MockWhatsAppService) IsClientConnected(sessionID string) bool {
	args := m.Called(sessionID)
	return args.Bool(0)
}

func (m *MockWhatsAppService) GetClientStatus(sessionID string) types.Status {
	args := m.Called(sessionID)
	return args.Get(0).(types.Status)
}

func (m *MockWhatsAppService) ConnectOnStartup(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}


func (m *MockWhatsAppService) DeleteMessage(ctx context.Context, sessionID, chatJID, messageID string, forEveryone bool) error {
	args := m.Called(ctx, sessionID, chatJID, messageID, forEveryone)
	return args.Error(0)
}

func (m *MockWhatsAppService) EditMessage(ctx context.Context, sessionID, chatJID, messageID, newText string) (*types.SendResponse, error) {
	args := m.Called(ctx, sessionID, chatJID, messageID, newText)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.SendResponse), args.Error(1)
}

func (m *MockWhatsAppService) DownloadMedia(ctx context.Context, sessionID, messageID string) ([]byte, string, error) {
	args := m.Called(ctx, sessionID, messageID)
	return args.Get(0).([]byte), args.String(1), args.Error(2)
}

func (m *MockWhatsAppService) ReactToMessage(ctx context.Context, sessionID, chatJID, messageID, emoji string) error {
	args := m.Called(ctx, sessionID, chatJID, messageID, emoji)
	return args.Error(0)
}

func TestGetSession_ByID(t *testing.T) {
	mockRepo := new(MockSessionRepository)
	mockWhatsApp := new(MockWhatsAppService)
	service := NewSessionService(mockRepo, mockWhatsApp)

	ctx := context.Background()
	sessionID := "test-session-id"
	expectedSession := &Session{
		ID:     sessionID,
		Name:   "Test Session",
		Status: types.StatusDisconnected,
	}


	mockRepo.On("GetByID", ctx, sessionID).Return(expectedSession, nil)


	result, err := service.GetSession(ctx, sessionID)


	assert.NoError(t, err)
	assert.Equal(t, expectedSession, result)
	mockRepo.AssertExpectations(t)
}

func TestGetSession_ByName(t *testing.T) {
	mockRepo := new(MockSessionRepository)
	mockWhatsApp := new(MockWhatsAppService)
	service := NewSessionService(mockRepo, mockWhatsApp)

	ctx := context.Background()
	sessionName := "Test Session"
	expectedSession := &Session{
		ID:     "test-session-id",
		Name:   sessionName,
		Status: types.StatusDisconnected,
	}


	mockRepo.On("GetByID", ctx, sessionName).Return(nil, ErrSessionNotFound)
	mockRepo.On("GetByName", ctx, sessionName).Return(expectedSession, nil)


	result, err := service.GetSession(ctx, sessionName)


	assert.NoError(t, err)
	assert.Equal(t, expectedSession, result)
	mockRepo.AssertExpectations(t)
}

func TestGetSession_NotFound(t *testing.T) {
	mockRepo := new(MockSessionRepository)
	mockWhatsApp := new(MockWhatsAppService)
	service := NewSessionService(mockRepo, mockWhatsApp)

	ctx := context.Background()
	sessionIDOrName := "non-existent"


	mockRepo.On("GetByID", ctx, sessionIDOrName).Return(nil, ErrSessionNotFound)
	mockRepo.On("GetByName", ctx, sessionIDOrName).Return(nil, ErrSessionNotFound)


	result, err := service.GetSession(ctx, sessionIDOrName)


	assert.Error(t, err)
	assert.Equal(t, ErrSessionNotFound, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateSession_ValidName(t *testing.T) {
	mockRepo := new(MockSessionRepository)
	mockWhatsApp := new(MockWhatsAppService)
	service := NewSessionService(mockRepo, mockWhatsApp)

	ctx := context.Background()
	validName := "MySession123"


	mockRepo.On("GetByName", ctx, validName).Return(nil, ErrSessionNotFound)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*session.Session")).Return(nil)


	result, err := service.CreateSession(ctx, validName)


	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, validName, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestCreateSession_InvalidNames(t *testing.T) {
	mockRepo := new(MockSessionRepository)
	mockWhatsApp := new(MockWhatsAppService)
	service := NewSessionService(mockRepo, mockWhatsApp)

	ctx := context.Background()

	testCases := []struct {
		name        string
		sessionName string
		expectedErr error
	}{
		{"empty name", "", ErrInvalidSessionName},
		{"too short", "ab", ErrSessionNameTooShort},
		{"too long", strings.Repeat("a", 51), ErrSessionNameTooLong},
		{"with spaces", "My Session", ErrInvalidSessionNameChar},
		{"with special chars", "session@123", ErrInvalidSessionNameChar},
		{"starts with hyphen", "-session", ErrInvalidSessionNameFormat},
		{"ends with underscore", "session_", ErrInvalidSessionNameFormat},
		{"reserved name", "admin", ErrReservedSessionName},
		{"reserved name case", "API", ErrReservedSessionName},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			result, err := service.CreateSession(ctx, tc.sessionName)


			assert.Error(t, err)
			assert.Equal(t, tc.expectedErr, err)
			assert.Nil(t, result)
		})
	}
}
