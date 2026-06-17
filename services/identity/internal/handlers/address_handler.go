package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/simplcommerce-go/pkg/response"
	"github.com/simplcommerce-go/services/identity/internal/models"
	"github.com/simplcommerce-go/services/identity/internal/services"
)

type AddressHandler struct {
	addressService services.AddressService
}

func NewAddressHandler(addressService services.AddressService) *AddressHandler {
	return &AddressHandler{addressService: addressService}
}

func (h *AddressHandler) Index(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(fiber.StatusUnauthorized, "unauthorized"))
	}

	addresses, err := h.addressService.FindByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(addresses))
}

func (h *AddressHandler) Create(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(fiber.StatusUnauthorized, "unauthorized"))
	}

	var req struct {
		ContactName       string `json:"contactName"`
		Phone             string `json:"phone"`
		AddressLine1      string `json:"addressLine1"`
		AddressLine2      string `json:"addressLine2"`
		City              string `json:"city"`
		ZipCode           string `json:"zipCode"`
		StateOrProvinceID uint   `json:"stateOrProvinceId"`
		CountryID         string `json:"countryId"`
		DistrictID        *uint  `json:"districtId"`
	}

	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "invalid request body"))
	}

	address := &models.Address{
		ContactName:       req.ContactName,
		Phone:             req.Phone,
		AddressLine1:      req.AddressLine1,
		AddressLine2:      req.AddressLine2,
		City:              req.City,
		ZipCode:           req.ZipCode,
		StateOrProvinceID: req.StateOrProvinceID,
		CountryID:         req.CountryID,
		DistrictID:        req.DistrictID,
	}

	result, err := h.addressService.Create(c.Context(), userID, address)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(result))
}

func (h *AddressHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "invalid address id"))
	}

	var req struct {
		ContactName       string `json:"contactName"`
		Phone             string `json:"phone"`
		AddressLine1      string `json:"addressLine1"`
		AddressLine2      string `json:"addressLine2"`
		City              string `json:"city"`
		ZipCode           string `json:"zipCode"`
		StateOrProvinceID uint   `json:"stateOrProvinceId"`
		CountryID         string `json:"countryId"`
		DistrictID        *uint  `json:"districtId"`
	}

	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "invalid request body"))
	}

	address := &models.Address{
		ContactName:       req.ContactName,
		Phone:             req.Phone,
		AddressLine1:      req.AddressLine1,
		AddressLine2:      req.AddressLine2,
		City:              req.City,
		ZipCode:           req.ZipCode,
		StateOrProvinceID: req.StateOrProvinceID,
		CountryID:         req.CountryID,
		DistrictID:        req.DistrictID,
	}

	result, err := h.addressService.Update(c.Context(), uint(id), address)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(result))
}

func (h *AddressHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(fiber.StatusBadRequest, "invalid address id"))
	}

	if err := h.addressService.Delete(c.Context(), uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(fiber.StatusInternalServerError, err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil))
}
