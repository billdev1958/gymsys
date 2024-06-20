package subscription

import (
	"context"
	"gymSystem/internal/domain/subscription/entities"

	"github.com/jackc/pgx/v5"
)

type Repository interface {
	RegisterSubscription(ctx context.Context, tx pgx.Tx, register *entities.Subscription) (int32, error)
}
