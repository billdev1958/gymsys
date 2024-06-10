package user

import (
	"context"
	"gymSystem/internal/domain/user/models"
)

type Usecase interface {
	RegisterUser(ctx context.Context, request models.RegisterUserRequest) (response models.RegisterUserResponse, err error)
}
