package usecase

import (
	"github.com/google/uuid"
	"github.com/zyxevls/internal/domain"
	"github.com/zyxevls/internal/infrastructure/midtrans"
	"github.com/zyxevls/internal/repository"
)

type PaymentUseCase interface {
	CreatePayment(invoiceID string, email string, amount int64) (string, error)
	HandleWebHook(orderID string, status string) error
}

type paymentUseCase struct {
	repo repository.PaymentRepository
	mid  *midtrans.MidtransService
}

func NewPaymentUseCase(r repository.PaymentRepository, m *midtrans.MidtransService) PaymentUseCase {
	return &paymentUseCase{r, m}
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
	return u.repo.UpdateStatus(orderID, status)
}
