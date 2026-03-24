package midtrans

import (
	"log"
	"strings"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/zyxevls/internal/config"
)

type MidtransService struct {
	client snap.Client
}

func NewMidtransService(cfg *config.Config) *MidtransService {
	var env midtrans.EnvironmentType

	if cfg.MidtransIsProduction {
		env = midtrans.Production
		if strings.HasPrefix(cfg.MidtransServerKey, "SB-") {
			log.Println("WARNING: Midtrans is set to Production but Server Key appears to be for Sandbox (starts with SB-)")
		}
	} else {
		env = midtrans.Sandbox
		if !strings.HasPrefix(cfg.MidtransServerKey, "SB-") {
			log.Println("WARNING: Midtrans is set to Sandbox but Server Key is missing 'SB-' prefix. This will likely cause 401 Unauthorized errors.")
		}
	}

	var s snap.Client
	s.New(cfg.MidtransServerKey, env)

	return &MidtransService{s}
}

func (m *MidtransService) CreateTransaction(orderID string, amount int64, email string) (string, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			Email: email,
		},
	}
	resp, err := m.client.CreateTransaction(req)
	if err != nil {
		return "", err
	}

	return resp.RedirectURL, nil
}
