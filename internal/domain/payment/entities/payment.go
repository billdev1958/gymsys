package entities

import (
	"time"

	"github.com/google/uuid"
)

type PaymentType struct {
	ID   int
	Name string
}

type Payment struct {
	ID            int
	AccountID     int
	PaymentTypeID int
	Cost          float32
	Status        string
	PaymentDate   time.Time
}

type RegisterUserTx struct {
	User
	AccountID          uuid.UUID
	AccountTypeID      int32
	SubscriptionCostID int32
	PaymentTypeID      int32
	Amount             float64
}

type User struct {
	ID        int32
	Name      string
	Lastname1 string
	Lastname2 string
	Email     string
	Phone     string
}
