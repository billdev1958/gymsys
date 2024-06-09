package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	postgres "gymSystem/internal/infrastructure/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {

	dbURL := "postgres://ax01:1a2s3d4f@localhost:5432/gym_sys?sslmode=disable"
	dbPool, err := pgxpool.New(context.Background(), (dbURL))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create conection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	storage := postgres.NewPgxStorage(dbPool)

	err = storage.SeedSubscriptionCost(context.Background())
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hola Mundo")
	})

	fmt.Println("Servidor escuchando en el puerto 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error iniciando el servidor:", err)
	}

}
