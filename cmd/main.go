package main

import (
	"context"
	"log"

	"github.com/Dnlbb/chat-server/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatal("failed to create app")
	}

	if err = a.Run(); err != nil {
		log.Fatal("failed to run app")
	}
}
