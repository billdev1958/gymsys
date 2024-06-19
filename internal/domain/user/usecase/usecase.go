package usecase

import (
	"context"
	"gymSystem/internal/domain/user"
	"gymSystem/internal/domain/user/entities"
	"gymSystem/internal/domain/user/models"
	"log"

	"github.com/google/uuid"
)

type usecase struct {
	repo user.Repository
}

func NewUsecase(repo user.Repository) user.Usecase {
	return &usecase{repo: repo}
}

func (u *usecase) RegisterUser(ctx context.Context, request models.RegisterUserRequest) (response models.RegisterUserResponse, err error) {

	registerUsertx := entities.RegisterUsertx{
		User: entities.User{
			Name:      request.Name,
			Lastname1: request.Lastname1,
			Lastname2: request.Lastname2,
			Email:     request.Email,
			Phone:     request.Phone,
		},
		AccountTypeID:      request.AccountTypeID,
		AccountID:          uuid.New(),
		SubscriptionCostID: request.SubscriptionCostID,
		PaymentTypeID:      request.PaymentTypeID,
		Amount:             request.Amount,
	}

	userID, err := u.repo.RegisterUser(ctx, &registerUsertx)
	if err != nil {
		log.Printf("error registering user: %v", err)
		return response, err
	}

	response = models.RegisterUserResponse{
		UserID: userID,
	}

	return response, nil
}
