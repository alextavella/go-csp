package domain

type Transaction struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
}
