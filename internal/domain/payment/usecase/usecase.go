package usecase

import (
	"context"
	"gymSystem/internal/domain/payment"
	"gymSystem/internal/domain/payment/entities"
	"gymSystem/internal/domain/payment/models"
	"log"

	"github.com/google/uuid"
)

type usecase struct {
	repo payment.Repository
}

func NewUsecase(repo payment.Repository) payment.Usecase {
	return &usecase{repo: repo}
}

func (u *usecase) RegisterPayment(ctx context.Context, request models.RegisterUserRequest) (response models.RegisterUserResponse, err error) {

	registerUserTx := entities.RegisterUserTx{
		User: entities.User{
			Name:      request.Name,
			Lastname1: request.Lastname1,
			Lastname2: request.Lastname2,
			Email:     request.Email,
			Phone:     request.Phone,
		},
		AccountID:          uuid.New(),
		AccountTypeID:      request.AccountTypeID,
		SubscriptionCostID: request.SubscriptionCostID,
		PaymentTypeID:      request.PaymentTypeID,
		Amount:             request.Amount,
	}
	userID, err := u.repo.RegisterUser(ctx, &registerUserTx)
	if err != nil {
		log.Printf("eror registering user: %v", err)
		return response, err
	}

	response = models.RegisterUserResponse{
		UserID:    userID,
		Name:      request.Name,
		Lastname1: request.Lastname1,
		Lastname2: request.Lastname2,
		Email:     request.Email,
		Amount:    request.Amount,
	}

	return response, nil
}
