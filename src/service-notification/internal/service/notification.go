package service

import "log"

// EventHandler - function which handle raw JSON
type EventHandler func(payload []byte) error

type NotificationService struct {
	// Event type -> handler func
	handlers map[string]EventHandler
}

func NewNotificationService() *NotificationService {
	svc := &NotificationService{
		handlers: make(map[string]EventHandler),
	}
	svc.registerHandlers()
	return svc
}

func (s *NotificationService) registerHandlers() {
	// Account events
	s.handlers["AccountCreatedEvent"] = s.handleAccountCreated
	s.handlers["AccountStatusChangedEvent"] = s.handleAccountStatusChanged
	s.handlers["CreditLimitChangedEvent"] = s.handleCreditLimitChanged

	s.handlers["TransferCreatedEvent"] = s.handleTransferCreated
	s.handlers["DepositCompletedEvent"] = s.handleDepositCompleted
	s.handlers["WithdrawalCompletedEvent"] = s.handleWithdrawalCompleted

	s.handlers["CardIssuedEvent"] = s.handleCardIssued
	s.handlers["CardStatusChangedEvent"] = s.handleCardStatusChanged
}

// HandleMessage - calls the required handler
func (s *NotificationService) HandleMessage(eventType string, payload []byte) error {
	handler, exists := s.handlers[eventType]
	if !exists {
		log.Printf("[WARNING] Event of unknown type was received: %s", eventType)
		return nil
	}

	return handler(payload)
}
