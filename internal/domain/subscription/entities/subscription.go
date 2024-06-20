package entities

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionCosts struct {
	ID               int32
	SubscriptionType int32
	SubscriptionDays int32
	Cost             float32
}

type Subscription struct {
	ID                 int32
	AccountID          uuid.UUID
	SubscriptionCostID int32
	StartDate          time.Time
	EndDate            time.Time
	Status             int
}
