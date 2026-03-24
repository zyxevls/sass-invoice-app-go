package usecase

import "github.com/zyxevls/internal/repository"

type DashboardUsecase interface {
	GetDashboard() (map[string]interface{}, error)
	GetTopCustomer() ([]map[string]interface{}, error)
	GetRecentTransaction() ([]map[string]interface{}, error)
}

type dashboardUsecase struct {
	repo repository.DashboardRepository
}

func NewDashboardUsecase(r repository.DashboardRepository) DashboardUsecase {
	return &dashboardUsecase{r}
}

func (u *dashboardUsecase) GetDashboard() (map[string]interface{}, error) {
	summary, err := u.repo.GetSummary()
	if err != nil {
		return nil, err
	}
	revenue, err := u.repo.GetRevenueChart()
	if err != nil {
		return nil, err
	}
	invoice, err := u.repo.GetInvoiceChart()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"summary":       summary,
		"revenue_chart": revenue,
		"invoice_chart": invoice,
	}, nil
}

func (u *dashboardUsecase) GetTopCustomer() ([]map[string]interface{}, error) {
	return u.repo.GetTopCustomer()
}

func (u *dashboardUsecase) GetRecentTransaction() ([]map[string]interface{}, error) {
	return u.repo.GetRecentTransaction()
}
