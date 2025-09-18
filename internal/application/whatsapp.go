package application

type MessageResponse struct {
	ID        string
	Timestamp int64
	Status    string
}

type GroupResponse struct {
	GroupJID string
	Name     string
	Members  []string
}

type MessageInfo struct {
	ID        string
	FromJID   string
	ToJID     string
	Content   string
	Type      string
	Timestamp int64
	IsFromMe  bool
}
