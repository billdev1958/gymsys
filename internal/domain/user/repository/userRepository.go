package user

import (
	"context"
	"fmt"
	"gymSystem/internal/domain/user/models"
	postgres "gymSystem/internal/infrastructure/db"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, register *RegisterUsertx) (int, error)
}

type userRepository struct {
	storage *postgres.PgxStorage
}

func NewUserRepository(storage *postgres.PgxStorage) UserRepository {
	return &userRepository{storage: storage}
}

type RegisterUsertx struct {
	models.User
	AccountTypeID      int32
	SubscriptionCostID int32
	PaymentTypeID      int32
	Ammount            float64
}

func (ur *userRepository) RegisterUser(ctx context.Context, register *RegisterUsertx) (int, error) {

	tx, err := ur.storage.DBPool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	var userID int
	query := "INSERT INTO users (name, lastname1, lastname2, email, phone, created_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id"
	err = tx.QueryRow(ctx, query, register.Name, register.Lastname1, register.Lastname2, register.Email, register.Phone, register.CreatedAt).Scan(&userID)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("insert user: %w", err)
	}
	register.ID = userID

	// Crear cuenta
	accountID := uuid.New()
	var account uuid.UUID
	query = "INSERT INTO accounts (user_id, account_id, account_type_id, created_at) VALUES($1, $2, $3, $4) RETURNING account_id"
	err = tx.QueryRow(ctx, query, userID, accountID, register.AccountTypeID, time.Now()).Scan(&account)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("insert account: %w", err)
	}

	var subscriptionDuration int
	query = "SELECT subscription_day FROM subscription_costs WHERE id = $1"
	err = tx.QueryRow(ctx, query, register.SubscriptionCostID).Scan(&subscriptionDuration)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("get subscription duration: %w", err)
	}

	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, subscriptionDuration)

	// Crear subscripcion
	query = "INSERT INTO subscriptions (account_id, subscription_cost_id, start_date, end_date) VALUES($1, $2, $3, $4)"
	_, err = tx.Exec(ctx, query, account, register.SubscriptionCostID, startDate, endDate)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("insert subscription: %w", err)
	}

	var expectedCost float64
	query = "SELECT cost from subscription_costs where id = $1"
	err = tx.QueryRow(ctx, query, register.SubscriptionCostID).Scan(&expectedCost)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("ammount: %w", err)
	}

	if register.Ammount != expectedCost {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("ammount incorrect: %w", err)
	}

	var status int32
	query = "SELECT id from status where id = 5"
	err = tx.QueryRow(ctx, query).Scan(&status)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("get status: %w", err)

	}

	query = "INSERT INTO payments (account_id, payment_type_id, cost, status_id, payment_date) VALUES($1, $2, $3, $4, $5)"
	_, err = tx.Exec(ctx, query, account, register.PaymentTypeID, register.Ammount, status, time.Now())
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("insert payment: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("commit transaction: %w", err)
	}

	return userID, nil
}
