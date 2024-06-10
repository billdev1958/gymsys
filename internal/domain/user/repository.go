package user

import (
	"context"
	"gymSystem/internal/domain/user/entities"
)

type Repository interface {
	RegisterUser(ctx context.Context, register *entities.RegisterUsertx) (int32, error)
}
