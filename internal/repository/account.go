package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/alextavella/bank-api/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository() (*AccountRepository, error) {
	db, err := sql.Open("mysql", config.DB_URI)
	db.SetMaxOpenConns(50)                 // Máximo de conexões abertas
	db.SetMaxIdleConns(5)                  // Máximo de conexões inativas
	db.SetConnMaxLifetime(time.Minute * 5) // Tempo máximo de uma conexão

	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao verificar conexão: %w", err)
	}

	return &AccountRepository{db: db}, nil
}

func (r *AccountRepository) Deposit(amount int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var balance int
	// Lock explícito para garantir consistência
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = 1 FOR UPDATE").Scan(&balance)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = 1", amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *AccountRepository) Withdraw(amount int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var balance int
	// Lock explícito para garantir consistência
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = 1 FOR UPDATE").Scan(&balance)
	if err != nil {
		return err
	}

	if balance < amount {
		return fmt.Errorf("saldo insuficiente")
	}

	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = 1", amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *AccountRepository) GetBalance() (int, error) {
	var balance int
	err := r.db.QueryRow("SELECT balance FROM accounts WHERE id = 1").Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
