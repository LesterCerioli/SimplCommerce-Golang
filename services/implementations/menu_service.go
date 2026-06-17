package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type MenuService struct {
	db *sql.DB
}

func NewMenuService(db *sql.DB) *MenuService {
	return &MenuService{db: db}
}

type MenuResponse struct {
	ID          uint               `json:"id"`
	Name        string             `json:"name"`
	IsPublished bool               `json:"isPublished"`
	IsSystem    bool               `json:"isSystem"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	Items       []MenuItemResponse `json:"items,omitempty"`
}

type MenuItemResponse struct {
	ID           uint      `json:"id"`
	ParentID     *uint     `json:"parentId,omitempty"`
	MenuID       uint      `json:"menuId"`
	EntityID     *uint     `json:"entityId,omitempty"`
	CustomLink   string    `json:"customLink"`
	Name         string    `json:"name"`
	DisplayOrder int       `json:"displayOrder"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (s *MenuService) List(ctx context.Context) ([]MenuResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, is_published, is_system, created_at, updated_at
		FROM cms_menus ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list menus: %w", err)
	}
	defer rows.Close()

	var menus []MenuResponse
	for rows.Next() {
		var m MenuResponse
		if err := rows.Scan(&m.ID, &m.Name, &m.IsPublished, &m.IsSystem, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		menus = append(menus, m)
	}
	if menus == nil {
		menus = []MenuResponse{}
	}
	return menus, nil
}

func (s *MenuService) GetByID(ctx context.Context, id uint) (*MenuResponse, error) {
	var m MenuResponse
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, is_published, is_system, created_at, updated_at
		FROM cms_menus WHERE id = $1
	`, id).Scan(&m.ID, &m.Name, &m.IsPublished, &m.IsSystem, &m.CreatedAt, &m.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("menu not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get menu: %w", err)
	}

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, parent_id, menu_id, entity_id, COALESCE(custom_link,''), name, display_order, created_at, updated_at
		FROM cms_menu_items WHERE menu_id = $1 ORDER BY display_order ASC
	`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item MenuItemResponse
		var parentID, entityID sql.NullInt64
		if err := rows.Scan(&item.ID, &parentID, &item.MenuID, &entityID, &item.CustomLink, &item.Name, &item.DisplayOrder, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		if parentID.Valid {
			v := uint(parentID.Int64)
			item.ParentID = &v
		}
		if entityID.Valid {
			v := uint(entityID.Int64)
			item.EntityID = &v
		}
		m.Items = append(m.Items, item)
	}
	if m.Items == nil {
		m.Items = []MenuItemResponse{}
	}
	return &m, nil
}

func (s *MenuService) Create(ctx context.Context, name string, isPublished bool) (*MenuResponse, error) {
	var m MenuResponse
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO cms_menus (name, is_published, created_at, updated_at)
		VALUES ($1,$2,NOW(),NOW()) RETURNING id, name, is_published, is_system, created_at, updated_at
	`, name, isPublished).Scan(&m.ID, &m.Name, &m.IsPublished, &m.IsSystem, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create menu: %w", err)
	}
	m.Items = []MenuItemResponse{}
	return &m, nil
}

func (s *MenuService) Update(ctx context.Context, id uint, name string, isPublished bool) (*MenuResponse, error) {
	_, err := s.db.ExecContext(ctx, `
		UPDATE cms_menus SET name=$1, is_published=$2, updated_at=NOW() WHERE id=$3
	`, name, isPublished, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update menu: %w", err)
	}
	return s.GetByID(ctx, id)
}

func (s *MenuService) Delete(ctx context.Context, id uint) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `DELETE FROM cms_menu_items WHERE menu_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete menu items: %w", err)
	}

	result, err := tx.ExecContext(ctx, `DELETE FROM cms_menus WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete menu: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("menu not found")
	}

	return tx.Commit()
}

func (s *MenuService) AddItem(ctx context.Context, menuID uint, parentID, entityID *uint, customLink, name string, displayOrder int) (*MenuItemResponse, error) {
	_, err := s.db.ExecContext(ctx, `SELECT 1 FROM cms_menus WHERE id = $1`, menuID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("menu not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find menu: %w", err)
	}

	var item MenuItemResponse
	var pID, eID sql.NullInt64
	err = s.db.QueryRowContext(ctx, `
		INSERT INTO cms_menu_items (parent_id, menu_id, entity_id, custom_link, name, display_order, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,NOW(),NOW())
		RETURNING id, parent_id, menu_id, entity_id, custom_link, name, display_order, created_at, updated_at
	`, parentID, menuID, entityID, customLink, name, displayOrder).Scan(
		&item.ID, &pID, &item.MenuID, &eID, &item.CustomLink, &item.Name, &item.DisplayOrder, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to add menu item: %w", err)
	}
	if pID.Valid {
		v := uint(pID.Int64)
		item.ParentID = &v
	}
	if eID.Valid {
		v := uint(eID.Int64)
		item.EntityID = &v
	}
	return &item, nil
}

func (s *MenuService) UpdateItem(ctx context.Context, id uint, parentID, entityID *uint, customLink, name string, displayOrder int) (*MenuItemResponse, error) {
	var item MenuItemResponse
	var pID, eID sql.NullInt64
	err := s.db.QueryRowContext(ctx, `
		UPDATE cms_menu_items SET parent_id=$1, entity_id=$2, custom_link=$3, name=$4, display_order=$5, updated_at=NOW()
		WHERE id=$6 RETURNING id, parent_id, menu_id, entity_id, custom_link, name, display_order, created_at, updated_at
	`, parentID, entityID, customLink, name, displayOrder, id).Scan(
		&item.ID, &pID, &item.MenuID, &eID, &item.CustomLink, &item.Name, &item.DisplayOrder, &item.CreatedAt, &item.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("menu item not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update menu item: %w", err)
	}
	if pID.Valid {
		v := uint(pID.Int64)
		item.ParentID = &v
	}
	if eID.Valid {
		v := uint(eID.Int64)
		item.EntityID = &v
	}
	return &item, nil
}

func (s *MenuService) DeleteItem(ctx context.Context, id uint) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM cms_menu_items WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete menu item: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("menu item not found")
	}
	return nil
}
