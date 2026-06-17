package services

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/simplcommerce-go/services/inventory/internal/models"
	"github.com/simplcommerce-go/services/inventory/internal/repositories"
)

type StockResult struct {
	ProductID   uint `json:"productId"`
	TotalStock  int  `json:"totalStock"`
	Available   int  `json:"available"`
	Warehouses  []WarehouseStock `json:"warehouses"`
}

type WarehouseStock struct {
	WarehouseID uint   `json:"warehouseId"`
	WarehouseName string `json:"warehouseName"`
	Quantity    int    `json:"quantity"`
	Reserved    int    `json:"reserved"`
	Available   int    `json:"available"`
}

type InventoryService interface {
	GetStock(ctx context.Context, productID uint) (*StockResult, error)
	UpdateStock(ctx context.Context, productID, warehouseID uint, quantity int, createdByID uint, note string) (*models.Stock, error)
	AdjustStock(ctx context.Context, productID, warehouseID uint, adjustment int, createdByID uint, note string) (*models.Stock, error)
	ReserveStock(ctx context.Context, productID, warehouseID uint, quantity int) (*models.Stock, error)
}

type inventoryService struct {
	stockRepo       repositories.StockRepository
	stockHistoryRepo repositories.StockHistoryRepository
	db              *gorm.DB
}

func NewInventoryService(
	stockRepo repositories.StockRepository,
	stockHistoryRepo repositories.StockHistoryRepository,
	db *gorm.DB,
) InventoryService {
	return &inventoryService{
		stockRepo:       stockRepo,
		stockHistoryRepo: stockHistoryRepo,
		db:              db,
	}
}

func (s *inventoryService) GetStock(ctx context.Context, productID uint) (*StockResult, error) {
	stocks, err := s.stockRepo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	result := &StockResult{
		ProductID:  productID,
		TotalStock: 0,
		Available:  0,
		Warehouses: make([]WarehouseStock, 0),
	}

	for _, stock := range stocks {
		available := stock.Quantity - stock.ReservedQuantity
		warehouseName := ""
		if stock.Warehouse.ID != 0 {
			warehouseName = stock.Warehouse.Name
		}

		result.TotalStock += stock.Quantity
		result.Available += available
		result.Warehouses = append(result.Warehouses, WarehouseStock{
			WarehouseID:   stock.WarehouseID,
			WarehouseName: warehouseName,
			Quantity:      stock.Quantity,
			Reserved:      stock.ReservedQuantity,
			Available:     available,
		})
	}

	return result, nil
}

func (s *inventoryService) UpdateStock(ctx context.Context, productID, warehouseID uint, quantity int, createdByID uint, note string) (*models.Stock, error) {
	var stock *models.Stock

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repo := repositories.NewStockRepository(tx)
		histRepo := repositories.NewStockHistoryRepository(tx)

		existing, err := repo.FindByProductAndWarehouse(ctx, productID, warehouseID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				stock = &models.Stock{
					ProductID:        productID,
					WarehouseID:      warehouseID,
					Quantity:         quantity,
					ReservedQuantity: 0,
				}
				if err := repo.Create(ctx, stock); err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			adjustment := quantity - existing.Quantity
			existing.Quantity = quantity
			if err := repo.Update(ctx, existing); err != nil {
				return err
			}
			stock = existing

			if adjustment != 0 {
				history := &models.StockHistory{
					ProductID:        productID,
					WarehouseID:      warehouseID,
					CreatedByID:      createdByID,
					AdjustedQuantity: adjustment,
					Note:             note,
				}
				if err := histRepo.Create(ctx, history); err != nil {
					return err
				}
			}
		}

		if stock != nil && stock.ID == 0 {
			history := &models.StockHistory{
				ProductID:        productID,
				WarehouseID:      warehouseID,
				CreatedByID:      createdByID,
				AdjustedQuantity: quantity,
				Note:             note,
			}
			if err := histRepo.Create(ctx, history); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.stockRepo.FindByID(ctx, stock.ID)
}

func (s *inventoryService) AdjustStock(ctx context.Context, productID, warehouseID uint, adjustment int, createdByID uint, note string) (*models.Stock, error) {
	var stock *models.Stock

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repo := repositories.NewStockRepository(tx)
		histRepo := repositories.NewStockHistoryRepository(tx)

		existing, err := repo.FindByProductAndWarehouse(ctx, productID, warehouseID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				stock = &models.Stock{
					ProductID:        productID,
					WarehouseID:      warehouseID,
					Quantity:         adjustment,
					ReservedQuantity: 0,
				}
				if err := repo.Create(ctx, stock); err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			existing.Quantity += adjustment
			if err := repo.Update(ctx, existing); err != nil {
				return err
			}
			stock = existing
		}

		history := &models.StockHistory{
			ProductID:        productID,
			WarehouseID:      warehouseID,
			CreatedByID:      createdByID,
			AdjustedQuantity: adjustment,
			Note:             note,
		}
		return histRepo.Create(ctx, history)
	})

	if err != nil {
		return nil, err
	}

	return s.stockRepo.FindByID(ctx, stock.ID)
}

func (s *inventoryService) ReserveStock(ctx context.Context, productID, warehouseID uint, quantity int) (*models.Stock, error) {
	var stock *models.Stock

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repo := repositories.NewStockRepository(tx)

		existing, err := repo.FindByProductAndWarehouse(ctx, productID, warehouseID)
		if err != nil {
			return err
		}

		available := existing.Quantity - existing.ReservedQuantity
		if available < quantity {
			return errors.New("insufficient stock")
		}

		existing.ReservedQuantity += quantity
		if err := repo.Update(ctx, existing); err != nil {
			return err
		}
		stock = existing
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.stockRepo.FindByID(ctx, stock.ID)
}
