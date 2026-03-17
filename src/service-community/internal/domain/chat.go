package domain

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID            uuid.UUID `json:"id"`
	TypeID        int32     `json:"type_id"`
	Title         *string   `json:"title,omitempty"`
	AvatarURL     *string   `json:"avatar_url,omitempty"`
	LastMessageAt time.Time `json:"last_message_at"`
	CreatedAt     time.Time `json:"created_at"`
}

type Message struct {
	ID               int64     `json:"id"`
	ChatID           uuid.UUID `json:"chat_id"`
	SenderID         uuid.UUID `json:"sender_id"`
	ReplyToMessageID *int64    `json:"reply_to_message_id,omitempty"`
	Content          *string   `json:"content,omitempty"`
	MediaAttachments []byte    `json:"media_attachments,omitempty"`

	IsTransfer            bool       `json:"is_transfer"`
	TransferAmount        *string    `json:"transfer_amount,omitempty"`
	TransferCurrency      *string    `json:"transfer_currency,omitempty"`
	IdempotencyKey        *string    `json:"idempotency_key,omitempty"`
	TransferTransactionID *uuid.UUID `json:"transfer_transaction_id,omitempty"`
	TransferStatus        *string    `json:"transfer_status,omitempty"`

	IsEdited  bool       `json:"is_edited"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type CreateMessageInput struct {
	ChatID           uuid.UUID
	SenderID         uuid.UUID
	ReplyToMessageID *int64
	Content          *string

	IsTransfer       bool
	TransferAmount   *string
	TransferCurrency *string
	IdempotencyKey   *string
}

type CreateChatInput struct {
	TypeID    int32
	Title     *string
	AvatarURL *string
	MemberIDs []uuid.UUID
}
