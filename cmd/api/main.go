package main

import (
	"log"

	"github.com/alextavella/bank-api/internal/handler"
	"github.com/alextavella/bank-api/internal/repository"
	"github.com/gofiber/fiber/v2"
)

func main() {
	repo, err := repository.NewAccountRepository()
	if err != nil {
		log.Fatalf("Erro ao iniciar repositório: %v", err)
	}

	bankHandler := handler.NewHandler(repo)

	app := fiber.New()
	app.Post("/deposit", bankHandler.Deposit)
	app.Post("/withdraw", bankHandler.Withdraw)
	app.Get("/balance", bankHandler.Balance)

	log.Fatal(app.Listen(":3000"))
}
