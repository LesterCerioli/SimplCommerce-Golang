package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type CountryService struct {
	db *sql.DB
}

func NewCountryService(db *sql.DB) *CountryService {
	return &CountryService{db: db}
}

type CountryResponse struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Code3             string `json:"code3"`
	IsBillingEnabled  bool   `json:"isBillingEnabled"`
	IsShippingEnabled bool   `json:"isShippingEnabled"`
	IsCityEnabled     bool   `json:"isCityEnabled"`
	IsZipCodeEnabled  bool   `json:"isZipCodeEnabled"`
	IsDistrictEnabled bool   `json:"isDistrictEnabled"`
}

type StateOrProvinceResponse struct {
	ID        string    `json:"id"`
	CountryID string    `json:"countryId"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type DistrictResponse struct {
	ID                string    `json:"id"`
	StateOrProvinceID string    `json:"stateOrProvinceId"`
	Name              string    `json:"name"`
	Type              string    `json:"type"`
	Location          string    `json:"location"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func (s *CountryService) FindAll(ctx context.Context) ([]CountryResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, COALESCE(code3, ''), is_billing_enabled, is_shipping_enabled, is_city_enabled, is_zip_code_enabled, is_district_enabled
		FROM identity_countries ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list countries: %w", err)
	}
	defer rows.Close()

	var countries []CountryResponse
	for rows.Next() {
		var c CountryResponse
		if err := rows.Scan(&c.ID, &c.Name, &c.Code3, &c.IsBillingEnabled, &c.IsShippingEnabled, &c.IsCityEnabled, &c.IsZipCodeEnabled, &c.IsDistrictEnabled); err != nil {
			return nil, fmt.Errorf("failed to scan country: %w", err)
		}
		countries = append(countries, c)
	}
	if countries == nil {
		countries = []CountryResponse{}
	}
	return countries, nil
}

func (s *CountryService) FindByID(ctx context.Context, id string) (*CountryResponse, error) {
	var c CountryResponse
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, COALESCE(code3, ''), is_billing_enabled, is_shipping_enabled, is_city_enabled, is_zip_code_enabled, is_district_enabled
		FROM identity_countries WHERE id = $1
	`, id).Scan(&c.ID, &c.Name, &c.Code3, &c.IsBillingEnabled, &c.IsShippingEnabled, &c.IsCityEnabled, &c.IsZipCodeEnabled, &c.IsDistrictEnabled)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("country not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find country: %w", err)
	}
	return &c, nil
}

func (s *CountryService) FindStatesByCountryID(ctx context.Context, countryID string) ([]StateOrProvinceResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, country_id, COALESCE(code, ''), name, COALESCE(type, ''), created_at, updated_at
		FROM identity_state_or_provinces WHERE country_id = $1 ORDER BY name ASC
	`, countryID)
	if err != nil {
		return nil, fmt.Errorf("failed to list states: %w", err)
	}
	defer rows.Close()

	var states []StateOrProvinceResponse
	for rows.Next() {
		var s StateOrProvinceResponse
		if err := rows.Scan(&s.ID, &s.CountryID, &s.Code, &s.Name, &s.Type, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan state: %w", err)
		}
		states = append(states, s)
	}
	if states == nil {
		states = []StateOrProvinceResponse{}
	}
	return states, nil
}

func (s *CountryService) FindStateByID(ctx context.Context, id string) (*StateOrProvinceResponse, error) {
	var state StateOrProvinceResponse
	err := s.db.QueryRowContext(ctx, `
		SELECT id, country_id, COALESCE(code, ''), name, COALESCE(type, ''), created_at, updated_at
		FROM identity_state_or_provinces WHERE id = $1
	`, id).Scan(&state.ID, &state.CountryID, &state.Code, &state.Name, &state.Type, &state.CreatedAt, &state.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("state not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find state: %w", err)
	}
	return &state, nil
}

func (s *CountryService) FindDistrictsByStateID(ctx context.Context, stateID string) ([]DistrictResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, state_or_province_id, name, COALESCE(type, ''), COALESCE(location, ''), created_at, updated_at
		FROM identity_districts WHERE state_or_province_id = $1 ORDER BY name ASC
	`, stateID)
	if err != nil {
		return nil, fmt.Errorf("failed to list districts: %w", err)
	}
	defer rows.Close()

	var districts []DistrictResponse
	for rows.Next() {
		var d DistrictResponse
		if err := rows.Scan(&d.ID, &d.StateOrProvinceID, &d.Name, &d.Type, &d.Location, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan district: %w", err)
		}
		districts = append(districts, d)
	}
	if districts == nil {
		districts = []DistrictResponse{}
	}
	return districts, nil
}

func (s *CountryService) FindDistrictByID(ctx context.Context, id string) (*DistrictResponse, error) {
	var d DistrictResponse
	err := s.db.QueryRowContext(ctx, `
		SELECT id, state_or_province_id, name, COALESCE(type, ''), COALESCE(location, ''), created_at, updated_at
		FROM identity_districts WHERE id = $1
	`, id).Scan(&d.ID, &d.StateOrProvinceID, &d.Name, &d.Type, &d.Location, &d.CreatedAt, &d.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("district not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find district: %w", err)
	}
	return &d, nil
}
