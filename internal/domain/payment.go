package domain

import "time"

type Payment struct {
	ID              string     `db:"id" json:"id"`
	InvoiceID       string     `db:"invoice_id" json:"invoice_id"`
	MidtransOrderID string     `db:"midtrans_order_id" json:"midtrans_order_id"`
	Status          string     `db:"status" json:"status"`
	PaymentURL      string     `db:"payment_url" json:"payment_url"`
	PaidAt          *time.Time `db:"paid_at" json:"paid_at,omitempty"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
}
