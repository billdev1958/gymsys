package repository

import (
	"context"
	"fmt"
	"gymSystem/internal/domain/user"
	"gymSystem/internal/domain/user/entities"
	postgres "gymSystem/internal/infrastructure/db"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	storage *postgres.PgxStorage
}

func NewUserRepository(storage *postgres.PgxStorage) user.Repository {
	return &userRepository{storage: storage}
}

func (ur *userRepository) RegisterUser(ctx context.Context, register *entities.RegisterUsertx) (userID int32, err error) {

	ctxTx, cancel := context.WithCancel(ctx)
	defer cancel()

	tx, err := ur.storage.DBPool.Begin(ctxTx)
	if err != nil {
		log.Printf("error beginning transaction: %v", err)
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctxTx); rbErr != nil {
				log.Printf("error rolling back transaction: %v", rbErr)
			}
		} else {
			err = tx.Commit(ctxTx)
			if err != nil {
				log.Printf("error committing transaction: %v", err)
			}
		}
	}()

	userID, err = ur.registerUser(ctxTx, tx, register)
	if err != nil {
		return 0, err
	}

	accountID, err := ur.createAccount(ctxTx, tx, register)
	if err != nil {
		return 0, err
	}

	err = ur.createSubscription(ctxTx, tx, accountID, register)
	if err != nil {
		return 0, err
	}

	err = ur.verifyAndInsertPayment(ctxTx, tx, accountID, register)
	if err != nil {
		return 0, err
	}

	err = ur.updateStatuses(ctxTx, tx, accountID, register)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (ur *userRepository) registerUser(ctx context.Context, tx pgx.Tx, register *entities.RegisterUsertx) (int32, error) {
	query := "INSERT INTO users (name, lastname1, lastname2, email, phone, created_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id"
	err := tx.QueryRow(ctx, query, register.Name, register.Lastname1, register.Lastname2, register.Email, register.Phone, register.CreatedAt).Scan(&register.ID)
	if err != nil {
		log.Printf("error inserting user: %v", err)
		return 0, fmt.Errorf("insert user: %w", err)
	}

	return register.ID, nil
}

func (ur *userRepository) createAccount(ctx context.Context, tx pgx.Tx, register *entities.RegisterUsertx) (uuid.UUID, error) {
	query := "INSERT INTO accounts (user_id, account_id, account_type_id, created_at) VALUES($1, $2, $3, $4) RETURNING account_id"
	err := tx.QueryRow(ctx, query, register.ID, register.AccountID, register.AccountTypeID, time.Now()).Scan(&register.AccountID)
	if err != nil {
		log.Printf("error inserting account: %v", err)
		return uuid.Nil, fmt.Errorf("insert account: %w", err)
	}
	return register.AccountID, nil
}

func (ur *userRepository) createSubscription(ctx context.Context, tx pgx.Tx, accountID uuid.UUID, register *entities.RegisterUsertx) error {
	var subscriptionDuration int
	query := "SELECT subscription_day FROM subscription_costs WHERE id = $1"
	err := tx.QueryRow(ctx, query, register.SubscriptionCostID).Scan(&subscriptionDuration)
	if err != nil {
		log.Printf("error getting subscription duration: %v", err)
		return fmt.Errorf("get subscription duration: %w", err)
	}

	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, subscriptionDuration)

	// Crear subscripcion
	query = "INSERT INTO subscriptions (account_id, subscription_cost_id, start_date, end_date) VALUES($1, $2, $3, $4) RETURNING id"
	err = tx.QueryRow(ctx, query, accountID, register.SubscriptionCostID, startDate, endDate).Scan(&register.SubscriptionID)
	if err != nil {
		log.Printf("error inserting subscription: %v", err)
		return fmt.Errorf("insert subscription: %w", err)
	}
	return nil
}

func (ur *userRepository) verifyAndInsertPayment(ctx context.Context, tx pgx.Tx, accountID uuid.UUID, register *entities.RegisterUsertx) error {
	var expectedCost float64
	query := "SELECT cost from subscription_costs where id = $1"
	err := tx.QueryRow(ctx, query, register.SubscriptionCostID).Scan(&expectedCost)
	if err != nil {
		log.Printf("error getting subscription cost: %v", err)
		return fmt.Errorf("amount: %w", err)
	}

	if register.Amount != expectedCost {
		log.Printf("amount incorrect: expected %v, got %v", expectedCost, register.Amount)
		return fmt.Errorf("amount incorrect: %w", err)
	}

	query = "INSERT INTO payments (account_id, payment_type_id, cost, payment_date) VALUES($1, $2, $3, $4)"
	_, err = tx.Exec(ctx, query, accountID, register.PaymentTypeID, register.Amount, time.Now())
	if err != nil {
		log.Printf("error inserting payment: %v", err)
		return fmt.Errorf("insert payment: %w", err)
	}
	return nil
}

func (ur *userRepository) updateStatuses(ctx context.Context, tx pgx.Tx, accountID uuid.UUID, register *entities.RegisterUsertx) error {
	var statusPayment, statusAccount int32

	query := `SELECT id FROM STATUS WHERE id = 5 OR id = 1`
	rows, err := tx.Query(ctx, query)
	if err != nil {
		log.Printf("error getting statuses: %v", err)
		return fmt.Errorf("get statuses: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			log.Printf("error scanning status: %v", err)
			return fmt.Errorf("scan status: %w", err)
		}

		if id == 5 {
			statusPayment = id
		} else {
			statusAccount = id
		}
	}
	if err := rows.Err(); err != nil {
		log.Printf("rows error: %v", err)
		return fmt.Errorf("rows error: %w", err)
	}

	query = "UPDATE payments SET status_id = $1 WHERE account_id = $2"
	result, err := tx.Exec(ctx, query, statusPayment, accountID)
	if err != nil {
		log.Printf("error updating payment status: %v", err)
		return fmt.Errorf("update payment status: %w", err)
	}
	if result.RowsAffected() == 0 {
		log.Printf("no rows affected when updating payment status")
		return fmt.Errorf("no rows affected when updating payment status")
	}

	query = "UPDATE accounts SET subscription_id = $1, status_id = $2 WHERE account_id = $3"
	result, err = tx.Exec(ctx, query, register.SubscriptionID, statusAccount, accountID)
	if err != nil {
		log.Printf("error updating account status: %v", err)
		return fmt.Errorf("update account status: %w", err)
	}
	if result.RowsAffected() == 0 {
		log.Printf("no rows affected when updating account status")
		return fmt.Errorf("no rows affected when updating account status")
	}
	return nil
}
