package implementations

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v3"
)

type ProductService struct {
	db *sql.DB
}

type CategoryBrief struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type AttributeValueResponse struct {
	ID            uint   `json:"id"`
	AttributeID   uint   `json:"attributeId"`
	AttributeName string `json:"attributeName"`
	Value         string `json:"value"`
}

type OptionValueResponse struct {
	ID          uint   `json:"id"`
	OptionID    uint   `json:"optionId"`
	OptionName  string `json:"optionName"`
	Value       string `json:"value"`
	DisplayType string `json:"displayType"`
	SortIndex   int    `json:"sortIndex"`
}

type ProductResponse struct {
	ID                     uint                    `json:"id"`
	Name                   string                  `json:"name"`
	Slug                   string                  `json:"slug"`
	ShortDescription       string                  `json:"shortDescription"`
	Description            string                  `json:"description"`
	Specification          string                  `json:"specification"`
	Price                  float64                 `json:"price"`
	OldPrice               *float64                `json:"oldPrice,omitempty"`
	SpecialPrice           *float64                `json:"specialPrice,omitempty"`
	SpecialPriceStart      *time.Time              `json:"specialPriceStart,omitempty"`
	SpecialPriceEnd        *time.Time              `json:"specialPriceEnd,omitempty"`
	HasOptions             bool                    `json:"hasOptions"`
	IsVisibleIndividually  bool                    `json:"isVisibleIndividually"`
	IsFeatured             bool                    `json:"isFeatured"`
	IsCallForPricing       bool                    `json:"isCallForPricing"`
	IsAllowToOrder         bool                    `json:"isAllowToOrder"`
	StockTrackingIsEnabled bool                    `json:"stockTrackingIsEnabled"`
	StockQuantity          int                     `json:"stockQuantity"`
	SKU                    string                  `json:"sku"`
	GTIN                   string                  `json:"gtin"`
	NormalizedName         string                  `json:"normalizedName"`
	DisplayOrder           int                     `json:"displayOrder"`
	ReviewsCount           int                     `json:"reviewsCount"`
	RatingAverage          *float64                `json:"ratingAverage,omitempty"`
	BrandID                *uint                   `json:"brandId,omitempty"`
	BrandName              string                  `json:"brandName,omitempty"`
	TaxClassID             *uint                   `json:"taxClassId,omitempty"`
	ThumbnailImageID       *uint                   `json:"thumbnailImageId,omitempty"`
	IsPublished            bool                    `json:"isPublished"`
	PublishedOn            *time.Time              `json:"publishedOn,omitempty"`
	CreatedAt              time.Time               `json:"createdAt"`
	UpdatedAt              time.Time               `json:"updatedAt"`
	Categories             []*CategoryBrief        `json:"categories,omitempty"`
	AttributeValues        []*AttributeValueResponse `json:"attributeValues,omitempty"`
	OptionValues           []*OptionValueResponse    `json:"optionValues,omitempty"`
	Medias                 []*MediaResponse          `json:"medias,omitempty"`
}

type ProductListItem struct {
	ID                     uint       `json:"id"`
	Name                   string     `json:"name"`
	Slug                   string     `json:"slug"`
	ShortDescription       string     `json:"shortDescription"`
	Price                  float64    `json:"price"`
	OldPrice               *float64   `json:"oldPrice,omitempty"`
	SpecialPrice           *float64   `json:"specialPrice,omitempty"`
	IsFeatured             bool       `json:"isFeatured"`
	StockQuantity          int        `json:"stockQuantity"`
	SKU                    string     `json:"sku"`
	GTIN                   string     `json:"gtin"`
	BrandID                *uint      `json:"brandId,omitempty"`
	BrandName              string     `json:"brandName,omitempty"`
	ThumbnailImageID       *uint      `json:"thumbnailImageId,omitempty"`
	IsPublished            bool       `json:"isPublished"`
	CreatedAt              time.Time  `json:"createdAt"`
	UpdatedAt              time.Time  `json:"updatedAt"`
}

