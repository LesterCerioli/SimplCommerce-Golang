package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type AddressService struct {
	db *sql.DB
}

func NewAddressService(db *sql.DB) *AddressService {
	return &AddressService{db: db}
}

type AddressResponse struct {
	ID                string    `json:"id"`
	ContactName       string    `json:"contactName"`
	Phone             string    `json:"phone"`
	AddressLine1      string    `json:"addressLine1"`
	AddressLine2      string    `json:"addressLine2"`
	City              string    `json:"city"`
	ZipCode           string    `json:"zipCode"`
	DistrictID        *string   `json:"districtId"`
	StateOrProvinceID uint      `json:"stateOrProvinceId"`
	CountryID         string    `json:"countryId"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type CreateAddressRequest struct {
	ContactName       string `json:"contactName"`
	Phone             string `json:"phone"`
	AddressLine1      string `json:"addressLine1"`
	AddressLine2      string `json:"addressLine2"`
	City              string `json:"city"`
	ZipCode           string `json:"zipCode"`
	DistrictID        *uint  `json:"districtId"`
	StateOrProvinceID uint   `json:"stateOrProvinceId"`
	CountryID         string `json:"countryId"`
}

type UpdateAddressRequest struct {
	ContactName       string `json:"contactName"`
	Phone             string `json:"phone"`
	AddressLine1      string `json:"addressLine1"`
	AddressLine2      string `json:"addressLine2"`
	City              string `json:"city"`
	ZipCode           string `json:"zipCode"`
	DistrictID        *uint  `json:"districtId"`
	StateOrProvinceID uint   `json:"stateOrProvinceId"`
	CountryID         string `json:"countryId"`
}

func (s *AddressService) Create(ctx context.Context, userID string, req CreateAddressRequest) (*AddressResponse, error) {
	var addr AddressResponse
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO identity_addresses (contact_name, phone, address_line1, address_line2, city, zip_code, district_id, state_or_province_id, country_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, contact_name, phone, address_line1, address_line2, city, zip_code, district_id, state_or_province_id, country_id, created_at, updated_at
	`, req.ContactName, req.Phone, req.AddressLine1, req.AddressLine2, req.City, req.ZipCode, req.DistrictID, req.StateOrProvinceID, req.CountryID).Scan(
		&addr.ID, &addr.ContactName, &addr.Phone, &addr.AddressLine1, &addr.AddressLine2, &addr.City, &addr.ZipCode, &addr.DistrictID, &addr.StateOrProvinceID, &addr.CountryID, &addr.CreatedAt, &addr.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO identity_user_addresses (user_id, address_id, address_type, created_at, updated_at)
		VALUES ($1, $2, 0, NOW(), NOW())
	`, userID, addr.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to link address to user: %w", err)
	}

	return &addr, nil
}

func (s *AddressService) Update(ctx context.Context, id string, req UpdateAddressRequest) (*AddressResponse, error) {
	_, err := s.db.ExecContext(ctx, `
		UPDATE identity_addresses SET contact_name = $1, phone = $2, address_line1 = $3, address_line2 = $4,
			city = $5, zip_code = $6, district_id = $7, state_or_province_id = $8, country_id = $9, updated_at = NOW()
		WHERE id = $10
	`, req.ContactName, req.Phone, req.AddressLine1, req.AddressLine2, req.City, req.ZipCode, req.DistrictID, req.StateOrProvinceID, req.CountryID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update address: %w", err)
	}

	return s.FindByID(ctx, id)
}

func (s *AddressService) Delete(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM identity_user_addresses WHERE address_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to remove address links: %w", err)
	}

	result, err := s.db.ExecContext(ctx, `DELETE FROM identity_addresses WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("address not found")
	}
	return nil
}

func (s *AddressService) FindByID(ctx context.Context, id string) (*AddressResponse, error) {
	var a AddressResponse
	var districtID sql.NullString
	err := s.db.QueryRowContext(ctx, `
		SELECT id, contact_name, phone, address_line1, address_line2, city, zip_code, district_id, state_or_province_id, country_id, created_at, updated_at
		FROM identity_addresses WHERE id = $1
	`, id).Scan(
		&a.ID, &a.ContactName, &a.Phone, &a.AddressLine1, &a.AddressLine2, &a.City, &a.ZipCode, &districtID, &a.StateOrProvinceID, &a.CountryID, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("address not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find address: %w", err)
	}
	if districtID.Valid {
		a.DistrictID = &districtID.String
	}
	return &a, nil
}

func (s *AddressService) FindByUserID(ctx context.Context, userID string) ([]AddressResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT a.id, a.contact_name, a.phone, a.address_line1, a.address_line2, a.city, a.zip_code, a.district_id, a.state_or_province_id, a.country_id, a.created_at, a.updated_at
		FROM identity_addresses a
		JOIN identity_user_addresses ua ON ua.address_id = a.id
		WHERE ua.user_id = $1
		ORDER BY a.created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list addresses: %w", err)
	}
	defer rows.Close()

	var addresses []AddressResponse
	for rows.Next() {
		var a AddressResponse
		var districtID sql.NullString
		if err := rows.Scan(&a.ID, &a.ContactName, &a.Phone, &a.AddressLine1, &a.AddressLine2, &a.City, &a.ZipCode, &districtID, &a.StateOrProvinceID, &a.CountryID, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan address: %w", err)
		}
		if districtID.Valid {
			a.DistrictID = &districtID.String
		}
		addresses = append(addresses, a)
	}
	if addresses == nil {
		addresses = []AddressResponse{}
	}
	return addresses, nil
}

func (s *AddressService) IsOwner(ctx context.Context, addressID, userID string) (bool, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM identity_user_addresses WHERE address_id = $1 AND user_id = $2
	`, addressID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check ownership: %w", err)
	}
	return count > 0, nil
}
