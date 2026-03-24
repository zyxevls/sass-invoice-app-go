# SaaS Invoice Application

A robust Go-based SaaS invoice management system featuring automated PDF generation, email notifications, and integrated Midtrans payment gateway.

## 🚀 Features

- **Invoice Management**: Create and track invoices with automated status updates.
- **Payment Integration**: Seamless payment processing using Midtrans Snap API.
- **PDF Generation**: Automatic generation of invoice PDFs with custom styling.
- **Email Notifications**: Instant email alerts for new invoices and successful payments, including PDF attachments.
- **Webhooks**: Secure handling of payment status notifications from Midtrans.
- **Clean Architecture**: Organized into layered components (Domain, Usecase, Repository, Infrastructure, Delivery).

## 🛠️ Tech Stack

- **Language**: [Go (Golang)](https://go.dev/)
- **Web Framework**: [Fiber v3](https://gofiber.io/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **SQL Library**: [sqlx](https://github.com/jmoiron/sqlx)
- **Payment Gateway**: [Midtrans SDK](https://github.com/midtrans/midtrans-go)
- **PDF Library**: [gofpdf](https://github.com/jung-kurt/gofpdf)
- **Email Service**: [gomail.v2](https://gopkg.in/gomail.v2)
- **Configuration**: [Viper](https://github.com/spf13/viper)
- **ID Generation**: [UUID](https://github.com/google/uuid)

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- SMTP Email account (for notifications)
- Midtrans Sandbox/Production account

## ⚙️ Setup & Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/zyxevls/sass-invoice-app-go.git
   cd sass-invoice-app-go
   ```

2. **Configuration**:
   Rename the `.env.example` (if exists) or create a `.env` file in the root directory:
   ```env
   DATABASE_URL=postgres://user:password@localhost:5432/invoice_db?sslmode=disable
   
   MIDTRANS_SERVER_KEY=your_server_key
   MIDTRANS_CLIENT_KEY=your_client_key
   MIDTRANS_IS_PRODUCTION=false

   SMTP_HOST=smtp.gmail.com
   SMTP_PORT=587
   SMTP_EMAIL=your_email@gmail.com
   SMTP_PASS=your_app_password
   ```

3. **Install dependencies**:
   ```bash
   go mod tidy
   ```

4. **Run the application**:
   ```bash
   go run cmd/main.go
   ```

## 📡 API Documentation

### Create Invoice
- **Endpoint**: `POST /api/v1/invoices`
- **Request Body**:
  ```json
  {
    "user_id": "user-uuid",
    "client_email": "client@example.com",
    "items": [
      {
        "name": "Web Hosting (1 Year)",
        "qty": 1,
        "price": 1500000
      }
    ]
  }
  ```

### Get All Invoices
- **Endpoint**: `GET /api/v1/invoices`

### Midtrans Webhook
- **Endpoint**: `POST /api/v1/payments/webhook`
- Handles `settlement`, `pending`, `expire`, and `cancel` status updates.

## 📁 Project Structure

```text
.
├── cmd/                # Entry point
├── internal/
│   ├── config/         # Configuration loading
│   ├── delivery/       # HTTP handlers & routing
│   ├── domain/         # Data models
│   ├── helpers/        # Utility functions (formatting, etc.)
│   ├── infrastructure/ # External services (Email, PDF, Midtrans)
│   ├── repository/     # Database operations
│   └── usecase/        # Business logic
├── .env                # Environment variables
└── go.mod              # Dependencies
```

## 📄 License
This project is licensed under the MIT License - see the LICENSE file for details.
