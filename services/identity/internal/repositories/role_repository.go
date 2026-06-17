package repositories

import (
	"context"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/identity/internal/models"
)

type RoleRepository interface {
	Create(ctx context.Context, role *models.Role) error
	FindByName(ctx context.Context, name string) (*models.Role, error)
	FindAll(ctx context.Context) ([]models.Role, error)
	AssignToUser(ctx context.Context, userID, roleID uint) error
	RemoveFromUser(ctx context.Context, userID, roleID uint) error
}

type GormRoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &GormRoleRepository{db: db}
}

func (r *GormRoleRepository) Create(ctx context.Context, role *models.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *GormRoleRepository) FindByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *GormRoleRepository) FindAll(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.WithContext(ctx).Find(&roles).Error
	return roles, err
}

func (r *GormRoleRepository) AssignToUser(ctx context.Context, userID, roleID uint) error {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}
	var role models.Role
	if err := r.db.WithContext(ctx).First(&role, roleID).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&user).Association("Roles").Append(&role)
}

func (r *GormRoleRepository) RemoveFromUser(ctx context.Context, userID, roleID uint) error {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}
	var role models.Role
	if err := r.db.WithContext(ctx).First(&role, roleID).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&user).Association("Roles").Delete(&role)
}
