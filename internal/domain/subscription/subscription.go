package subscription

import "time"

type SubscriptionCosts struct {
	ID               int32
	SubscriptionType int32
	SubscriptionDays int32
	Cost             float32
}

type Subscription struct {
	ID                 int32
	AccountID          int32
	SubscriptionCostID int32
	StartDate          time.Time
	EndDate            time.Time
	Status             int
}
