package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type TaxService struct {
	db *sql.DB
}

func NewTaxService(db *sql.DB) *TaxService {
	return &TaxService{db: db}
}

type TaxClassResponse struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
	Rates     []TaxRateResponse `json:"rates,omitempty"`
}

type TaxRateResponse struct {
	ID                string    `json:"id"`
	TaxClassID        string    `json:"taxClassId"`
	CountryID         string    `json:"countryId"`
	StateOrProvinceID *string   `json:"stateOrProvinceId,omitempty"`
	Rate              float64   `json:"rate"`
	ZipCode           string    `json:"zipCode"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func (s *TaxService) GetTaxClasses(ctx context.Context) ([]TaxClassResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, created_at, updated_at FROM tax_tax_classes ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list tax classes: %w", err)
	}
	defer rows.Close()

	var classes []TaxClassResponse
	for rows.Next() {
		var tc TaxClassResponse
		if err := rows.Scan(&tc.ID, &tc.Name, &tc.CreatedAt, &tc.UpdatedAt); err != nil {
			return nil, err
		}
		classes = append(classes, tc)
	}
	if classes == nil {
		classes = []TaxClassResponse{}
	}
	return classes, nil
}

func (s *TaxService) CreateTaxClass(ctx context.Context, name string) (*TaxClassResponse, error) {
	var tc TaxClassResponse
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO tax_tax_classes (name, created_at, updated_at) VALUES ($1,NOW(),NOW())
		RETURNING id, name, created_at, updated_at
	`, name).Scan(&tc.ID, &tc.Name, &tc.CreatedAt, &tc.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create tax class: %w", err)
	}
	tc.Rates = []TaxRateResponse{}
	return &tc, nil
}

func (s *TaxService) UpdateTaxClass(ctx context.Context, id string, name string) (*TaxClassResponse, error) {
	var tc TaxClassResponse
	err := s.db.QueryRowContext(ctx, `
		UPDATE tax_tax_classes SET name=$1, updated_at=NOW() WHERE id=$2
		RETURNING id, name, created_at, updated_at
	`, name, id).Scan(&tc.ID, &tc.Name, &tc.CreatedAt, &tc.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tax class not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update tax class: %w", err)
	}
	return &tc, nil
}

func (s *TaxService) DeleteTaxClass(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM tax_tax_classes WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete tax class: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("tax class not found")
	}
	return nil
}

func (s *TaxService) GetTaxRates(ctx context.Context, page, pageSize int) ([]TaxRateResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tax_tax_rates`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tax rates: %w", err)
	}

	offset := (page - 1) * pageSize
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, tax_class_id, country_id, state_or_province_id, rate, COALESCE(zip_code,''), created_at, updated_at
		FROM tax_tax_rates ORDER BY country_id ASC LIMIT $1 OFFSET $2
	`, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tax rates: %w", err)
	}
	defer rows.Close()

	var rates []TaxRateResponse
	for rows.Next() {
		var tr TaxRateResponse
		var stateID sql.NullString
		if err := rows.Scan(&tr.ID, &tr.TaxClassID, &tr.CountryID, &stateID, &tr.Rate, &tr.ZipCode, &tr.CreatedAt, &tr.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if stateID.Valid {
			tr.StateOrProvinceID = &stateID.String
		}
		rates = append(rates, tr)
	}
	if rates == nil {
		rates = []TaxRateResponse{}
	}
	return rates, total, nil
}

func (s *TaxService) CreateTaxRate(ctx context.Context, taxClassID string, countryID string, stateOrProvinceID *string, rate float64, zipCode string) (*TaxRateResponse, error) {
	var tr TaxRateResponse
	var stateID sql.NullString
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO tax_tax_rates (tax_class_id, country_id, state_or_province_id, rate, zip_code, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,NOW(),NOW())
		RETURNING id, tax_class_id, country_id, state_or_province_id, rate, zip_code, created_at, updated_at
	`, taxClassID, countryID, stateOrProvinceID, rate, zipCode).Scan(
		&tr.ID, &tr.TaxClassID, &tr.CountryID, &stateID, &tr.Rate, &tr.ZipCode, &tr.CreatedAt, &tr.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create tax rate: %w", err)
	}
	if stateID.Valid {
		tr.StateOrProvinceID = &stateID.String
	}
	return &tr, nil
}

func (s *TaxService) UpdateTaxRate(ctx context.Context, id string, taxClassID *string, countryID *string, stateOrProvinceID *string, rate *float64, zipCode *string) (*TaxRateResponse, error) {
	existing, err := s.GetTaxRate(ctx, id)
	if err != nil {
		return nil, err
	}

	tcID := existing.TaxClassID
	if taxClassID != nil {
		tcID = *taxClassID
	}
	cID := existing.CountryID
	if countryID != nil {
		cID = *countryID
	}
	sopID := existing.StateOrProvinceID
	if stateOrProvinceID != nil {
		sopID = stateOrProvinceID
	}
	r := existing.Rate
	if rate != nil {
		r = *rate
	}
	z := existing.ZipCode
	if zipCode != nil {
		z = *zipCode
	}

	var tr TaxRateResponse
	var stateID sql.NullString
	err = s.db.QueryRowContext(ctx, `
		UPDATE tax_tax_rates SET tax_class_id=$1, country_id=$2, state_or_province_id=$3, rate=$4, zip_code=$5, updated_at=NOW()
		WHERE id=$6 RETURNING id, tax_class_id, country_id, state_or_province_id, rate, zip_code, created_at, updated_at
	`, tcID, cID, sopID, r, z, id).Scan(
		&tr.ID, &tr.TaxClassID, &tr.CountryID, &stateID, &tr.Rate, &tr.ZipCode, &tr.CreatedAt, &tr.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tax rate not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update tax rate: %w", err)
	}
	if stateID.Valid {
		tr.StateOrProvinceID = &stateID.String
	}
	return &tr, nil
}

func (s *TaxService) DeleteTaxRate(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM tax_tax_rates WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete tax rate: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("tax rate not found")
	}
	return nil
}

func (s *TaxService) GetTaxRate(ctx context.Context, id string) (*TaxRateResponse, error) {
	var tr TaxRateResponse
	var stateID sql.NullString
	err := s.db.QueryRowContext(ctx, `
		SELECT id, tax_class_id, country_id, state_or_province_id, rate, COALESCE(zip_code,''), created_at, updated_at
		FROM tax_tax_rates WHERE id = $1
	`, id).Scan(&tr.ID, &tr.TaxClassID, &tr.CountryID, &stateID, &tr.Rate, &tr.ZipCode, &tr.CreatedAt, &tr.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tax rate not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get tax rate: %w", err)
	}
	if stateID.Valid {
		tr.StateOrProvinceID = &stateID.String
	}
	return &tr, nil
}
