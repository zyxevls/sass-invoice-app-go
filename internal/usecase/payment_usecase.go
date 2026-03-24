package usecase

import (
	"github.com/google/uuid"
	"github.com/zyxevls/internal/domain"
	"github.com/zyxevls/internal/infrastructure/email"
	"github.com/zyxevls/internal/infrastructure/midtrans"
	"github.com/zyxevls/internal/infrastructure/pdf"
	"github.com/zyxevls/internal/repository"
)

type PaymentUseCase interface {
	CreatePayment(invoiceID string, email string, amount int64) (string, error)
	HandleWebHook(orderID string, status string) error
}

type paymentUseCase struct {
	repo        repository.PaymentRepository
	invoiceRepo repository.InvoiceRepository
	mid         *midtrans.MidtransService
	email       *email.EmailService
	pdf         *pdf.PDFService
}

func NewPaymentUseCase(r repository.PaymentRepository, ir repository.InvoiceRepository, m *midtrans.MidtransService, e *email.EmailService, p *pdf.PDFService) PaymentUseCase {
	return &paymentUseCase{r, ir, m, e, p}
}

func (u *paymentUseCase) CreatePayment(invoiceID string, email string, amount int64) (string, error) {
	orderID := "ORDER-" + uuid.NewString()

	url, err := u.mid.CreateTransaction(orderID, amount, email)
	if err != nil {
		return "", err
	}

	payment := &domain.Payment{
		ID:              uuid.NewString(),
		InvoiceID:       invoiceID,
		MidtransOrderID: orderID,
		Status:          "pending",
		PaymentURL:      url,
	}

	if err := u.repo.Create(payment); err != nil {
		return "", err
	}

	return url, nil
}

func (u *paymentUseCase) HandleWebHook(orderID string, status string) error {
	err := u.repo.UpdateStatus(orderID, status)
	if err != nil {
		return err
	}

	if status == "settlement" {
		payment, err := u.repo.GetByOrderID(orderID)
		if err != nil {
			return err
		}

		invoice, err := u.invoiceRepo.GetByID(payment.InvoiceID)
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

		clientEmail := ""
		if invoice.ClientEmail != nil {
			clientEmail = *invoice.ClientEmail
		}

		pdfFile, err := u.pdf.GenerateInvoice(invoiceCode, clientEmail, totalAmount)
		if err != nil {
			return err
		}

		emailBody := email.InvoiceTemplate(invoiceCode, totalAmount, payment.PaymentURL)

		err = u.email.Send(clientEmail, "Payment Success - Invoice #"+invoiceCode, emailBody, pdfFile)
		if err != nil {
			return err
		}
	}

	return nil
}
