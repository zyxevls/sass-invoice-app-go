package email

import (
	"github.com/zyxevls/internal/config"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{cfg}
}

func (e *EmailService) Send(to, subject, body string, attachment string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", e.cfg.SMTPEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if attachment != "" {
		m.Attach(attachment)
	}

	d := gomail.NewDialer(e.cfg.SMTPHost, e.cfg.SMTPPort, e.cfg.SMTPEmail, e.cfg.SMTPPass)

	return d.DialAndSend(m)
}
