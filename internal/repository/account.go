package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/alextavella/bank-api/internal/config"
	"github.com/alextavella/bank-api/internal/domain"
	"github.com/bsm/redislock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

type AccountRepository struct {
	db     *sql.DB
	cache  *redis.Client
	locker *redislock.Client
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
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.CACHE_URI,
		Password: "", // Sem senha
		DB:       0,  // Banco padrão
	})

	locker := redislock.New(redisClient)

	return &AccountRepository{db: db, cache: redisClient, locker: locker}, nil
}

func (r *AccountRepository) Deposit(ctx context.Context, transaction domain.Transaction) error {
	amount := transaction.Amount
	userID := transaction.UserID

	lockKey := "lock:user:" + userID
	lock, err := r.locker.Obtain(ctx, lockKey, 5*time.Second, nil)
	if err == redislock.ErrNotObtained {
		return fmt.Errorf("outra operação em andamento para este usuário")
	} else if err != nil {
		return fmt.Errorf("erro ao tentar obter lock")
	}
	defer lock.Release(ctx)

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var balance float64
	// Lock explícito para garantir consistência
	err = tx.QueryRow("SELECT balance FROM accounts WHERE user_id = ? FOR UPDATE", userID).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			if amount < 0 {
				return fmt.Errorf("saldo insuficiente")
			}
			_, err = tx.Exec("INSERT INTO accounts (user_id, balance) VALUES (?, ?)", userID, amount)
			if err != nil {
				return fmt.Errorf("erro ao persistir operação")
			}
		} else {
			return err
		}
	}

	newBalance := balance + amount
	_, err = tx.Exec("UPDATE accounts SET balance = ? WHERE user_id = ?", newBalance, userID)
	if err != nil {
		return err
	}

	// Atualiza o saldo no cache
	r.cache.Set(ctx, CACHE_KEY, newBalance, time.Minute)

	return tx.Commit()
}

func (r *AccountRepository) Withdraw(ctx context.Context, transaction domain.Transaction) error {
	amount := transaction.Amount
	userID := transaction.UserID

	lockKey := "lock:user:" + userID
	lock, err := r.locker.Obtain(ctx, lockKey, 5*time.Second, nil)
	if err == redislock.ErrNotObtained {
		return fmt.Errorf("outra operação em andamento para este usuário")
	} else if err != nil {
		return fmt.Errorf("erro ao tentar obter lock")
	}
	defer lock.Release(ctx)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação")
	}
	defer tx.Rollback()

	var balance float64
	// Lock explícito para garantir consistência
	err = tx.QueryRow("SELECT balance FROM accounts WHERE user_id = ? FOR UPDATE", userID).Scan(&balance)
	if err != nil {
		return fmt.Errorf("saldo insuficiente")
	}

	if balance < amount {
		return fmt.Errorf("saldo insuficiente")
	}

	newBalance := balance - amount
	_, err = tx.Exec("UPDATE accounts SET balance = ? WHERE user_id = ?", newBalance, userID)
	if err != nil {
		return fmt.Errorf("erro ao persistir operação")
	}

	// Atualiza o saldo no cache
	r.cache.Set(ctx, CACHE_KEY, newBalance, time.Minute)

	return tx.Commit()
}

func (r *AccountRepository) GetBalance(ctx context.Context, userID string) (float64, error) {
	var balance float64

	// Tenta buscar no cache
	balance, err := r.cache.Get(ctx, "balance").Float64()
	if err == nil {
		return balance, nil
	}

	// Se não estiver no cache, busca no banco
	err = r.db.QueryRow("SELECT balance FROM accounts WHERE user_id = ?", userID).Scan(&balance)
	if err != nil {
		return 0, err
	}

	// Atualiza o cache com o saldo obtido
	r.cache.Set(ctx, CACHE_KEY, balance, time.Minute)

	return balance, nil
}
