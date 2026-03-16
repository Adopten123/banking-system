package domain

type UserCreatedEvent struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

func (e UserCreatedEvent) EventName() string { return "UserCreatedEvent" }