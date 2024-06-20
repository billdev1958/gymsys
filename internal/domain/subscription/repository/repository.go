package repository

import (
	"context"
	"fmt"
	"gymSystem/internal/domain/subscription"
	"gymSystem/internal/domain/subscription/entities"
	postgres "gymSystem/internal/infrastructure/db"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

type SubscriptionRepository struct {
	Storage *postgres.PgxStorage
}

func NewSubscriptionRepository(storage *postgres.PgxStorage) subscription.Repository {
	return &SubscriptionRepository{Storage: storage}
}

func (sr SubscriptionRepository) RegisterSubscription(ctx context.Context, tx pgx.Tx, register *entities.Subscription) (int32, error) {
	var subscriptionDuration int
	query := "SELECT subscription_day FROM subscription_costs WHERE id = $1"
	err := tx.QueryRow(ctx, query, register.SubscriptionCostID).Scan(&subscriptionDuration)
	if err != nil {
		log.Printf("error getting subscription duration: %v", err)
		return 0, fmt.Errorf("get subscription duration: %w", err)
	}

	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, subscriptionDuration)

	query = "INSERT INTO subscriptions (account_id, subscription_cost_id, start_date, end_date) VALUES($1, $2, $3, $4) RETURNING id"
	err = tx.QueryRow(ctx, query, register.AccountID, register.SubscriptionCostID, startDate, endDate).Scan(&register.ID)
	if err != nil {
		log.Printf("error inserting subscription: %v", err)
		return 0, fmt.Errorf("insert subscription: %w", err)
	}
	return register.ID, nil
}
