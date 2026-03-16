package domain

import "errors"

var (
	ErrEmptyContent = errors.New("post content cannot be empty")
	ErrInvalidType  = errors.New("invalid post type")
)