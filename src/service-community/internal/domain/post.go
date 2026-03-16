package domain

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID                 int64
	AuthorID           uuid.UUID
	TypeID             int32
	Content            string
	MediaAttachments   []byte
	RelatedAssetTicker *string
	Status             string
	LikesCount         int32
	CommentsCount      int32
	IsPinned           bool
	IsEdited           bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type CreatePostInput struct {
	AuthorID           uuid.UUID
	TypeID             int32
	Content            string
	MediaAttachments   []byte
	RelatedAssetTicker *string
	Status             string
}
