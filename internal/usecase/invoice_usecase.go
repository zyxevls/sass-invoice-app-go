package usecase

import (
	"github.com/zyxevls/internal/domain"
	"github.com/zyxevls/internal/repository"
)

type InvoiceUsecase interface {
	CreateInvoice(req CreateInvoiceRequest) error
	GetInvoices() ([]domain.Invoice, error)
}

type invoiceUsecase struct {
	repo repository.InvoiceRepository
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

func NewInvoiceUsecase(r repository.InvoiceRepository) InvoiceUsecase {
	return &invoiceUsecase{r}
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

	return u.repo.CreateInvoice(repository.CreateInvoiceRequest{
		UserID:      req.UserID,
		ClientEmail: req.ClientEmail,
		Items:       repoItems,
	})
}

func (u *invoiceUsecase) GetInvoices() ([]domain.Invoice, error) {
	return u.repo.GetInvoices()
}
