package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Transaction struct {
	Amount int `json:"amount"`
}

var (
	balance int
	mu      sync.Mutex
)

func main() {
	app := fiber.New()

	// Canal para transações financeiras
	depositCh := make(chan int)
	withdrawCh := make(chan int)

	// Goroutine para processamento de transações
	go func() {
		for {
			select {
			case amount := <-depositCh:
				mu.Lock()
				balance += amount
				mu.Unlock()
				fmt.Println("Depósito:", amount, "Novo saldo:", balance)
			case amount := <-withdrawCh:
				mu.Lock()
				if balance >= amount {
					balance -= amount
					fmt.Println("Saque:", amount, "Novo saldo:", balance)
				} else {
					fmt.Println("Saque falhou: saldo insuficiente")
				}
				mu.Unlock()
			}
		}
	}()

	app.Post("/deposit", func(c *fiber.Ctx) error {
		var tx Transaction
		if err := c.BodyParser(&tx); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
		}

		// Envia a transação para o canal de depósito
		depositCh <- tx.Amount

		return c.JSON(fiber.Map{"status": "success", "message": "Deposit successful"})
	})

	app.Post("/withdraw", func(c *fiber.Ctx) error {
		var tx Transaction
		if err := c.BodyParser(&tx); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
		}

		// Envia a transação para o canal de saque
		withdrawCh <- tx.Amount

		return c.JSON(fiber.Map{"status": "success", "message": "Withdraw attempted"})
	})

	app.Get("/balance", func(c *fiber.Ctx) error {
		mu.Lock()
		currentBalance := balance
		mu.Unlock()

		return c.JSON(fiber.Map{"balance": currentBalance})
	})

	log.Fatal(app.Listen(":3000"))
}
