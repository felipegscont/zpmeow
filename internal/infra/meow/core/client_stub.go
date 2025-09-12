package core

// Temporary stub for MeowClient to make build work
type MeowClient struct {
	SessionID string
}

// NewMeowClient creates a new meow client
func NewMeowClient(sessionID string) *MeowClient {
	return &MeowClient{
		SessionID: sessionID,
	}
}

// Placeholder methods
func (c *MeowClient) Connect() error { return nil }
func (c *MeowClient) Disconnect() error { return nil }
func (c *MeowClient) IsConnected() bool { return false }
func (c *MeowClient) GetQRCode() (string, error) { return "", nil }
func (c *MeowClient) PairPhone(phoneNumber string) (string, error) { return "", nil }
