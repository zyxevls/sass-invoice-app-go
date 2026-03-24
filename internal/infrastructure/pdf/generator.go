package pdf

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

func (p *PDFService) GenerateInvoice(code string, email string, amount int64) (string, error) {
	fileName := fmt.Sprintf("invoice_%s.pdf", code)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Invoice")

	// Invoice Number
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Invoice #%s", code))

	// Amount
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(40, 10, fmt.Sprintf("Amount: Rp %d", amount))

	// Save PDF
	err := pdf.OutputFileAndClose(fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
