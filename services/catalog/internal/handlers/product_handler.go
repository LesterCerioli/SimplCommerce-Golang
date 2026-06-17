package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"

	"github.com/simplcommerce-go/pkg/response"
	catalogmodels "github.com/simplcommerce-go/services/catalog/internal/models"
	"github.com/simplcommerce-go/services/catalog/internal/services"
)

type ProductHandler struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

type createProductRequest struct {
	Name                 string             `json:"name" validate:"required"`
	Slug                 string             `json:"slug,omitempty"`
	ShortDescription     string             `json:"shortDescription,omitempty"`
	Description          string             `json:"description,omitempty"`
	Specification        string             `json:"specification,omitempty"`
	Price                float64            `json:"price" validate:"required"`
	OldPrice             *float64           `json:"oldPrice,omitempty"`
	SpecialPrice         *float64           `json:"specialPrice,omitempty"`
	SpecialPriceStart    *string            `json:"specialPriceStart,omitempty"`
	SpecialPriceEnd      *string            `json:"specialPriceEnd,omitempty"`
	HasOptions           bool               `json:"hasOptions,omitempty"`
	IsVisibleIndividually bool              `json:"isVisibleIndividually,omitempty"`
	IsFeatured           bool               `json:"isFeatured,omitempty"`
	IsCallForPricing     bool               `json:"isCallForPricing,omitempty"`
	IsAllowToOrder       bool               `json:"isAllowToOrder,omitempty"`
	StockTrackingIsEnabled bool             `json:"stockTrackingIsEnabled,omitempty"`
	StockQuantity        int                `json:"stockQuantity,omitempty"`
	SKU                  string             `json:"sku,omitempty"`
	GTIN                 string             `json:"gtin,omitempty"`
	DisplayOrder         int                `json:"displayOrder,omitempty"`
	BrandID              *uint              `json:"brandId,omitempty"`
	TaxClassID           *uint              `json:"taxClassId,omitempty"`
	ThumbnailImageID     *uint              `json:"thumbnailImageId,omitempty"`
	IsPublished          bool               `json:"isPublished,omitempty"`
	CategoryIDs          []uint             `json:"categoryIds,omitempty"`
}

type updateProductRequest struct {
	Name                 string     `json:"name"`
	Slug                 string     `json:"slug,omitempty"`
	ShortDescription     string     `json:"shortDescription,omitempty"`
	Description          string     `json:"description,omitempty"`
	Specification        string     `json:"specification,omitempty"`
	Price                float64    `json:"price"`
	OldPrice             *float64   `json:"oldPrice,omitempty"`
	SpecialPrice         *float64   `json:"specialPrice,omitempty"`
	SpecialPriceStart    *string    `json:"specialPriceStart,omitempty"`
	SpecialPriceEnd      *string    `json:"specialPriceEnd,omitempty"`
	HasOptions           bool       `json:"hasOptions,omitempty"`
	IsVisibleIndividually bool      `json:"isVisibleIndividually,omitempty"`
	IsFeatured           bool       `json:"isFeatured,omitempty"`
	IsCallForPricing     bool       `json:"isCallForPricing,omitempty"`
	IsAllowToOrder       bool       `json:"isAllowToOrder,omitempty"`
	StockTrackingIsEnabled bool     `json:"stockTrackingIsEnabled,omitempty"`
	StockQuantity        int        `json:"stockQuantity,omitempty"`
	SKU                  string     `json:"sku,omitempty"`
	GTIN                 string     `json:"gtin,omitempty"`
	DisplayOrder         int        `json:"displayOrder,omitempty"`
	BrandID              *uint      `json:"brandId,omitempty"`
	TaxClassID           *uint      `json:"taxClassId,omitempty"`
	ThumbnailImageID     *uint      `json:"thumbnailImageId,omitempty"`
	IsPublished          bool       `json:"isPublished,omitempty"`
	CategoryIDs          []uint     `json:"categoryIds,omitempty"`
}

func (h *ProductHandler) Index(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))

	products, total, err := h.service.GetProducts(c.Context(), page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(products, page, pageSize, total))
}

