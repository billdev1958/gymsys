package payment

import "time"

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
