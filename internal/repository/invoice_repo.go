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
	CreateInvoice(req CreateInvoiceRequest) (*domain.Invoice, error)
	GetInvoices() ([]domain.Invoice, error)
	GetByID(id string) (*domain.Invoice, error)
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
	if invoice.UserID != nil && *invoice.UserID != "" {
		userID = *invoice.UserID
	}

	query := `
	INSERT INTO invoices (id, user_id, invoice_code, client_email, status, total_amount, expired_at, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = tx.Exec(query,
		invoice.ID,
		userID,
		invoice.InvoiceCode,
		invoice.ClientEmail,
		invoice.Status,
		invoice.TotalAmount,
		invoice.ExpiredAt,
		invoice.CreatedAt,
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

func (u *invoiceRepository) CreateInvoice(req CreateInvoiceRequest) (*domain.Invoice, error) {
	invoiceID := uuid.New().String()

	var total int64
	var items []domain.InvoiceItem

	for _, i := range req.Items {
		total += int64(i.Qty) * i.Price

		name := i.Name
		price := float64(i.Price)

		items = append(items, domain.InvoiceItem{
			ID:        uuid.NewString(),
			InvoiceID: invoiceID,
			Name:      &name,
			Qty:       i.Qty,
			Price:     &price,
		})
	}

	invoiceCode := "INV-" + time.Now().Format("20060102150405")
	status := "pending"
	totalAmount := float64(total)
	expiredAt := time.Now().Add(24 * time.Hour)
	createdAt := time.Now()

	invoice := &domain.Invoice{
		ID:          invoiceID,
		UserID:      &req.UserID,
		InvoiceCode: &invoiceCode,
		ClientEmail: &req.ClientEmail,
		Status:      &status,
		TotalAmount: &totalAmount,
		ExpiredAt:   &expiredAt,
		CreatedAt:   &createdAt,
	}

	err := u.Create(invoice, items)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (u *invoiceRepository) GetInvoices() ([]domain.Invoice, error) {
	invoices := make([]domain.Invoice, 0)

	err := u.db.Select(&invoices, "SELECT * FROM invoices ORDER BY created_at DESC")
	return invoices, err
}

func (u *invoiceRepository) GetByID(id string) (*domain.Invoice, error) {
	invoice := &domain.Invoice{}
	err := u.db.Get(invoice, "SELECT * FROM invoices WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return invoice, nil
}
