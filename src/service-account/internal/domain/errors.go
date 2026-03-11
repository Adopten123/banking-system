package domain

import "errors"

// ErrSystem - Problem is not related to business rules
var (
	ErrSystem = errors.New("internal server error")
)

// Account Errors
var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrAccountInactive   = errors.New("account is not active")
	ErrAccountHasBalance = errors.New("account balance is not zero")
)

// Transaction Errors
var (
	ErrInsufficientFunds    = errors.New("insufficient funds")
	ErrDuplicateTransaction = errors.New("duplicate transaction: idempotency key already exists")
)

// Validation Errors
var (
	ErrInvalidDepositAmount = errors.New("deposit amount must be greater than zero")
	ErrInvalidAmountFormat  = errors.New("invalid amount format")
)

// Cards Errors
var (
	ErrCardNotFound = errors.New("card not found")
	ErrCardBlocked  = errors.New("card is blocked")
	ErrInvalidCardStatus = errors.New("invalid card status")
)
