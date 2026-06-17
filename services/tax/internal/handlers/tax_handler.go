package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/tax/internal/services"
)

type TaxClassHandler struct {
	taxService services.TaxService
}

func NewTaxClassHandler(taxService services.TaxService) *TaxClassHandler {
	return &TaxClassHandler{taxService: taxService}
}

func (h *TaxClassHandler) GetAll(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	classes, total, err := h.taxService.GetAllTaxClasses(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(classes, page, pageSize, total))
}

func (h *TaxClassHandler) Create(c fiber.Ctx) error {
	var input struct {
		Name string `json:"name"`
	}

	if err := c.Bind().Body(&input); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	class, err := h.taxService.CreateTaxClass(c.Context(), input.Name)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(class))
}

func (h *TaxClassHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid tax class id")
	}

	var input struct {
		Name string `json:"name"`
	}

	if err := c.Bind().Body(&input); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	class, err := h.taxService.UpdateTaxClass(c.Context(), uint(id), input.Name)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(response.Success(class))
}

func (h *TaxClassHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid tax class id")
	}

	if err := h.taxService.DeleteTaxClass(c.Context(), uint(id)); err != nil {
		return response.Error(fiber.StatusNotFound, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

type TaxRateHandler struct {
	taxService services.TaxService
}

func NewTaxRateHandler(taxService services.TaxService) *TaxRateHandler {
	return &TaxRateHandler{taxService: taxService}
}

func (h *TaxRateHandler) GetAll(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	rates, total, err := h.taxService.GetAllTaxRates(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(rates, page, pageSize, total))
}

func (h *TaxRateHandler) Create(c fiber.Ctx) error {
	var input services.CreateTaxRateInput
	if err := c.Bind().Body(&input); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	rate, err := h.taxService.CreateTaxRate(c.Context(), input)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(rate))
}

func (h *TaxRateHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid tax rate id")
	}

	var input services.UpdateTaxRateInput
	if err := c.Bind().Body(&input); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	rate, err := h.taxService.UpdateTaxRate(c.Context(), uint(id), input)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(response.Success(rate))
}

func (h *TaxRateHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid tax rate id")
	}

	if err := h.taxService.DeleteTaxRate(c.Context(), uint(id)); err != nil {
		return response.Error(fiber.StatusNotFound, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
