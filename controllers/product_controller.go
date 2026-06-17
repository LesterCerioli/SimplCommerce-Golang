package controllers

import (
	"github.com/gofiber/fiber/v3"
	"simplcommerce/internal/response"
	"simplcommerce/services/implementations"
)

type ProductController struct {
	productService  *implementations.ProductService
	categoryService *implementations.CategoryService
	brandService    *implementations.BrandService
}

func NewProductController(
	productService *implementations.ProductService,
	categoryService *implementations.CategoryService,
	brandService *implementations.BrandService,
) *ProductController {
	return &ProductController{
		productService:  productService,
		categoryService: categoryService,
		brandService:    brandService,
	}
}

func (h *ProductController) Index(c fiber.Ctx) error {
	page := parseInt(c.Query("page", "1"), 1)
	pageSize := parseInt(c.Query("pageSize", "20"), 20)
	search := c.Query("search", "")

	var categoryID *string
	if v := c.Query("categoryId"); v != "" {
		categoryID = &v
	}

	var brandID *string
	if v := c.Query("brandId"); v != "" {
		brandID = &v
	}

	var minPrice, maxPrice *float64
	if v := c.Query("minPrice"); v != "" {
		if p, err := parseFloat(v); err == nil {
			minPrice = &p
		}
	}
	if v := c.Query("maxPrice"); v != "" {
		if p, err := parseFloat(v); err == nil {
			maxPrice = &p
		}
	}

	products, total, err := h.productService.GetProducts(c.Context(), page, pageSize, search, categoryID, brandID, minPrice, maxPrice)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(products, page, pageSize, total))
}

func (h *ProductController) Featured(c fiber.Ctx) error {
	count := parseInt(c.Query("count", "10"), 10)

	products, err := h.productService.GetFeaturedProducts(c.Context(), count)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(products))
}

func (h *ProductController) Show(c fiber.Ctx) error {
	slug := c.Params("slug")

	product, err := h.productService.GetProductBySlug(c.Context(), slug)
	if err != nil {
		return response.Error(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(response.Success(product))
}

func (h *ProductController) Create(c fiber.Ctx) error {
	var req implementations.CreateProductRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	product, err := h.productService.CreateProduct(c.Context(), &req)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.Created(product))
}

func (h *ProductController) Update(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(fiber.StatusBadRequest, "invalid product id")
	}

	var req implementations.UpdateProductRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(fiber.StatusBadRequest, "invalid request body")
	}

	product, err := h.productService.UpdateProduct(c.Context(), id, &req)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(product))
}

func (h *ProductController) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Error(fiber.StatusBadRequest, "invalid product id")
	}

	if err := h.productService.DeleteProduct(c.Context(), id); err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Success(nil))
}

func (h *ProductController) Search(c fiber.Ctx) error {
	query := c.Query("q")
	page := parseInt(c.Query("page", "1"), 1)
	pageSize := parseInt(c.Query("pageSize", "20"), 20)

	var minPrice, maxPrice *float64
	if v := c.Query("minPrice"); v != "" {
		if p, err := parseFloat(v); err == nil {
			minPrice = &p
		}
	}
	if v := c.Query("maxPrice"); v != "" {
		if p, err := parseFloat(v); err == nil {
			maxPrice = &p
		}
	}

	var categoryID *string
	if v := c.Query("categoryId"); v != "" {
		categoryID = &v
	}

	products, total, err := h.productService.SearchProducts(c.Context(), query, page, pageSize, categoryID, minPrice, maxPrice)
	if err != nil {
		return response.Error(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(response.Paginated(products, page, pageSize, total))
}