func (h *ProductHandler) Featured(c fiber.Ctx) error {
	count, _ := strconv.Atoi(c.Query("count", "10"))

	products, err := h.service.GetFeaturedProducts(c.Context(), count)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(products))
}

func (h *ProductHandler) Show(c fiber.Ctx) error {
	slug := c.Params("slug")

	product, err := h.service.GetProductBySlug(c.Context(), slug)
	if err != nil {
		return response.Error(fiber.StatusNotFound, "product not found")
	}

	return c.JSON(response.Success(product))
}

func (h *ProductHandler) Create(c fiber.Ctx) error {
	var req createProductRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	product := &catalogmodels.Product{
		Name:                   req.Name,
		Slug:                   req.Slug,
		ShortDescription:       req.ShortDescription,
		Description:            req.Description,
		Specification:          req.Specification,
		Price:                  req.Price,
		OldPrice:               req.OldPrice,
		SpecialPrice:           req.SpecialPrice,
		HasOptions:             req.HasOptions,
		IsVisibleIndividually:  req.IsVisibleIndividually,
		IsFeatured:             req.IsFeatured,
		IsCallForPricing:       req.IsCallForPricing,
		IsAllowToOrder:         req.IsAllowToOrder,
		StockTrackingIsEnabled: req.StockTrackingIsEnabled,
		StockQuantity:          req.StockQuantity,
		SKU:                    req.SKU,
		GTIN:                   req.GTIN,
		DisplayOrder:           req.DisplayOrder,
		BrandID:                req.BrandID,
		TaxClassID:             req.TaxClassID,
		ThumbnailImageID:       req.ThumbnailImageID,
		IsPublished:            req.IsPublished,
	}

	if err := h.service.CreateProduct(c.Context(), product); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(product))
}

func (h *ProductHandler) Update(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid product id")
	}

	var req updateProductRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	product, err := h.service.GetProductByID(c.Context(), uint(id))
	if err != nil {
		return response.Error(fiber.StatusNotFound, "product not found")
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Slug != "" {
		product.Slug = req.Slug
	}
	product.ShortDescription = req.ShortDescription
	product.Description = req.Description
	product.Specification = req.Specification
	product.Price = req.Price
	product.OldPrice = req.OldPrice
	product.SpecialPrice = req.SpecialPrice
	product.HasOptions = req.HasOptions
	product.IsVisibleIndividually = req.IsVisibleIndividually
	product.IsFeatured = req.IsFeatured
	product.IsCallForPricing = req.IsCallForPricing
	product.IsAllowToOrder = req.IsAllowToOrder
	product.StockTrackingIsEnabled = req.StockTrackingIsEnabled
	product.StockQuantity = req.StockQuantity
	product.SKU = req.SKU
	product.GTIN = req.GTIN
	product.DisplayOrder = req.DisplayOrder
	product.BrandID = req.BrandID
	product.TaxClassID = req.TaxClassID
	product.ThumbnailImageID = req.ThumbnailImageID
	product.IsPublished = req.IsPublished

	if err := h.service.UpdateProduct(c.Context(), product); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(product))
}

func (h *ProductHandler) Delete(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid product id")
	}

	if err := h.service.DeleteProduct(c.Context(), uint(id)); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(nil))
}

func (h *ProductHandler) Search(c fiber.Ctx) error {
	query := c.Query("q")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "20"))

	var minPrice, maxPrice *float64
	if v := c.Query("minPrice"); v != "" {
		if p, err := strconv.ParseFloat(v, 64); err == nil {
			minPrice = &p
		}
	}
	if v := c.Query("maxPrice"); v != "" {
		if p, err := strconv.ParseFloat(v, 64); err == nil {
			maxPrice = &p
		}
	}

	var categoryID *uint
	if v := c.Query("categoryId"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			uid := uint(id)
			categoryID = &uid
		}
	}

	products, total, err := h.service.SearchProducts(c.Context(), query, minPrice, maxPrice, categoryID, page, pageSize)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(products, page, pageSize, total))
}
