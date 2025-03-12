package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository() (*AccountRepository, error) {
	dsn := "root:root@tcp(db:3306)/bank"
	db, err := sql.Open("mysql", dsn)
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
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = 1").Scan(&balance)
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
