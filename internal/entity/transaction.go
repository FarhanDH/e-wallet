package entity

import "time"

type Transaction struct {
	ID              int       `json:"id"`
	FromAccountID   int       `json:"from_account_id"`
	ToAccountID     int       `json:"to_account_id"`
	Amount          int       `json:"json"`
	TransactionType string    `json:"transaction_type"`
	CreatedAt       time.Time `json:"created_at"`
}
