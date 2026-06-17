package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/cms/internal/models"
)

type MenuItemRepository interface {
	Create(ctx context.Context, item *models.MenuItem) error
	Update(ctx context.Context, item *models.MenuItem) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.MenuItem, error)
	FindByMenuID(ctx context.Context, menuID uint) ([]models.MenuItem, error)
}

type GormMenuItemRepository struct {
	db *gorm.DB
}

func NewMenuItemRepository(db *gorm.DB) MenuItemRepository {
	return &GormMenuItemRepository{db: db}
}

func (r *GormMenuItemRepository) Create(ctx context.Context, item *models.MenuItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *GormMenuItemRepository) Update(ctx context.Context, item *models.MenuItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *GormMenuItemRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.MenuItem{}, id).Error
}

func (r *GormMenuItemRepository) FindByID(ctx context.Context, id uint) (*models.MenuItem, error) {
	var item models.MenuItem
	err := r.db.WithContext(ctx).First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *GormMenuItemRepository) FindByMenuID(ctx context.Context, menuID uint) ([]models.MenuItem, error) {
	var items []models.MenuItem
	err := r.db.WithContext(ctx).Where("menu_id = ?", menuID).Order("display_order asc").Find(&items).Error
	return items, err
}
