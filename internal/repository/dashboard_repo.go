package repository

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	redisService "github.com/zyxevls/internal/infrastructure/redis"
)

type DashboardRepository interface {
	GetSummary() (map[string]interface{}, error)
	GetRevenueChart() ([]map[string]interface{}, error)
	GetInvoiceChart() ([]map[string]interface{}, error)
	GetTopCustomer() ([]map[string]interface{}, error)
	GetRecentTransaction() ([]map[string]interface{}, error)
}

type dashboardRepository struct {
	db  *sqlx.DB
	rdb *redis.Client
}

func NewDashboardRepository(db *sqlx.DB) DashboardRepository {
	return &dashboardRepository{db: db, rdb: nil}
}

func NewDashboardRepositoryWithRedis(db *sqlx.DB, rdb *redis.Client) DashboardRepository {
	return &dashboardRepository{db: db, rdb: rdb}
}

func (r *dashboardRepository) GetSummary() (map[string]interface{}, error) {
	cacheKey := "dashboard:summary"
	cacheDuration := 5 * time.Minute

	// Try to get from cache first
	if r.rdb != nil {
		var cached map[string]interface{}
		if err := redisService.CacheGet(r.rdb, cacheKey, &cached); err == nil {
			return cached, nil
		}
	}

	result := make(map[string]interface{})

	query := `
	SELECT 
	COUNT(*) as total_invoice,
	COALESCE(SUM(CASE WHEN status='paid' THEN total_amount ELSE 0 END), 0) AS total_revenue,
	COALESCE(SUM(CASE WHEN status='paid' THEN 1 ELSE 0 END), 0) AS total_paid,
	COALESCE(SUM(CASE WHEN status='pending' THEN 1 ELSE 0 END), 0) AS total_pending
	FROM invoices;
	`

	var totalInvoice, totalPaid, totalPending int64
	var totalRevenue float64

	err := r.db.QueryRow(query).Scan(
		&totalInvoice,
		&totalRevenue,
		&totalPaid,
		&totalPending,
	)

	if err != nil {
		return nil, err
	}

	result["total_invoice"] = totalInvoice
	result["total_revenue"] = totalRevenue
	result["total_paid"] = totalPaid
	result["total_pending"] = totalPending

	// Cache the result
	if r.rdb != nil {
		_ = redisService.CacheSet(r.rdb, cacheKey, result, cacheDuration)
	}

	return result, nil
}

func (r *dashboardRepository) GetRevenueChart() ([]map[string]interface{}, error) {
	cacheKey := "dashboard:revenue_chart"
	cacheDuration := 5 * time.Minute

	// Try to get from cache first
	if r.rdb != nil {
		var cached []map[string]interface{}
		if err := redisService.CacheGet(r.rdb, cacheKey, &cached); err == nil {
			return cached, nil
		}
	}

	rows, err := r.db.Query(`
	SELECT COALESCE(DATE(created_at)::text, ''), SUM(total_amount)
	FROM invoices
	WHERE status='paid'
	GROUP BY DATE(created_at)
	ORDER BY DATE(created_at)
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var date string
		var revenue int64

		if err := rows.Scan(&date, &revenue); err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"date":    date,
			"revenue": revenue,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Cache the result
	if r.rdb != nil {
		_ = redisService.CacheSet(r.rdb, cacheKey, result, cacheDuration)
	}

	return result, nil
}

func (r *dashboardRepository) GetInvoiceChart() ([]map[string]interface{}, error) {
	cacheKey := "dashboard:invoice_chart"
	cacheDuration := 5 * time.Minute

	// Try to get from cache first
	if r.rdb != nil {
		var cached []map[string]interface{}
		if err := redisService.CacheGet(r.rdb, cacheKey, &cached); err == nil {
			return cached, nil
		}
	}

	rows, err := r.db.Query(`
	SELECT COALESCE(DATE(created_at)::text, ''), COUNT(*)
	FROM invoices
	GROUP BY DATE(created_at)
	ORDER BY DATE(created_at)
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var date string
		var total int

		if err := rows.Scan(&date, &total); err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"date":  date,
			"total": total,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Cache the result
	if r.rdb != nil {
		_ = redisService.CacheSet(r.rdb, cacheKey, result, cacheDuration)
	}

	return result, nil
}

func (r *dashboardRepository) GetTopCustomer() ([]map[string]interface{}, error) {
	cacheKey := "dashboard:top_customer"
	cacheDuration := 5 * time.Minute

	// Try to get from cache first
	if r.rdb != nil {
		var cached []map[string]interface{}
		if err := redisService.CacheGet(r.rdb, cacheKey, &cached); err == nil {
			return cached, nil
		}
	}

	rows, err := r.db.Query(`
	SELECT client_email, COUNT(*) as total_invoice, SUM(total_amount) as total_revenue
	FROM invoices
	GROUP BY client_email
	ORDER BY total_revenue DESC
	LIMIT 10
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var clientEmail string
		var totalInvoice int
		var totalRevenue int64

		if err := rows.Scan(&clientEmail, &totalInvoice, &totalRevenue); err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"client_email":  clientEmail,
			"total_invoice": totalInvoice,
			"total_revenue": totalRevenue,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Cache the result
	if r.rdb != nil {
		_ = redisService.CacheSet(r.rdb, cacheKey, result, cacheDuration)
	}

	return result, nil
}

func (r *dashboardRepository) GetRecentTransaction() ([]map[string]interface{}, error) {
	cacheKey := "dashboard:recent_transaction"
	cacheDuration := 5 * time.Minute

	// Try to get from cache first
	if r.rdb != nil {
		var cached []map[string]interface{}
		if err := redisService.CacheGet(r.rdb, cacheKey, &cached); err == nil {
			return cached, nil
		}
	}

	rows, err := r.db.Query(`
	SELECT COALESCE(invoice_code, ''), COALESCE(client_email, ''), COALESCE(total_amount, 0), COALESCE(status, ''), COALESCE(created_at::text, '')
	FROM invoices
	ORDER BY created_at DESC
	LIMIT 10
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var invoiceCode string
		var clientEmail string
		var totalAmount int64
		var status string
		var createdAt string

		if err := rows.Scan(&invoiceCode, &clientEmail, &totalAmount, &status, &createdAt); err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{
			"invoice_code": invoiceCode,
			"client_email": clientEmail,
			"total_amount": totalAmount,
			"status":       status,
			"created_at":   createdAt,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Cache the result
	if r.rdb != nil {
		_ = redisService.CacheSet(r.rdb, cacheKey, result, cacheDuration)
	}

	return result, nil
}
