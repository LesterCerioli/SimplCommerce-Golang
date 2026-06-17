package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/cms/internal/models"
)

type MenuRepository interface {
	Create(ctx context.Context, menu *models.Menu) error
	Update(ctx context.Context, menu *models.Menu) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*models.Menu, error)
	FindAll(ctx context.Context) ([]models.Menu, error)
}

type GormMenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &GormMenuRepository{db: db}
}

func (r *GormMenuRepository) Create(ctx context.Context, menu *models.Menu) error {
	return r.db.WithContext(ctx).Create(menu).Error
}

func (r *GormMenuRepository) Update(ctx context.Context, menu *models.Menu) error {
	return r.db.WithContext(ctx).Save(menu).Error
}

func (r *GormMenuRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Menu{}, id).Error
}

func (r *GormMenuRepository) FindByID(ctx context.Context, id uint) (*models.Menu, error) {
	var menu models.Menu
	err := r.db.WithContext(ctx).Preload("Items").First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *GormMenuRepository) FindAll(ctx context.Context) ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.WithContext(ctx).Preload("Items").Find(&menus).Error
	return menus, err
}
