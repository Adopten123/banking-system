package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Adopten123/banking-system/service-notification/internal/domain"
)

// Helper для красивого форматирования источника (card/account)
func formatEntity(entityType string, entityID string) string {
	return fmt.Sprintf("%s [%s]", strings.ToUpper(entityType), entityID)
}

func (s *NotificationService) handleTransferCreated(payload []byte) error {
	var event domain.TransferCreatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal TransferCreatedEvent: %w", err)
	}

	sourceStr := formatEntity(event.SourceType, event.SourceID.String())
	destStr := formatEntity(event.DestinationType, event.DestinationID.String())

	log.Printf("[PUSH] Transfer %s: %s -> %s. Amount: %s %s\n",
		event.TransactionID, sourceStr, destStr, event.SenderAmount, event.SenderCurrency)
	return nil
}

func (s *NotificationService) handleDepositCompleted(payload []byte) error {
	var event domain.DepositCompletedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal DepositCompletedEvent: %w", err)
	}

	destStr := formatEntity(event.DestinationType, event.DestinationID.String())
	log.Printf("[PUSH] Replenished! %s received %s %s\n", destStr, event.Amount, event.Currency)
	return nil
}

func (s *NotificationService) handleWithdrawalCompleted(payload []byte) error {
	var event domain.WithdrawalCompletedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal WithdrawalCompletedEvent: %w", err)
	}

	sourceStr := formatEntity(event.SourceType, event.SourceID.String())
	log.Printf("[PUSH] Debited! %s spent %s %s\n", sourceStr, event.Amount, event.Currency)
	return nil
}
