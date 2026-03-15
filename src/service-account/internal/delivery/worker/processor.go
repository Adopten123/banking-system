package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Adopten123/banking-system/service-account/internal/domain"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

func (w *RecurringWorker) processDuePayments(ctx context.Context) {
	payments, err := w.repo.GetDueRecurringPayments(ctx, 50)
	if err != nil {
		log.Printf("[WORKER ERROR] Failed to fetch due payments: %v", err)
		return
	}

	if len(payments) > 0 {
		log.Printf("[WORKER] Found %d payments to process", len(payments))
	}

	for _, p := range payments {
		w.executePayment(ctx, p)
	}
}

func (w *RecurringWorker) executePayment(ctx context.Context, p domain.RecurringPayment) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(p.CronExpression)
	if err != nil {
		log.Printf("[WORKER ERROR] Invalid cron expression for payment %s: %v", p.ID, err)
		return
	}
	nextRun := schedule.Next(time.Now().UTC())

	idempKey := fmt.Sprintf("rec-%s-%d", p.ID, p.NextExecutionTime.Unix())
	amountStr := p.Amount.String()

	var processErr error

	if p.DestinationTypeID != 0 && p.DestinationID != uuid.Nil {
		input := domain.TransferInput{
			SourceType:      resolveEntityType(p.SourceTypeID),
			SourceID:        p.SourceID.String(),
			DestinationType: resolveEntityType(p.DestinationTypeID),
			DestinationID:   p.DestinationID.String(),
			Amount:          amountStr,
			Currency:        p.CurrencyCode,
			IdempotencyKey:  idempKey,
			Description:     p.Description,
		}

		_, processErr = w.service.Transfer(ctx, input)

	} else {
		input := domain.ServiceWithdrawInput{
			SourceType:     resolveEntityType(p.SourceTypeID),
			SourceValue:    p.SourceID.String(),
			AmountStr:      amountStr,
			IdempotencyKey: idempKey,
		}

		_, processErr = w.service.Withdraw(ctx, input)
	}

	if processErr != nil {
		log.Printf("[WORKER WARNING] Payment %s failed: %v", p.ID, processErr)
	} else {
		log.Printf("[WORKER SUCCESS] Payment %s executed successfully", p.ID)
	}

	err = w.repo.UpdateRecurringPaymentNextRun(ctx, p.ID, nextRun)

	if err != nil {
		log.Printf("[WORKER CRITICAL] Failed to update next run time for %s: %v", p.ID, err)
	}
}
