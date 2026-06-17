package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PaymentService struct {
	db *sql.DB
}

func NewPaymentService(db *sql.DB) *PaymentService {
	return &PaymentService{db: db}
}

type PaymentResponse struct {
	ID                   string    `json:"id"`
	OrderID              string    `json:"orderId"`
	Amount               float64   `json:"amount"`
	PaymentFee           float64   `json:"paymentFee"`
	PaymentMethod        string    `json:"paymentMethod"`
	GatewayTransactionID string    `json:"gatewayTransactionId"`
	Status               string    `json:"status"`
	FailureMessage       string    `json:"failureMessage"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

type PaymentProviderResponse struct {
	ID                       string `json:"id"`
	Name                     string `json:"name"`
	IsEnabled                bool   `json:"isEnabled"`
	ConfigureURL             string `json:"configureUrl"`
	LandingViewComponentName string `json:"landingViewComponentName"`
	AdditionalSettings       string `json:"additionalSettings"`
}

func (s *PaymentService) CreatePayment(ctx context.Context, payment *PaymentResponse) error {
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO payments_payments (order_id, amount, payment_fee, payment_method, gateway_transaction_id, status, failure_message, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())
		RETURNING id, created_at, updated_at`,
		payment.OrderID, payment.Amount, payment.PaymentFee, payment.PaymentMethod,
		payment.GatewayTransactionID, payment.Status, payment.FailureMessage,
	).Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	return nil
}

func (s *PaymentService) GetPayment(ctx context.Context, id string) (*PaymentResponse, error) {
	var p PaymentResponse
	err := s.db.QueryRowContext(ctx, `
		SELECT id, order_id, amount, payment_fee, payment_method, gateway_transaction_id, status, failure_message, created_at, updated_at
		FROM payments_payments WHERE id = $1
	`, id).Scan(&p.ID, &p.OrderID, &p.Amount, &p.PaymentFee, &p.PaymentMethod,
		&p.GatewayTransactionID, &p.Status, &p.FailureMessage, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	return &p, nil
}

func (s *PaymentService) GetAllPayments(ctx context.Context, page, pageSize int) ([]PaymentResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var total int64
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM payments_payments`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count payments: %w", err)
	}

	offset := (page - 1) * pageSize
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, order_id, amount, payment_fee, payment_method, gateway_transaction_id, status, failure_message, created_at, updated_at
		FROM payments_payments ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list payments: %w", err)
	}
	defer rows.Close()

	var payments []PaymentResponse
	for rows.Next() {
		var p PaymentResponse
		if err := rows.Scan(&p.ID, &p.OrderID, &p.Amount, &p.PaymentFee, &p.PaymentMethod,
			&p.GatewayTransactionID, &p.Status, &p.FailureMessage, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		payments = append(payments, p)
	}
	if payments == nil {
		payments = []PaymentResponse{}
	}
	return payments, total, nil
}

func (s *PaymentService) GetPaymentsByOrderID(ctx context.Context, orderID string) ([]PaymentResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, order_id, amount, payment_fee, payment_method, gateway_transaction_id, status, failure_message, created_at, updated_at
		FROM payments_payments WHERE order_id = $1 ORDER BY created_at DESC
	`, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %w", err)
	}
	defer rows.Close()

	var payments []PaymentResponse
	for rows.Next() {
		var p PaymentResponse
		if err := rows.Scan(&p.ID, &p.OrderID, &p.Amount, &p.PaymentFee, &p.PaymentMethod,
			&p.GatewayTransactionID, &p.Status, &p.FailureMessage, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	if payments == nil {
		payments = []PaymentResponse{}
	}
	return payments, nil
}

func (s *PaymentService) ListProviders(ctx context.Context) ([]PaymentProviderResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, is_enabled, configure_url, landing_view_component_name, additional_settings
		FROM payments_payment_providers ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list providers: %w", err)
	}
	defer rows.Close()

	var providers []PaymentProviderResponse
	for rows.Next() {
		var p PaymentProviderResponse
		if err := rows.Scan(&p.ID, &p.Name, &p.IsEnabled, &p.ConfigureURL, &p.LandingViewComponentName, &p.AdditionalSettings); err != nil {
			return nil, err
		}
		providers = append(providers, p)
	}
	if providers == nil {
		providers = []PaymentProviderResponse{}
	}
	return providers, nil
}
