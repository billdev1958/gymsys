package main

import (
	"log"

	"gymSystem/internal/app"
)

func main() {
	port := "8080" // El puerto en el que quieres que corra tu servidor HTTP

	application, err := app.NewApp(port)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}
	defer application.DB.Close()

	if err := application.Run(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
