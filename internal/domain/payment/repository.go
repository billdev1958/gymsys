package payment

import (
	"context"
	"gymSystem/internal/domain/payment/entities"
)

type Repository interface {
	RegisterUser(ctx context.Context, register *entities.RegisterUserTx) (int32, error)
}
