package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/zyxevls/internal/domain"
)

type InvoiceRepository interface {
	Create(invoice *domain.Invoice, items []domain.InvoiceItem) error
	FindAll() ([]domain.Invoice, error)
	CreateInvoice(req CreateInvoiceRequest) error
	GetInvoices() ([]domain.Invoice, error)
}

type invoiceRepository struct {
	db *sqlx.DB
}

type CreateInvoiceItemRequest struct {
	Name  string
	Qty   int
	Price int64
}

type CreateInvoiceRequest struct {
	UserID      string
	ClientEmail string
	Items       []CreateInvoiceItemRequest
}

func NewInvoiceRepository(db *sqlx.DB) InvoiceRepository {
	return &invoiceRepository{db}
}

func (r *invoiceRepository) Create(invoice *domain.Invoice, items []domain.InvoiceItem) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	var userID any
	if invoice.UserID != "" {
		userID = invoice.UserID
	}

	query := `
	INSERT INTO invoices (id, user_id, invoice_code, client_email, status, total_amount, expired_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = tx.Exec(query,
		invoice.ID,
		userID,
		invoice.InvoiceCode,
		invoice.ClientEmail,
		invoice.Status,
		invoice.TotalAmount,
		invoice.ExpiredAt,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	for _, item := range items {
		_, err := tx.Exec(`
		INSERT INTO invoice_items (id, invoice_id, name, qty, price)
		VALUES ($1, $2, $3, $4, $5)
		`,
			item.ID,
			item.InvoiceID,
			item.Name,
			item.Qty,
			item.Price,
		)

		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *invoiceRepository) FindAll() ([]domain.Invoice, error) {
	invoices := make([]domain.Invoice, 0)

	err := r.db.Select(&invoices, "SELECT * FROM invoices ORDER BY created_at DESC")
	return invoices, err
}

func (u *invoiceRepository) CreateInvoice(req CreateInvoiceRequest) error {
	invoiceID := uuid.New().String()

	var total int64
	var items []domain.InvoiceItem

	for _, i := range req.Items {
		total += int64(i.Qty) * i.Price

		items = append(items, domain.InvoiceItem{
			ID:        uuid.NewString(),
			InvoiceID: invoiceID,
			Name:      i.Name,
			Qty:       i.Qty,
			Price:     float64(i.Price),
		})
	}

	invoice := &domain.Invoice{
		ID:          invoiceID,
		UserID:      req.UserID,
		InvoiceCode: "INV-" + time.Now().Format("20060102150405"),
		ClientEmail: req.ClientEmail,
		Status:      "pending",
		TotalAmount: float64(total),
		ExpiredAt:   time.Now().Add(24 * time.Hour),
		CreatedAt:   time.Now(),
	}

	return u.Create(invoice, items)
}

func (u *invoiceRepository) GetInvoices() ([]domain.Invoice, error) {
	return u.FindAll()
}
