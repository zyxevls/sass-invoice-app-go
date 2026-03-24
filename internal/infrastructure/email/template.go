package email

import (
	"fmt"

	"github.com/zyxevls/internal/helpers"
)

func InvoiceTemplate(code string, amount int64, paymentURL string) string {
	rupiah := helpers.FormatRupiah(amount)
	return fmt.Sprintf(`
	<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Invoice #%[1]s</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f4f4f4;
        }
        .container {
            max-width: 600px;
            margin: 20px auto;
            background-color: #ffffff;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        .header {
            background-color: #4CAF50;
            color: white;
            padding: 20px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 24px;
        }
        .content {
            padding: 30px;
        }
        .invoice-details {
            margin-bottom: 20px;
            background-color: #e8f5e9;
            padding: 15px;
            border-radius: 5px;
        }
        .invoice-details p {
            margin: 5px 0;
            font-size: 16px;
        }
        .amount {
            font-size: 28px;
            font-weight: bold;
            color: #4CAF50;
            margin: 20px 0;
            text-align: center;
        }
        .button-container {
            text-align: center;
            margin: 30px 0;
        }
        .button {
            background-color: #4CAF50;
            color: white;
            padding: 15px 25px;
            text-decoration: none;
            border-radius: 5px;
            font-size: 18px;
            font-weight: bold;
        }
        .footer {
            background-color: #f1f1f1;
            padding: 20px;
            text-align: center;
            font-size: 12px;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Invoice #%[1]s</h1>
        </div>
        <div class="content">
            <div class="invoice-details">
                <p><strong>Invoice Code:</strong> %[1]s</p>
                <p><strong>Amount:</strong> %[2]s</p>
            </div>
            <div class="amount">
                Total: %[2]s
            </div>
            <div class="button-container">
                <a href="%[3]s" class="button">Pay Now</a>
            </div>
            <p style="text-align: center; color: #666;">Click the button above to complete your payment.</p>
        </div>
        <div class="footer">
            <p>This is an automated email. Please do not reply.</p>
        </div>
    </div>
</body>
</html>
	`,
		code, rupiah, paymentURL)
}
