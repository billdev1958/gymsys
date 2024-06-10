package models

type RegisterUserRequest struct {
	Name               string
	Lastname1          string
	Lastname2          string
	Email              string
	Phone              string
	AccountTypeID      int32
	SubscriptionCostID int32
	PaymentTypeID      int32
	Ammount            float64
}

type RegisterUserResponse struct {
	UserID int32
}
