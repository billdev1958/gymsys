package repository_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"gymSystem/internal/domain/user/entities"
	user "gymSystem/internal/domain/user/repository"
	postgres "gymSystem/internal/infrastructure/db"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	dbURL := "postgres://ax01:1a2s3d4f@localhost:5432/gym_sys?sslmode=disable"

	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("failed to connect to the database: %v", err)
	}

	_, err = dbPool.Exec(context.Background(), "TRUNCATE TABLE users, accounts, subscriptions, payments RESTART IDENTITY CASCADE;")
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}

	storage := postgres.NewPgxStorage(dbPool)

	// Sembrar datos necesarios
	err = storage.SeedStatus(context.Background())
	if err != nil {
		t.Fatalf("failed to seed status: %v", err)
	}

	err = storage.SeedAccountValues(context.Background())
	if err != nil {
		t.Fatalf("failed to seed account values: %v", err)
	}

	err = storage.SeedPaymentTypes(context.Background())
	if err != nil {
		t.Fatalf("failed to seed payment types: %v", err)
	}

	err = storage.SeedSubscriptionCost(context.Background())
	if err != nil {
		t.Fatalf("failed to seed payment types: %v", err)
	}

	return dbPool
}

func TestRegisterUser(t *testing.T) {
	dbPool := setupTestDB(t)
	defer dbPool.Close() // Asegura que la conexión se cierre después de las pruebas

	userRepo := user.NewUserRepository(&postgres.PgxStorage{DBPool: dbPool})

	t.Run("Success", func(t *testing.T) {
		register := &entities.RegisterUsertx{
			User: entities.User{
				Name:      "John",
				Lastname1: "Doe",
				Lastname2: "Smith",
				Email:     "john.doe@example.com",
				Phone:     "1234567890",
				CreatedAt: time.Now(),
			},
			AccountTypeID:      1,
			SubscriptionCostID: 1, // Correspondiente a 'Sencilla', 1 día, 29.00
			PaymentTypeID:      1,
			Amount:             29.0,
			AccountID:          uuid.New(), // Generar un nuevo UUID para la cuenta
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Añadir un contexto con timeout
		defer cancel()

		userID, err := userRepo.RegisterUser(ctx, register)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if userID == 0 {
			t.Fatalf("expected valid userID, got %v", userID)
		}

		// Verificar que el usuario se ha registrado correctamente
		var count int
		err = dbPool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE id = $1", userID).Scan(&count)
		if err != nil || count != 1 {
			t.Fatalf("expected user to be registered, got count %v, err %v", count, err)
		}
	})

	t.Run("Fail_InsertUser", func(t *testing.T) {
		register := &entities.RegisterUsertx{
			User: entities.User{
				Name:      "John",
				Lastname1: "Doe",
				Lastname2: "Smith",
				Email:     "", // Invalid email to trigger an error
				Phone:     "1234567890",
				CreatedAt: time.Now(),
			},
			AccountTypeID:      1,
			SubscriptionCostID: 1, // Correspondiente a 'Sencilla', 1 día, 29.00
			PaymentTypeID:      1,
			Amount:             29.0,
			AccountID:          uuid.New(), // Generar un nuevo UUID para la cuenta
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Añadir un contexto con timeout
		defer cancel()

		_, err := userRepo.RegisterUser(ctx, register)
		if err == nil || !strings.Contains(err.Error(), "insert user") {
			t.Fatalf("expected insert user error, got %v", err)
		}
	})

	t.Run("Fail_CostMismatch", func(t *testing.T) {
		register := &entities.RegisterUsertx{
			User: entities.User{
				Name:      "John",
				Lastname1: "Doe",
				Lastname2: "Smith",
				Email:     "john.doe2@example.com", // Using a different email to avoid unique constraint error
				Phone:     "1234567890",
				CreatedAt: time.Now(),
			},
			AccountTypeID:      1,
			SubscriptionCostID: 1, // Correspondiente a 'Sencilla', 1 día, 29.00
			PaymentTypeID:      1,
			Amount:             9999.99,    // Mismatch cost to trigger an error
			AccountID:          uuid.New(), // Generar un nuevo UUID para la cuenta
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Añadir un contexto con timeout
		defer cancel()

		_, err := userRepo.RegisterUser(ctx, register)
		if err == nil || !strings.Contains(err.Error(), "amount incorrect") {
			t.Fatalf("expected amount incorrect error, got %v", err)
		}
	})
}
