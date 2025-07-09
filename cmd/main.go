package main

import (
	"fmt"
	"lesta-battleship/cli/internal/app"
	"log"
	"os"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatalf("Ошибка инициализации: %v", err)
	}

	fmt.Println("CLI клиент запущен.")
	if err := app.Run(); err != nil {
		log.Printf("Ошибка выполнения: %v", err)
		os.Exit(1)
	}
}
