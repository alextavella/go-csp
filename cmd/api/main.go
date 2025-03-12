package main

import (
	"log"

	"github.com/alextavella/bank-api/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type Transaction struct {
	Amount int `json:"amount"`
}

func main() {
	repo, err := repository.NewRepository()
	if err != nil {
		log.Fatalf("Erro ao iniciar repositório: %v", err)
	}

	app := fiber.New()

	app.Post("/deposit", func(c *fiber.Ctx) error {
		var tx Transaction
		if err := c.BodyParser(&tx); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
		}

		if err := repo.Deposit(tx.Amount); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Deposit successful"})
	})

	app.Post("/withdraw", func(c *fiber.Ctx) error {
		var tx Transaction
		if err := c.BodyParser(&tx); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
		}

		if err := repo.Withdraw(tx.Amount); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"status": "success", "message": "Withdraw successful"})
	})

	app.Get("/balance", func(c *fiber.Ctx) error {
		balance, err := repo.GetBalance()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"balance": balance})
	})

	log.Fatal(app.Listen(":3000"))
}