type CreateProductRequest struct {
	Name                  string             `json:"name"`
	Slug                  string             `json:"slug,omitempty"`
	ShortDescription      string             `json:"shortDescription,omitempty"`
	Description           string             `json:"description,omitempty"`
	Specification         string             `json:"specification,omitempty"`
	Price                 float64            `json:"price"`
	OldPrice              *float64           `json:"oldPrice,omitempty"`
	SpecialPrice          *float64           `json:"specialPrice,omitempty"`
	SpecialPriceStart     *string            `json:"specialPriceStart,omitempty"`
	SpecialPriceEnd       *string            `json:"specialPriceEnd,omitempty"`
	HasOptions            bool               `json:"hasOptions,omitempty"`
	IsVisibleIndividually bool               `json:"isVisibleIndividually,omitempty"`
	IsFeatured            bool               `json:"isFeatured,omitempty"`
	IsCallForPricing      bool               `json:"isCallForPricing,omitempty"`
	IsAllowToOrder        bool               `json:"isAllowToOrder,omitempty"`
	StockTrackingIsEnabled bool              `json:"stockTrackingIsEnabled,omitempty"`
	StockQuantity         int                `json:"stockQuantity,omitempty"`
	SKU                   string             `json:"sku,omitempty"`
	GTIN                  string             `json:"gtin,omitempty"`
	DisplayOrder          int                `json:"displayOrder,omitempty"`
	BrandID               *uint              `json:"brandId,omitempty"`
	TaxClassID            *uint              `json:"taxClassId,omitempty"`
	ThumbnailImageID      *uint              `json:"thumbnailImageId,omitempty"`
	IsPublished           bool               `json:"isPublished,omitempty"`
	CategoryIDs           []uint             `json:"categoryIds,omitempty"`
	AttributeValues       []AttributeValueReq `json:"attributeValues,omitempty"`
	OptionValues          []OptionValueReq    `json:"optionValues,omitempty"`
}

type AttributeValueReq struct {
	AttributeID uint   `json:"attributeId"`
	Value       string `json:"value"`
}

type OptionValueReq struct {
	OptionID    uint   `json:"optionId"`
	Value       string `json:"value"`
	DisplayType string `json:"displayType"`
	SortIndex   int    `json:"sortIndex"`
}

type UpdateProductRequest struct {
	Name                  string             `json:"name"`
	Slug                  string             `json:"slug,omitempty"`
	ShortDescription      string             `json:"shortDescription,omitempty"`
	Description           string             `json:"description,omitempty"`
	Specification         string             `json:"specification,omitempty"`
	Price                 float64            `json:"price"`
	OldPrice              *float64           `json:"oldPrice,omitempty"`
	SpecialPrice          *float64           `json:"specialPrice,omitempty"`
	SpecialPriceStart     *string            `json:"specialPriceStart,omitempty"`
	SpecialPriceEnd       *string            `json:"specialPriceEnd,omitempty"`
	HasOptions            bool               `json:"hasOptions,omitempty"`
	IsVisibleIndividually bool               `json:"isVisibleIndividually,omitempty"`
	IsFeatured            bool               `json:"isFeatured,omitempty"`
	IsCallForPricing      bool               `json:"isCallForPricing,omitempty"`
	IsAllowToOrder        bool               `json:"isAllowToOrder,omitempty"`
	StockTrackingIsEnabled bool              `json:"stockTrackingIsEnabled,omitempty"`
	StockQuantity         int                `json:"stockQuantity,omitempty"`
	SKU                   string             `json:"sku,omitempty"`
	GTIN                  string             `json:"gtin,omitempty"`
	DisplayOrder          int                `json:"displayOrder,omitempty"`
	BrandID               *uint              `json:"brandId,omitempty"`
	TaxClassID            *uint              `json:"taxClassId,omitempty"`
	ThumbnailImageID      *uint              `json:"thumbnailImageId,omitempty"`
	IsPublished           bool               `json:"isPublished,omitempty"`
	CategoryIDs           []uint             `json:"categoryIds,omitempty"`
	AttributeValues       []AttributeValueReq `json:"attributeValues,omitempty"`
	OptionValues          []OptionValueReq    `json:"optionValues,omitempty"`
}

func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{db: db}
}

func parseProductTime(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return nil
	}
	return &t
}

