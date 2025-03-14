package handler

import (
	"github.com/alextavella/bank-api/internal/domain"
	"github.com/alextavella/bank-api/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type BankHandler struct {
	repo *repository.AccountRepository
}

func NewHandler(repo *repository.AccountRepository) *BankHandler {
	return &BankHandler{
		repo: repo,
	}
}

func (h *BankHandler) Deposit(c *fiber.Ctx) error {
	var tx domain.Transaction
	if err := c.BodyParser(&tx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := h.repo.Deposit(c.Context(), tx.Amount); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Deposit successful"})
}

func (h *BankHandler) Withdraw(c *fiber.Ctx) error {
	var tx domain.Transaction
	if err := c.BodyParser(&tx); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	if err := h.repo.Withdraw(c.Context(), tx.Amount); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Withdraw successful"})
}

func (h *BankHandler) Balance(c *fiber.Ctx) error {
	balance, err := h.repo.GetBalance(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"balance": balance})
}
