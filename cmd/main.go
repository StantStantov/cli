package main

import (
	"context"
	"fmt"
	"lesta-battleship/cli/internal/app"
	"log"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app, err := app.New(ctx)
	if err != nil {
		log.Fatalf("Ошибка инициализации: %v", err)
	}

	fmt.Println("CLI клиент запущен.")
	if err := app.Run(); err != nil {
		log.Printf("Ошибка выполнения: %v", err)
		os.Exit(1)
	}
}
