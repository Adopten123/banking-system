package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Adopten123/banking-system/service-notification/internal/domain"
)

func (s *NotificationService) handleTransferCreated(payload []byte) error {
	var event domain.TransferCreatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal TransferCreatedEvent: %w", err)
	}

	log.Printf("[PUSH] Перевод %s: со счета %d на счет %d, сумма: %s %s\n",
		event.TransactionID, event.FromAccountID, event.ToAccountID, event.Amount, event.Currency)
	return nil
}

func (s *NotificationService) handleAccountCreated(payload []byte) error {
	var event domain.AccountCreatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal AccountCreatedEvent: %w", err)
	}

	log.Printf("[EMAIL] Добро пожаловать! Создан счет %d (%s)\n", event.AccountID, event.Currency)
	return nil
}

func (s *NotificationService) handleAccountStatusChanged(payload []byte) error {
	var event domain.AccountStatusChangedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal AccountStatusChangedEvent: %w", err)
	}

	log.Printf("[SMS] Статус вашего счета %d изменен с %d на %d\n", event.AccountID, event.OldStatus, event.NewStatus)
	return nil
}

func (s *NotificationService) handleDepositCompleted(payload []byte) error {
	var event domain.DepositCompletedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal DepositCompletedEvent: %w", err)
	}

	log.Printf("[PUSH] Счет %d успешно пополнен на %s %s\n", event.AccountID, event.Amount, event.Currency)
	return nil
}
