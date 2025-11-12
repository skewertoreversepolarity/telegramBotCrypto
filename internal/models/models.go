package models

import (
	"time"
)

// Deposit представляет запись из таблицы deposit
type Deposit struct {
	ID          int       `json:"id" db:"id"`
	Network     string    `json:"network" db:"network"`
	Currency    string    `json:"currency" db:"currency"`
	FromAddress string    `json:"from_address" db:"from_address"`
	ToAddress   string    `json:"to_address" db:"to_address"`
	Amount      float64   `json:"amount" db:"amount"`
	TxnHash     string    `json:"txn_hash" db:"txn_hash"`
	BlockNumber int64     `json:"block_number" db:"block_number"`
	Result      string    `json:"result" db:"result"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	TranCode    *string   `json:"tran_code" db:"tran_code"`
}

// Outgoing представляет запись из таблицы outgoing
type Outgoing struct {
	ID          int       `json:"id" db:"id"`
	From        string    `json:"from" db:"from"`
	To          string    `json:"to" db:"to"`
	Currency    string    `json:"currency" db:"currency"`
	Network     string    `json:"network" db:"network"`
	Amount      float64   `json:"amount" db:"amount"`
	Commission  float64   `json:"commission" db:"commission"`
	TranCode    *string   `json:"trancode" db:"trancode"`
	TranHash    *string   `json:"tran_hash" db:"tran_hash"`
	Status      string    `json:"status" db:"status"`
	Response    *string   `json:"response" db:"response"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	TotalAmount float64   `json:"total_amount" db:"total_amount"`
	Finoper     *int      `json:"finoper" db:"finoper"`
}
