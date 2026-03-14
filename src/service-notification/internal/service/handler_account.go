package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Adopten123/banking-system/service-notification/internal/domain"
)

func (s *NotificationService) handleAccountCreated(payload []byte) error {
	var event domain.AccountCreatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal AccountCreatedEvent: %w", err)
	}

	log.Printf("[EMAIL] Welcome! Account %d (%s) created successfully\n", event.AccountID, event.Currency)
	return nil
}

func (s *NotificationService) handleAccountStatusChanged(payload []byte) error {
	var event domain.AccountStatusChangedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal AccountStatusChangedEvent: %w", err)
	}

	log.Printf("[SMS] Your account %d status has changed from %d to %d\n", event.AccountID, event.OldStatus, event.NewStatus)
	return nil
}

func (s *NotificationService) handleCreditLimitChanged(payload []byte) error {
	var event domain.CreditLimitChangedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal CreditLimitChangedEvent: %w", err)
	}

	log.Printf("[EMAIL] Credit limit for account %s updated: old limit %s, new limit %s %s\n",
		event.AccountID, event.OldLimit, event.NewLimit, event.Currency)
	return nil
}
