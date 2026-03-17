package domain

import (
	"time"

	"github.com/google/uuid"
)

type Report struct {
	ID         int64
	ReporterID uuid.UUID
	TargetType string
	TargetID   int64
	Reason     string
	Status     string
	CreatedAt  time.Time
}

type CreateReportInput struct {
	ReporterID uuid.UUID
	TargetType string
	TargetID   int64
	Reason     string
}