func (s *ProductService) scanProductListItem(row scannable) (*ProductListItem, error) {
	var p ProductListItem
	var oldPrice, specialPrice sql.NullFloat64
	var brandID, thumbID sql.NullInt64
	var brandName sql.NullString

	err := row.Scan(
		&p.ID, &p.Name, &p.Slug, &p.ShortDescription,
		&p.Price, &oldPrice, &specialPrice,
		&p.IsFeatured, &p.StockQuantity,
		&p.SKU, &p.GTIN,
		&brandID, &brandName, &thumbID,
		&p.IsPublished,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if oldPrice.Valid {
		p.OldPrice = &oldPrice.Float64
	}
	if specialPrice.Valid {
		p.SpecialPrice = &specialPrice.Float64
	}
	if brandID.Valid {
		bid := uint(brandID.Int64)
		p.BrandID = &bid
	}
	if brandName.Valid {
		p.BrandName = brandName.String
	}
	if thumbID.Valid {
		tid := uint(thumbID.Int64)
		p.ThumbnailImageID = &tid
	}
	return &p, nil
}

func (s *ProductService) GetProducts(ctx context.Context, page, pageSize int, search string, categoryID, brandID *uint, minPrice, maxPrice *float64) ([]*ProductListItem, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	where := "p.deleted_at IS NULL"
	args := []interface{}{}
	argIdx := 1

	if search != "" {
		where += " AND (p.name ILIKE '%' || $" + itoa(argIdx) + " || '%' OR p.short_description ILIKE '%' || $" + itoa(argIdx) + " || '%' OR p.sku ILIKE '%' || $" + itoa(argIdx) + " || '%')"
		args = append(args, search)
		argIdx++
	}
	if categoryID != nil && *categoryID > 0 {
		where += " AND EXISTS (SELECT 1 FROM catalog_product_categories pc WHERE pc.product_id = p.id AND pc.category_id = $" + itoa(argIdx) + ")"
		args = append(args, *categoryID)
		argIdx++
	}
	if brandID != nil && *brandID > 0 {
		where += " AND p.brand_id = $" + itoa(argIdx)
		args = append(args, *brandID)
		argIdx++
	}
	if minPrice != nil {
		where += " AND p.price >= $" + itoa(argIdx)
		args = append(args, *minPrice)
		argIdx++
	}
	if maxPrice != nil {
		where += " AND p.price <= $" + itoa(argIdx)
		args = append(args, *maxPrice)
		argIdx++
	}

	var total int64
	countQuery := "SELECT COUNT(*) FROM catalog_products p WHERE " + where
	err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)
	dataQuery := `SELECT p.id, p.name, p.slug, COALESCE(p.short_description, ''),
		p.price, p.old_price, p.special_price,
		p.is_featured, p.stock_quantity,
		COALESCE(p.sku, ''), COALESCE(p.gtin, ''),
		p.brand_id, b.name, p.thumbnail_image_id,
		p.is_published, p.created_at, p.updated_at
		FROM catalog_products p
		LEFT JOIN catalog_brands b ON p.brand_id = b.id
		WHERE ` + where + ` ORDER BY p.created_at DESC LIMIT $` + itoa(argIdx) + ` OFFSET $` + itoa(argIdx+1)
	argIdx += 2

	rows, err := s.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*ProductListItem
	for rows.Next() {
		p, err := s.scanProductListItem(rows)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, p)
	}
	if products == nil {
		products = []*ProductListItem{}
	}
	return products, total, nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

