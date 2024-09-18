package main

import (
	"log"

	"github.com/cvzm/go-web-project/bootstrap"
)

func main() {
	// Initialize App using Wire
	app, err := bootstrap.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// Set up and run the application
	log.Printf("Starting server")
	if err := app.SetupAndRun(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}
}
