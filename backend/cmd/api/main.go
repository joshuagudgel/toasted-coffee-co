package main

import (
	"log"

	"github.com/joshuagudgel/toasted-coffee/backend/internal/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer app.Close()

	if err := app.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
