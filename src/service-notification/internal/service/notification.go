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
	s.handlers["TransferCreatedEvent"] = s.handleTransferCreated
	s.handlers["AccountCreatedEvent"] = s.handleAccountCreated
	s.handlers["AccountStatusChangedEvent"] = s.handleAccountStatusChanged
	s.handlers["DepositCompletedEvent"] = s.handleDepositCompleted
}

// HandleMessage - calls the required handler
func (s *NotificationService) HandleMessage(eventType string, payload []byte) error {
	handler, exists := s.handlers[eventType]
	if !exists {
		log.Printf("[WARNING] Получено событие неизвестного типа: %s", eventType)
		return nil
	}

	return handler(payload)
}
