package models

import "github.com/google/uuid"

type RegisterUserRequest struct {
	Name               string
	Lastname1          string
	Lastname2          string
	Email              string
	Phone              string
	AccountTypeID      int32
	AccountID          uuid.UUID
	SubscriptionCostID int32
	PaymentTypeID      int32
	Amount             float64
}

type RegisterUserResponse struct {
	UserID int32
}
