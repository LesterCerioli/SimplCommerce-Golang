package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/inventory/internal/repositories"
	"github.com/simplcommerce-go/services/inventory/internal/services"
)

type StockHandler struct {
	inventoryService services.InventoryService
	stockRepo        repositories.StockRepository
	stockHistoryRepo repositories.StockHistoryRepository
}

func NewStockHandler(
	inventoryService services.InventoryService,
	stockRepo repositories.StockRepository,
	stockHistoryRepo repositories.StockHistoryRepository,
) *StockHandler {
	return &StockHandler{
		inventoryService: inventoryService,
		stockRepo:        stockRepo,
		stockHistoryRepo: stockHistoryRepo,
	}
}

type updateStockRequest struct {
	WarehouseID uint   `json:"warehouseId"`
	Quantity    int    `json:"quantity"`
	Note        string `json:"note"`
}

func (h *StockHandler) GetStock(c fiber.Ctx) error {
	productIDStr := c.Params("productId")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid product id")
	}

	stock, err := h.inventoryService.GetStock(c.Context(), uint(productID))
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	if stock == nil {
		return response.Error(fiber.StatusNotFound, "stock not found")
	}

	return c.JSON(response.Success(stock))
}

func (h *StockHandler) UpdateStock(c fiber.Ctx) error {
	productIDStr := c.Params("productId")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid product id")
	}

	var req updateStockRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	userID, _ := c.Locals("userID").(uint)

	stock, err := h.inventoryService.UpdateStock(c.Context(), uint(productID), req.WarehouseID, req.Quantity, userID, req.Note)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(stock))
}

func (h *StockHandler) GetStockHistory(c fiber.Ctx) error {
	productIDStr := c.Params("productId")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid product id")
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))

	histories, total, err := h.stockHistoryRepo.PaginateByProductID(c.Context(), uint(productID), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(histories, page, pageSize, total))
}
