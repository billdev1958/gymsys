package payment

import (
	"context"
	"gymSystem/internal/domain/payment/models"
)

type Usecase interface {
	RegisterPayment(ctx context.Context, request models.RegisterUserRequest) (response models.RegisterUserResponse, err error)
}
