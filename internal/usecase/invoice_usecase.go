package usecase

import (
	"github.com/zyxevls/internal/config"
	"github.com/zyxevls/internal/domain"
	"github.com/zyxevls/internal/infrastructure/email"
	"github.com/zyxevls/internal/infrastructure/midtrans"
	"github.com/zyxevls/internal/infrastructure/pdf"
	"github.com/zyxevls/internal/repository"
)

type InvoiceUsecase interface {
	CreateInvoice(req CreateInvoiceRequest) error
	GetInvoices() ([]domain.Invoice, error)
	GetInvoiceByID(id string) (*domain.Invoice, error)
}

type invoiceUsecase struct {
	repo     repository.InvoiceRepository
	email    *email.EmailService
	pdf      *pdf.PDFService
	midtrans *midtrans.MidtransService
}

type CreateInvoiceRequest struct {
	UserID      string             `json:"user_id"`
	ClientEmail string             `json:"client_email"`
	Items       []InvoiceItemInput `json:"items"`
}

type InvoiceItemInput struct {
	Name  string `json:"name"`
	Qty   int    `json:"qty"`
	Price int64  `json:"price"`
}

func NewInvoiceUsecase(r repository.InvoiceRepository, cfg *config.Config, m *midtrans.MidtransService) InvoiceUsecase {
	return &invoiceUsecase{
		repo:     r,
		email:    email.NewEmailService(cfg),
		pdf:      pdf.NewPDFService(),
		midtrans: m,
	}
}

func (u *invoiceUsecase) CreateInvoice(req CreateInvoiceRequest) error {
	repoItems := make([]repository.CreateInvoiceItemRequest, 0, len(req.Items))
	for _, item := range req.Items {
		repoItems = append(repoItems, repository.CreateInvoiceItemRequest{
			Name:  item.Name,
			Qty:   item.Qty,
			Price: item.Price,
		})
	}

	invoice, err := u.repo.CreateInvoice(repository.CreateInvoiceRequest{
		UserID:      req.UserID,
		ClientEmail: req.ClientEmail,
		Items:       repoItems,
	})
	if err != nil {
		return err
	}

	totalAmount := int64(0)
	if invoice.TotalAmount != nil {
		totalAmount = int64(*invoice.TotalAmount)
	}

	invoiceCode := ""
	if invoice.InvoiceCode != nil {
		invoiceCode = *invoice.InvoiceCode
	}

	paymentURL, err := u.midtrans.CreateTransaction(invoice.ID, totalAmount, req.ClientEmail)
	if err != nil {
		return err
	}

	pdfFile, err := u.pdf.GenerateInvoice(invoiceCode, req.ClientEmail, totalAmount)
	if err != nil {
		return err
	}

	emailBody := email.InvoiceTemplate(invoiceCode, totalAmount, paymentURL)

	err = u.email.Send(req.ClientEmail, "Invoice #"+invoiceCode, emailBody, pdfFile)
	if err != nil {
		return err
	}

	return nil
}

func (u *invoiceUsecase) GetInvoices() ([]domain.Invoice, error) {
	return u.repo.GetInvoices()
}

func (u *invoiceUsecase) GetInvoiceByID(id string) (*domain.Invoice, error) {
	return u.repo.GetByID(id)
}
