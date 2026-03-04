package domain

import "errors"

var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrInternal          = errors.New("internal server error")
	ErrAccountInactive   = errors.New("account is not active")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrAccountHasBalance = errors.New("account balance is not zero")
	ErrInvalidDepositAmount = errors.New("deposit amount must be greater than zero")
	ErrInvalidAmountFormat = errors.New("invalid amount format")
)
