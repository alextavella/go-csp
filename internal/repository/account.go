package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/alextavella/bank-api/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

type AccountRepository struct {
	db    *sql.DB
	cache *redis.Client
}

const (
	CACHE_KEY = "balance"
)

func NewAccountRepository() (*AccountRepository, error) {
	db, err := sql.Open("mysql", config.DB_URI)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao verificar conexão: %w", err)
	}

	db.SetMaxOpenConns(50)                 // Máximo de conexões abertas
	db.SetMaxIdleConns(5)                  // Máximo de conexões inativas
	db.SetConnMaxLifetime(time.Minute * 5) // Tempo máximo de uma conexão

	// Conexão com o Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.CACHE_URI,
		Password: "", // Sem senha
		DB:       0,  // Banco padrão
	})

	return &AccountRepository{db: db, cache: rdb}, nil
}

func (r *AccountRepository) Deposit(ctx context.Context, amount int) error {
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

	newBalance := balance + amount
	_, err = tx.Exec("UPDATE accounts SET balance = ? WHERE id = 1", newBalance)
	if err != nil {
		return err
	}

	// Atualiza o saldo no cache
	r.cache.Set(ctx, CACHE_KEY, newBalance, time.Minute)

	return tx.Commit()
}

func (r *AccountRepository) Withdraw(ctx context.Context, amount int) error {
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

	newBalance := balance - amount
	_, err = tx.Exec("UPDATE accounts SET balance = ? WHERE id = 1", newBalance)
	if err != nil {
		return err
	}

	// Atualiza o saldo no cache
	r.cache.Set(ctx, CACHE_KEY, newBalance, time.Minute)

	return tx.Commit()
}

func (r *AccountRepository) GetBalance(ctx context.Context) (int, error) {
	var balance int

	// Tenta buscar no cache
	balance, err := r.cache.Get(ctx, "balance").Int()
	if err == nil {
		return balance, nil
	}

	// Se não estiver no cache, busca no banco
	err = r.db.QueryRow("SELECT balance FROM accounts WHERE id = 1").Scan(&balance)
	if err != nil {
		return 0, err
	}

	// Atualiza o cache com o saldo obtido
	r.cache.Set(ctx, CACHE_KEY, balance, time.Minute)

	return balance, nil
}
