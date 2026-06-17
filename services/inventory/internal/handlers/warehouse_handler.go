package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/inventory/internal/models"
	"github.com/simplcommerce-go/services/inventory/internal/services"
)

type WarehouseHandler struct {
	warehouseService services.WarehouseService
}

func NewWarehouseHandler(warehouseService services.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{warehouseService: warehouseService}
}

type createWarehouseRequest struct {
	Name      string `json:"name"`
	AddressID uint   `json:"addressId"`
	VendorID  *uint  `json:"vendorId,omitempty"`
}

type updateWarehouseRequest struct {
	Name      string `json:"name"`
	AddressID uint   `json:"addressId"`
	VendorID  *uint  `json:"vendorId,omitempty"`
}

func (h *WarehouseHandler) List(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))

	warehouses, total, err := h.warehouseService.Paginate(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(warehouses, page, pageSize, total))
}

func (h *WarehouseHandler) Create(c fiber.Ctx) error {
	var req createWarehouseRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	warehouse := &models.Warehouse{
		Name:      req.Name,
		AddressID: req.AddressID,
		VendorID:  req.VendorID,
	}

	if err := h.warehouseService.Create(c.Context(), warehouse); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(warehouse))
}

func (h *WarehouseHandler) Update(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid warehouse id")
	}

	var req updateWarehouseRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	warehouse, err := h.warehouseService.FindByID(c.Context(), uint(id))
	if err != nil {
		return response.Error(fiber.StatusNotFound, "warehouse not found")
	}

	warehouse.Name = req.Name
	warehouse.AddressID = req.AddressID
	warehouse.VendorID = req.VendorID

	if err := h.warehouseService.Update(c.Context(), warehouse); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(warehouse))
}

func (h *WarehouseHandler) Delete(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid warehouse id")
	}

	if err := h.warehouseService.Delete(c.Context(), uint(id)); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(nil))
}
