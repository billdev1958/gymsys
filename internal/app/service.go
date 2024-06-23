package app

import (
	"context"
	v1 "gymSystem/internal/domain/user/http"
	"gymSystem/internal/domain/user/repository"
	"gymSystem/internal/domain/user/usecase"
	postgres "gymSystem/internal/infrastructure/db"

	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func StartUserService(ctx context.Context, db *pgxpool.Pool, router *http.ServeMux) error {
	storage := postgres.NewPgxStorage(db)

	repo := repository.NewUserRepository(storage)

	uc := usecase.NewUsecase(repo)

	h := v1.NewHandler(uc)

	h.UserRoutes(router)

	return nil
}
