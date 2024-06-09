package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        int
	Name      string
	Lastname1 string
	Lastname2 string
	Email     string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccountType struct {
	ID   int
	Name string
}

type Accounts struct {
	ID             int
	UserID         int
	AccountID      uuid.UUID
	AccountTypeID  int
	SubscriptionID int
	Status         int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
