package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        int32
	Name      string
	Lastname1 string
	Lastname2 string
	Email     string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccountType struct {
	ID   int32
	Name string
}

type Accounts struct {
	ID             int32
	UserID         int32
	AccountID      uuid.UUID
	AccountTypeID  int32
	SubscriptionID int32
	Status         int32
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type RegisterUsertx struct {
	User
	AccountID          uuid.UUID
	AccountTypeID      int32
	SubscriptionCostID int32
	SubscriptionID     int32
	PaymentTypeID      int32
	Amount             float64
}
