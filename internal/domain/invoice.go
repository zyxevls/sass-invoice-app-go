package domain

import "time"

type Invoice struct {
	ID          string    `db:"id" json:"id"`
	UserID      *string   `db:"user_id" json:"user_id"`
	InvoiceCode *string    `db:"invoice_code" json:"invoice_code"`
	ClientEmail *string    `db:"client_email" json:"client_email"`
	Status      *string    `db:"status" json:"status"`
	TotalAmount *float64   `db:"total_amount" json:"total_amount"`
	ExpiredAt   *time.Time `db:"expired_at" json:"expired_at"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at"`
}

type InvoiceItem struct {
	ID        string  `db:"id" json:"id"`
	InvoiceID string  `db:"invoice_id" json:"invoice_id"`
	Name      *string `db:"name" json:"name"`
	Qty       int     `db:"qty" json:"qty"`
	Price     *float64 `db:"price" json:"price"`
}
