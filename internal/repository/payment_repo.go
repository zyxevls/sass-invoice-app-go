package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/zyxevls/internal/domain"
)

type PaymentRepository interface {
	Create(p *domain.Payment) error
	UpdateStatus(orderID string, status string) error
	GetByOrderID(orderID string) (*domain.Payment, error)
}

type paymentRepository struct {
	db *sqlx.DB
}

type CreatePaymentRequest struct {
	OrderID string
	Amount  int64
}

func NewPaymentRepository(db *sqlx.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) Create(p *domain.Payment) error {
	_, err := r.db.Exec(`
	INSERT INTO payments (id, invoice_id, midtrans_order_id, status, payment_url)
	VALUES ($1, $2, $3, $4, $5)`,
		p.ID,
		p.InvoiceID,
		p.MidtransOrderID,
		p.Status,
		p.PaymentURL,
	)
	return err
}

func (r *paymentRepository) UpdateStatus(orderID string, status string) error {
	_, err := r.db.Exec(`
	UPDATE payments SET status=$1, paid_at=NOW()
	WHERE midtrans_order_id=$2`,
		status,
		orderID,
	)
	if err != nil {
		return err
	}
	return nil
}
func (r *paymentRepository) GetByOrderID(orderID string) (*domain.Payment, error) {
	payment := &domain.Payment{}
	err := r.db.Get(payment, "SELECT * FROM payments WHERE midtrans_order_id=$1", orderID)
	if err != nil {
		return nil, err
	}
	return payment, nil
}
