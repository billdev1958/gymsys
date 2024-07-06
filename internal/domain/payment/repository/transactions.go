package repository

import (
	"context"
	"fmt"
	"gymSystem/internal/domain/payment"
	"gymSystem/internal/domain/payment/entities"
	postgres "gymSystem/internal/infrastructure/db"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type paymentRepository struct {
	storage *postgres.PgxStorage
}

func NewPaymentRepository(storage *postgres.PgxStorage) payment.Repository {
	return &paymentRepository{storage: storage}
}

func (pr *paymentRepository) RegisterUser(ctx context.Context, register *entities.RegisterUserTx) (userID int32, err error) {

	ctxTx, cancel := context.WithTimeout(ctx, 5*time.Second) // Añadir un contexto con timeout de 5 segundos	defer cancel()
	defer cancel()

	tx, err := pr.storage.DBPool.Begin(ctxTx)
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
				log.Printf("error commiting transaction: %v", err)
			}

		}
	}()

	var expectedCost float64
	var subscriptionDays int
	err = tx.QueryRow(ctxTx, "SELECT cost, subscription_day from subscription_costs where id = $1", register.SubscriptionCostID).Scan(&expectedCost, &subscriptionDays)
	if err != nil {
		log.Printf("error getting subscription cost and days: %v", err)
		cancel()
		return 0, fmt.Errorf("get subscription cost and days: %v", err)
	}

	if register.Amount != expectedCost {
		log.Printf("invalid cost: expected %v, got %v", expectedCost, register.Amount)
		cancel()
		return 0, fmt.Errorf("amount incorrect: expected %v, got %v", expectedCost, register.Amount)
	}

	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, subscriptionDays)

	err = tx.QueryRow(ctxTx, "INSERT INTO users (name, lastname1, lastname2, email, phone, created_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		register.Name, register.Lastname1, register.Lastname2, register.Email, register.Phone, time.Now()).Scan(&register.ID)
	if err != nil {
		log.Printf("error inserting user: %v", err)
		cancel()
		return 0, fmt.Errorf("insert user: %w", err)
	}

	batch := &pgx.Batch{}

	// Crear cuenta
	batch.Queue("INSERT INTO accounts (user_id, account_id, account_type_id, created_at) VALUES($1, $2, $3, $4)",
		register.ID, register.AccountID, register.AccountTypeID, time.Now())

	// Registrar subscripcion
	batch.Queue("INSERT INTO subscriptions (account_id, subscription_cost_id, start_date, end_date) VALUES($1, $2, $3, $4)",
		register.AccountID, register.SubscriptionCostID, startDate, endDate)

	// Registrar pago
	batch.Queue("INSERT INTO payments (account_id, payment_type_id, cost, payment_date) VALUES($1, $2, $3, $4)",
		register.AccountID, register.PaymentTypeID, register.Amount, time.Now())

	br := tx.SendBatch(ctxTx, batch)
	defer br.Close()

	for i := 0; i < 2; i++ {
		_, err := br.Exec()
		if err != nil {
			log.Printf("error executing batch: %v", err)
			cancel()
			return 0, fmt.Errorf("execute batch: %w", err)
		}
	}

	return register.ID, nil
}