func (s *ProductService) GetProductBySlug(ctx context.Context, slug string) (*ProductResponse, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT p.id, p.name, p.slug, COALESCE(p.short_description, ''), COALESCE(p.description, ''), COALESCE(p.specification, ''),
			p.price, p.old_price, p.special_price, p.special_price_start, p.special_price_end,
			p.has_options, p.is_visible_individually, p.is_featured, p.is_call_for_pricing, p.is_allow_to_order,
			p.stock_tracking_is_enabled, p.stock_quantity, COALESCE(p.sku, ''), COALESCE(p.gtin, ''), COALESCE(p.normalized_name, ''),
			p.display_order, p.reviews_count, p.rating_average,
			p.brand_id, p.tax_class_id, p.thumbnail_image_id,
			p.is_published, p.published_on,
			p.created_at, p.updated_at,
			COALESCE(b.name, '') as brand_name
		 FROM catalog_products p
		 LEFT JOIN catalog_brands b ON p.brand_id = b.id
		 WHERE p.slug = $1 AND p.deleted_at IS NULL`,
		slug,
	)

	var p ProductResponse
	var oldPrice, specialPrice, ratingAvg sql.NullFloat64
	var spStart, spEnd, pubOn sql.NullTime
	var brandID, taxClassID, thumbID sql.NullInt64

	err := row.Scan(
		&p.ID, &p.Name, &p.Slug, &p.ShortDescription, &p.Description, &p.Specification,
		&p.Price, &oldPrice, &specialPrice, &spStart, &spEnd,
		&p.HasOptions, &p.IsVisibleIndividually, &p.IsFeatured, &p.IsCallForPricing, &p.IsAllowToOrder,
		&p.StockTrackingIsEnabled, &p.StockQuantity, &p.SKU, &p.GTIN, &p.NormalizedName,
		&p.DisplayOrder, &p.ReviewsCount, &ratingAvg,
		&brandID, &taxClassID, &thumbID,
		&p.IsPublished, &pubOn,
		&p.CreatedAt, &p.UpdatedAt,
		&p.BrandName,
	)
	if err == sql.ErrNoRows {
		return nil, fiber.NewError(fiber.StatusNotFound, "Product not found")
	}
	if err != nil {
		return nil, err
	}

	if oldPrice.Valid {
		p.OldPrice = &oldPrice.Float64
	}
	if specialPrice.Valid {
		p.SpecialPrice = &specialPrice.Float64
	}
	if ratingAvg.Valid {
		p.RatingAverage = &ratingAvg.Float64
	}
	if spStart.Valid {
		p.SpecialPriceStart = &spStart.Time
	}
	if spEnd.Valid {
		p.SpecialPriceEnd = &spEnd.Time
	}
	if pubOn.Valid {
		p.PublishedOn = &pubOn.Time
	}
	if brandID.Valid {
		bid := uint(brandID.Int64)
		p.BrandID = &bid
	}
	if taxClassID.Valid {
		tid := uint(taxClassID.Int64)
		p.TaxClassID = &tid
	}
	if thumbID.Valid {
		tid := uint(thumbID.Int64)
		p.ThumbnailImageID = &tid
	}

	if err := s.loadProductRelations(ctx, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *ProductService) loadProductRelations(ctx context.Context, p *ProductResponse) error {
	catRows, err := s.db.QueryContext(ctx,
		"SELECT c.id, c.name, c.slug FROM catalog_categories c JOIN catalog_product_categories pc ON c.id = pc.category_id WHERE pc.product_id = $1 AND c.is_deleted = false ORDER BY c.display_order ASC",
		p.ID,
	)
	if err != nil {
		return err
	}
	defer catRows.Close()
	for catRows.Next() {
		var c CategoryBrief
		if err := catRows.Scan(&c.ID, &c.Name, &c.Slug); err != nil {
			return err
		}
		p.Categories = append(p.Categories, &c)
	}
	if p.Categories == nil {
		p.Categories = []*CategoryBrief{}
	}

	attrRows, err := s.db.QueryContext(ctx,
		`SELECT pav.id, pav.attribute_id, COALESCE(pa.name, ''), pav.value
		 FROM catalog_product_attribute_values pav
		 LEFT JOIN catalog_product_attributes pa ON pav.attribute_id = pa.id
		 WHERE pav.product_id = $1 ORDER BY pav.id ASC`,
		p.ID,
	)
	if err != nil {
		return err
	}
	defer attrRows.Close()
	for attrRows.Next() {
		var a AttributeValueResponse
		if err := attrRows.Scan(&a.ID, &a.AttributeID, &a.AttributeName, &a.Value); err != nil {
			return err
		}
		p.AttributeValues = append(p.AttributeValues, &a)
	}
	if p.AttributeValues == nil {
		p.AttributeValues = []*AttributeValueResponse{}
	}

	optRows, err := s.db.QueryContext(ctx,
		`SELECT pov.id, pov.option_id, COALESCE(po.name, ''), pov.value, COALESCE(pov.display_type, ''), pov.sort_index
		 FROM catalog_product_option_values pov
		 LEFT JOIN catalog_product_options po ON pov.option_id = po.id
		 WHERE pov.product_id = $1 ORDER BY pov.sort_index ASC`,
		p.ID,
	)
	if err != nil {
		return err
	}
	defer optRows.Close()
	for optRows.Next() {
		var o OptionValueResponse
		if err := optRows.Scan(&o.ID, &o.OptionID, &o.OptionName, &o.Value, &o.DisplayType, &o.SortIndex); err != nil {
			return err
		}
		p.OptionValues = append(p.OptionValues, &o)
	}
	if p.OptionValues == nil {
		p.OptionValues = []*OptionValueResponse{}
	}

	medRows, err := s.db.QueryContext(ctx,
		`SELECT m.id, COALESCE(m.caption, ''), COALESCE(m.file_name, ''), m.file_size, m.media_type
		 FROM catalog_product_medias pm
		 JOIN catalog_media m ON pm.media_id = m.id
		 WHERE pm.product_id = $1 ORDER BY pm.display_order ASC`,
		p.ID,
	)
	if err != nil {
		return err
	}
	defer medRows.Close()
	for medRows.Next() {
		var m MediaResponse
		if err := medRows.Scan(&m.ID, &m.Caption, &m.FileName, &m.FileSize, &m.MediaType); err != nil {
			return err
		}
		p.Medias = append(p.Medias, &m)
	}
	if p.Medias == nil {
		p.Medias = []*MediaResponse{}
	}

	return nil
}

func (s *ProductService) GetProductByID(ctx context.Context, id uint) (*ProductResponse, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT p.id, p.name, p.slug, COALESCE(p.short_description, ''), COALESCE(p.description, ''), COALESCE(p.specification, ''),
			p.price, p.old_price, p.special_price, p.special_price_start, p.special_price_end,
			p.has_options, p.is_visible_individually, p.is_featured, p.is_call_for_pricing, p.is_allow_to_order,
			p.stock_tracking_is_enabled, p.stock_quantity, COALESCE(p.sku, ''), COALESCE(p.gtin, ''), COALESCE(p.normalized_name, ''),
			p.display_order, p.reviews_count, p.rating_average,
			p.brand_id, p.tax_class_id, p.thumbnail_image_id,
			p.is_published, p.published_on,
			p.created_at, p.updated_at,
			COALESCE(b.name, '') as brand_name
		 FROM catalog_products p
		 LEFT JOIN catalog_brands b ON p.brand_id = b.id
		 WHERE p.id = $1 AND p.deleted_at IS NULL`,
		id,
	)

	var p ProductResponse
	var oldPrice, specialPrice, ratingAvg sql.NullFloat64
	var spStart, spEnd, pubOn sql.NullTime
	var brandID, taxClassID, thumbID sql.NullInt64

	err := row.Scan(
		&p.ID, &p.Name, &p.Slug, &p.ShortDescription, &p.Description, &p.Specification,
		&p.Price, &oldPrice, &specialPrice, &spStart, &spEnd,
		&p.HasOptions, &p.IsVisibleIndividually, &p.IsFeatured, &p.IsCallForPricing, &p.IsAllowToOrder,
		&p.StockTrackingIsEnabled, &p.StockQuantity, &p.SKU, &p.GTIN, &p.NormalizedName,
		&p.DisplayOrder, &p.ReviewsCount, &ratingAvg,
		&brandID, &taxClassID, &thumbID,
		&p.IsPublished, &pubOn,
		&p.CreatedAt, &p.UpdatedAt,
		&p.BrandName,
	)
	if err == sql.ErrNoRows {
		return nil, fiber.NewError(fiber.StatusNotFound, "Product not found")
	}
	if err != nil {
		return nil, err
	}

	if oldPrice.Valid {
		p.OldPrice = &oldPrice.Float64
	}
	if specialPrice.Valid {
		p.SpecialPrice = &specialPrice.Float64
	}
	if ratingAvg.Valid {
		p.RatingAverage = &ratingAvg.Float64
	}
	if spStart.Valid {
		p.SpecialPriceStart = &spStart.Time
	}
	if spEnd.Valid {
		p.SpecialPriceEnd = &spEnd.Time
	}
	if pubOn.Valid {
		p.PublishedOn = &pubOn.Time
	}
	if brandID.Valid {
		bid := uint(brandID.Int64)
		p.BrandID = &bid
	}
	if taxClassID.Valid {
		tid := uint(taxClassID.Int64)
		p.TaxClassID = &tid
	}
	if thumbID.Valid {
		tid := uint(thumbID.Int64)
		p.ThumbnailImageID = &tid
	}

	if err := s.loadProductRelations(ctx, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *ProductService) GetFeaturedProducts(ctx context.Context, limit int) ([]*ProductListItem, error) {
	if limit < 1 {
		limit = 10
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT p.id, p.name, p.slug, COALESCE(p.short_description, ''),
			p.price, p.old_price, p.special_price,
			p.is_featured, p.stock_quantity,
			COALESCE(p.sku, ''), COALESCE(p.gtin, ''),
			p.brand_id, b.name, p.thumbnail_image_id,
			p.is_published, p.created_at, p.updated_at
		 FROM catalog_products p
		 LEFT JOIN catalog_brands b ON p.brand_id = b.id
		 WHERE p.is_featured = true AND p.is_published = true AND p.deleted_at IS NULL
		 ORDER BY p.created_at DESC LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*ProductListItem
	for rows.Next() {
		p, err := s.scanProductListItem(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if products == nil {
		products = []*ProductListItem{}
	}
	return products, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, req *CreateProductRequest) (*ProductResponse, error) {
	slug := req.Slug
	if slug == "" {
		slug = GenerateSlug(req.Name)
	}

	normalizedName := req.Name
	for _, r := range normalizedName {
		if r >= 'a' && r <= 'z' {
			normalizedName = string(r - 32) + normalizedName[1:]
		}
		break
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var productID uint
	var createdAt, updatedAt time.Time

	var spStart, spEnd *time.Time
	if req.SpecialPriceStart != nil {
		spStart = parseProductTime(req.SpecialPriceStart)
	}
	if req.SpecialPriceEnd != nil {
		spEnd = parseProductTime(req.SpecialPriceEnd)
	}

	var publishedOn *time.Time
	if req.IsPublished {
		now := time.Now()
		publishedOn = &now
	}

	err = tx.QueryRowContext(ctx,
		`INSERT INTO catalog_products (name, slug, short_description, description, specification,
			price, old_price, special_price, special_price_start, special_price_end,
			has_options, is_visible_individually, is_featured, is_call_for_pricing, is_allow_to_order,
			stock_tracking_is_enabled, stock_quantity, sku, gtin, normalized_name,
			display_order, brand_id, tax_class_id, thumbnail_image_id,
			is_published, published_on, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,NOW(),NOW())
		 RETURNING id, created_at, updated_at`,
		req.Name, slug, req.ShortDescription, req.Description, req.Specification,
		req.Price, req.OldPrice, req.SpecialPrice, spStart, spEnd,
		req.HasOptions, req.IsVisibleIndividually, req.IsFeatured, req.IsCallForPricing, req.IsAllowToOrder,
		req.StockTrackingIsEnabled, req.StockQuantity, req.SKU, req.GTIN, normalizedName,
		req.DisplayOrder, req.BrandID, req.TaxClassID, req.ThumbnailImageID,
		req.IsPublished, publishedOn,
	).Scan(&productID, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	for _, catID := range req.CategoryIDs {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO catalog_product_categories (product_id, category_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
			productID, catID,
		)
		if err != nil {
			return nil, err
		}
	}

	for _, av := range req.AttributeValues {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO catalog_product_attribute_values (attribute_id, product_id, value, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())",
			av.AttributeID, productID, av.Value,
		)
		if err != nil {
			return nil, err
		}
	}

	for _, ov := range req.OptionValues {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO catalog_product_option_values (option_id, product_id, value, display_type, sort_index, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())",
			ov.OptionID, productID, ov.Value, ov.DisplayType, ov.SortIndex,
		)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetProductByID(ctx, productID)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id uint, req *UpdateProductRequest) (*ProductResponse, error) {
	existing, err := s.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}

	name := existing.Name
	if req.Name != "" {
		name = req.Name
	}
	slug := req.Slug
	if slug == "" && req.Name != "" {
		slug = GenerateSlug(name)
	} else if slug == "" {
		slug = existing.Slug
	}

	shortDesc := req.ShortDescription
	desc := req.Description
	spec := req.Specification
	price := req.Price
	oldPrice := req.OldPrice
	specialPrice := req.SpecialPrice
	hasOptions := req.HasOptions
	isVisible := req.IsVisibleIndividually
	isFeatured := req.IsFeatured
	isCallPricing := req.IsCallForPricing
	isAllowOrder := req.IsAllowToOrder
	stockTracking := req.StockTrackingIsEnabled
	stockQty := req.StockQuantity
	sku := req.SKU
	gtin := req.GTIN
	dispOrder := req.DisplayOrder
	brandID := req.BrandID
	taxClassID := req.TaxClassID
	thumbID := req.ThumbnailImageID
	isPub := req.IsPublished

	var spStart, spEnd *time.Time
	if req.SpecialPriceStart != nil {
		spStart = parseProductTime(req.SpecialPriceStart)
	}
	if req.SpecialPriceEnd != nil {
		spEnd = parseProductTime(req.SpecialPriceEnd)
	}

	var publishedOn *time.Time
	if req.IsPublished && existing.PublishedOn == nil {
		now := time.Now()
		publishedOn = &now
	} else {
		publishedOn = existing.PublishedOn
	}

	normalizedName := name
	for i, r := range name {
		if r >= 'a' && r <= 'z' {
			normalizedName = string(r-32) + name[i+1:]
		} else {
			normalizedName = name
		}
		break
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx,
		`UPDATE catalog_products SET
			name = $1, slug = $2, short_description = $3, description = $4, specification = $5,
			price = $6, old_price = $7, special_price = $8, special_price_start = $9, special_price_end = $10,
			has_options = $11, is_visible_individually = $12, is_featured = $13, is_call_for_pricing = $14, is_allow_to_order = $15,
			stock_tracking_is_enabled = $16, stock_quantity = $17, sku = $18, gtin = $19, normalized_name = $20,
			display_order = $21, brand_id = $22, tax_class_id = $23, thumbnail_image_id = $24,
			is_published = $25, published_on = $26, updated_at = NOW()
		 WHERE id = $27 AND deleted_at IS NULL`,
		name, slug, shortDesc, desc, spec,
		price, oldPrice, specialPrice, spStart, spEnd,
		hasOptions, isVisible, isFeatured, isCallPricing, isAllowOrder,
		stockTracking, stockQty, sku, gtin, normalizedName,
		dispOrder, brandID, taxClassID, thumbID,
		isPub, publishedOn, id,
	)
	if err != nil {
		return nil, err
	}
	rowsAff, _ := result.RowsAffected()
	if rowsAff == 0 {
		return nil, fiber.NewError(fiber.StatusNotFound, "Product not found")
	}

	if req.CategoryIDs != nil {
		_, err = tx.ExecContext(ctx, "DELETE FROM catalog_product_categories WHERE product_id = $1", id)
		if err != nil {
			return nil, err
		}
		for _, catID := range req.CategoryIDs {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO catalog_product_categories (product_id, category_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
				id, catID,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	if req.AttributeValues != nil {
		_, err = tx.ExecContext(ctx, "DELETE FROM catalog_product_attribute_values WHERE product_id = $1", id)
		if err != nil {
			return nil, err
		}
		for _, av := range req.AttributeValues {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO catalog_product_attribute_values (attribute_id, product_id, value, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())",
				av.AttributeID, id, av.Value,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	if req.OptionValues != nil {
		_, err = tx.ExecContext(ctx, "DELETE FROM catalog_product_option_values WHERE product_id = $1", id)
		if err != nil {
			return nil, err
		}
		for _, ov := range req.OptionValues {
			_, err = tx.ExecContext(ctx,
				"INSERT INTO catalog_product_option_values (option_id, product_id, value, display_type, sort_index, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())",
				ov.OptionID, id, ov.Value, ov.DisplayType, ov.SortIndex,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetProductByID(ctx, id)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uint) error {
	result, err := s.db.ExecContext(ctx, "UPDATE catalog_products SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Product not found")
	}
	return nil
}

func (s *ProductService) SearchProducts(ctx context.Context, query string, page, pageSize int, categoryID *uint, minPrice, maxPrice *float64) ([]*ProductListItem, int64, error) {
	return s.GetProducts(ctx, page, pageSize, query, categoryID, nil, minPrice, maxPrice)
}
