package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "gymSystem/internal/domain/user/http"
	"gymSystem/internal/domain/user/repository"
	"gymSystem/internal/domain/user/usecase"
	postgres "gymSystem/internal/infrastructure/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	dbPool := setupDatabase()
	defer dbPool.Close()

	seedDatabase(ctx, dbPool)

	storage := postgres.NewPgxStorage(dbPool)
	repo := repository.NewUserRepository(storage)
	uc := usecase.NewUsecase(repo)
	h := v1.NewHandler(uc)

	mux := http.NewServeMux()
	h.UserRoutes(mux)

	go func() {
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	fmt.Println("Server is running on port 8080")

	// Graceful shutdown
	gracefulShutdown(ctx, dbPool)
}

func setupDatabase() *pgxpool.Pool {
	dbURL := "postgres://ax01:1a2s3d4f@localhost:5432/gym_sys?sslmode=disable"
	dbPool, err := pgxpool.New(context.Background(), (dbURL))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create conection pool: %v\n", err)
		os.Exit(1)
	}
	return dbPool
}

func seedDatabase(ctx context.Context, dbPool *pgxpool.Pool) {
	storage := postgres.NewPgxStorage(dbPool)

	err := storage.SeedSubscriptionCost(context.Background())
	if err != nil {
		log.Fatalf("error seeding subscription costs %v\n", err)
	}

	err = storage.SeedStatus(context.Background())
	if err != nil {
		log.Fatalf("error seeding status %v\n", err)
	}

	err = storage.SeedAccountValues(context.Background())
	if err != nil {
		log.Fatalf("error seeding account values %v\n", err)
	}

	err = storage.SeedPaymentTypes(context.Background())
	if err != nil {
		log.Fatalf("error seeding paymentTypes values")
	}

}

// gracefulShutdown handles the graceful shutdown of the server and database connection
func gracefulShutdown(ctx context.Context, dbPool *pgxpool.Pool) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM)
	signal.Notify(quit, syscall.SIGINT)

	<-quit
	fmt.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	dbPool.Close()
	fmt.Println("Database connection closed")
}
