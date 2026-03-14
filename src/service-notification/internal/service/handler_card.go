package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Adopten123/banking-system/service-notification/internal/domain"
)

func (s *NotificationService) handleCardIssued(payload []byte) error {
	var event domain.CardIssuedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal CardIssuedEvent: %w", err)
	}

	log.Printf("[SMS] Good news! Your new card %s has been successfully issued\n", event.PanMask)
	return nil
}

func (s *NotificationService) handleCardStatusChanged(payload []byte) error {
	var event domain.CardStatusChangedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal CardStatusChangedEvent: %w", err)
	}

	log.Printf("[PUSH] Attention! Card %s status changed from '%s' to '%s'\n",
		event.CardID, event.OldStatus, event.NewStatus)
	return nil
}