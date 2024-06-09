package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxStorage struct {
	DBPool *pgxpool.Pool
}

func NewPgxStorage(dbpool *pgxpool.Pool) *PgxStorage {
	return &PgxStorage{DBPool: dbpool}
}

// Account Types constants and structure

type account struct {
	Name string
}

type subscriptionCost struct {
	Type string
	Days int32
	Cost float64
}

type data struct {
	SubscriptionCosts []subscriptionCost
}

func (storage *PgxStorage) SeedSubscriptionCost(ctx context.Context) (err error) {
	d := data{
		SubscriptionCosts: []subscriptionCost{
			{Type: "Sencilla", Days: 1, Cost: 29},
			{Type: "Sencilla", Days: 30, Cost: 299},
			{Type: "Sencilla", Days: 90, Cost: 879},
			{Type: "Normal", Days: 180, Cost: 1619},
			{Type: "Premium", Days: 365, Cost: 3109},
		},
	}

	var count int
	err = storage.DBPool.QueryRow(ctx, "SELECT COUNT(*) FROM subscription_costs").Scan(&count)
	if count > 0 {
		fmt.Println("la tabla subscription_costs ya contiene datos.")
		return nil
	}

	query := "INSERT INTO subscription_costs (subscription_type, subscription_day, cost) VALUES($1, $2, $3)"

	for i := range d.SubscriptionCosts {
		sub := &d.SubscriptionCosts[i]
		_, err := storage.DBPool.Exec(ctx, query, sub.Type, sub.Days, sub.Cost)
		if err != nil {
			return fmt.Errorf("insert subscription cost: %w", err)
		}
	}

	fmt.Println("Valores insertados correctamente en suscription_costs")
	return nil
}

type status struct {
	Name string
}

func (storage *PgxStorage) SeedStatus(ctx context.Context) (err error) {
	statusValues := [5]string{"activa", "inactiva", "cancelada", "en proceso", "exitosa"}

	var count int
	err = storage.DBPool.QueryRow(ctx, "SELECT COUNT(*) FROM status").Scan(&count)
	if count > 0 {
		fmt.Println("la tabla status ya contiene datos")
		return nil
	}

	query := "INSERT INTO status (name) VALUES($1)"
	for _, value := range statusValues {
		_, err := storage.DBPool.Exec(ctx, query, value)
		if err != nil {
			return fmt.Errorf("isert status: %w", err)
		}
	}

	fmt.Println("Valores insertados correctamente en status")
	return nil
}

func (storage *PgxStorage) SeedAccountValues(ctx context.Context) (err error) {
	accountValues := [3]string{"admin", "empleado", "cliente"}

	var count int
	err = storage.DBPool.QueryRow(ctx, "SELECT COUNT(*) FROM account_types").Scan(&count)

	if count > 0 {
		fmt.Println("la tabla account_types ya contiene datos.")
		return nil
	}

	query := "INSERT INTO account_types (name) VALUES ($1)"
	for _, value := range accountValues {
		_, err := storage.DBPool.Exec(ctx, query, value)
		if err != nil {
			return fmt.Errorf("insert accountType value: %w", err)
		}
	}

	fmt.Println("Valores insertados correctamente en account_types")
	return nil
}

func (storage *PgxStorage) SeedPaymentTypes(ctx context.Context) (err error) {
	paymentTypes := [3]string{"efectivo", "tarjeta", "transferencia"}

	var count int
	err = storage.DBPool.QueryRow(ctx, "SELECT COUNT(*) FROM payment_types").Scan(&count)

	if count > 0 {
		fmt.Println("la tabla payment_types ya contiene datos")
		return nil
	}

	query := "INSERT INTO payment_types (name) VALUES ($1)"
	for _, value := range paymentTypes {
		_, err := storage.DBPool.Exec(ctx, query, value)
		if err != nil {
			return fmt.Errorf("insert payment_type value: %w\n", err)
		}
	}

	fmt.Println("Valores insertados correctamente en payment_types")
	return nil
}
