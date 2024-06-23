package app

import (
	"context"
	"fmt"
	postgres "gymSystem/internal/infrastructure/db"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	DB     *pgxpool.Pool
	port   string
	router *http.ServeMux
}

func NewApp(port string) (*App, error) {
	dsn := "postgres://ax01:1a2s3d4f@localhost:5432/gym_sys?sslmode=disable"

	db, err := setupDatabase(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to setup database: %w", err)
	}

	if err := seedDatabase(context.Background(), db); err != nil {
		return nil, fmt.Errorf("failed to seed database: %w", err)
	}

	router := http.NewServeMux()

	return &App{
		DB:     db,
		port:   port,
		router: router,
	}, nil
}

func (app *App) Run() error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.port),
		Handler: app.router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-quit
		log.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		server.Shutdown(ctx)
		app.DB.Close()

	}()

	log.Printf("Starting server on port %s...", app.port)
	return server.ListenAndServe()
}

func setupDatabase(dsn string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	return dbpool, nil
}

func seedDatabase(ctx context.Context, dbPool *pgxpool.Pool) error {
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
	return nil
}
