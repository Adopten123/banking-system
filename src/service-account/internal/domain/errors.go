package domain

import "errors"

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrInternal        = errors.New("internal server error")
	ErrAccountInactive = errors.New("account is not active")
	ErrInsufficientFunds = errors.New("insufficient funds")
)
