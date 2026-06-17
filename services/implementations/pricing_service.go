package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PricingService struct {
	db *sql.DB
}

func NewPricingService(db *sql.DB) *PricingService {
	return &PricingService{db: db}
}

type CartRuleResponse struct {
	ID                   uint                      `json:"id"`
	Name                 string                    `json:"name"`
	Description          string                    `json:"description"`
	IsActive             bool                      `json:"isActive"`
	StartOn              *time.Time                `json:"startOn,omitempty"`
	EndOn                *time.Time                `json:"endOn,omitempty"`
	IsCouponRequired     bool                      `json:"isCouponRequired"`
	RuleToApply          string                    `json:"ruleToApply"`
	DiscountAmount       float64                   `json:"discountAmount"`
	MaxDiscountAmount    *float64                  `json:"maxDiscountAmount,omitempty"`
	DiscountStep         *int                      `json:"discountStep,omitempty"`
	UsageLimitPerCoupon  *int                      `json:"usageLimitPerCoupon,omitempty"`
	UsageLimitPerCustomer *int                     `json:"usageLimitPerCustomer,omitempty"`
	CreatedAt            time.Time                 `json:"createdAt"`
	UpdatedAt            time.Time                 `json:"updatedAt"`
	Coupons              []CouponResponse          `json:"coupons,omitempty"`
	Categories           []uint                    `json:"categories,omitempty"`
	Products             []uint                    `json:"products,omitempty"`
	CustomerGroups       []uint                    `json:"customerGroups,omitempty"`
}

type CouponResponse struct {
	ID         uint      `json:"id"`
	CartRuleID uint      `json:"cartRuleId"`
	Code       string    `json:"code"`
	CreatedAt  time.Time `json:"createdAt"`
}

type CatalogRuleResponse struct {
	ID                uint       `json:"id"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	IsActive          bool       `json:"isActive"`
	StartOn           *time.Time `json:"startOn,omitempty"`
	EndOn             *time.Time `json:"endOn,omitempty"`
	RuleToApply       string     `json:"ruleToApply"`
	DiscountAmount    float64    `json:"discountAmount"`
	MaxDiscountAmount *float64   `json:"maxDiscountAmount,omitempty"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	CustomerGroups    []uint     `json:"customerGroups,omitempty"`
}

func (s *PricingService) GetCartRules(ctx context.Context) ([]CartRuleResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, COALESCE(description,''), is_active, start_on, end_on, is_coupon_required,
			rule_to_apply, discount_amount, max_discount_amount, discount_step, usage_limit_per_coupon,
			usage_limit_per_customer, created_at, updated_at
		FROM pricing_cart_rules ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list cart rules: %w", err)
	}
	defer rows.Close()

	var rules []CartRuleResponse
	for rows.Next() {
		var cr CartRuleResponse
		var startOn, endOn sql.NullTime
		var maxDiscount sql.NullFloat64
		var discountStep, usagePerCoupon, usagePerCustomer sql.NullInt64
		if err := rows.Scan(&cr.ID, &cr.Name, &cr.Description, &cr.IsActive, &startOn, &endOn,
			&cr.IsCouponRequired, &cr.RuleToApply, &cr.DiscountAmount, &maxDiscount, &discountStep,
			&usagePerCoupon, &usagePerCustomer, &cr.CreatedAt, &cr.UpdatedAt); err != nil {
			return nil, err
		}
		if startOn.Valid {
			cr.StartOn = &startOn.Time
		}
		if endOn.Valid {
			cr.EndOn = &endOn.Time
		}
		if maxDiscount.Valid {
			cr.MaxDiscountAmount = &maxDiscount.Float64
		}
		if discountStep.Valid {
			v := int(discountStep.Int64)
			cr.DiscountStep = &v
		}
		if usagePerCoupon.Valid {
			v := int(usagePerCoupon.Int64)
			cr.UsageLimitPerCoupon = &v
		}
		if usagePerCustomer.Valid {
			v := int(usagePerCustomer.Int64)
			cr.UsageLimitPerCustomer = &v
		}
		rules = append(rules, cr)
	}
	if rules == nil {
		rules = []CartRuleResponse{}
	}
	return rules, nil
}

func (s *PricingService) GetCartRule(ctx context.Context, id uint) (*CartRuleResponse, error) {
	var cr CartRuleResponse
	var startOn, endOn sql.NullTime
	var maxDiscount sql.NullFloat64
	var discountStep, usagePerCoupon, usagePerCustomer sql.NullInt64
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, COALESCE(description,''), is_active, start_on, end_on, is_coupon_required,
			rule_to_apply, discount_amount, max_discount_amount, discount_step, usage_limit_per_coupon,
			usage_limit_per_customer, created_at, updated_at
		FROM pricing_cart_rules WHERE id = $1
	`, id).Scan(&cr.ID, &cr.Name, &cr.Description, &cr.IsActive, &startOn, &endOn,
		&cr.IsCouponRequired, &cr.RuleToApply, &cr.DiscountAmount, &maxDiscount, &discountStep,
		&usagePerCoupon, &usagePerCustomer, &cr.CreatedAt, &cr.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cart rule not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cart rule: %w", err)
	}
	if startOn.Valid {
		cr.StartOn = &startOn.Time
	}
	if endOn.Valid {
		cr.EndOn = &endOn.Time
	}
	if maxDiscount.Valid {
		cr.MaxDiscountAmount = &maxDiscount.Float64
	}
	if discountStep.Valid {
		v := int(discountStep.Int64)
		cr.DiscountStep = &v
	}
	if usagePerCoupon.Valid {
		v := int(usagePerCoupon.Int64)
		cr.UsageLimitPerCoupon = &v
	}
	if usagePerCustomer.Valid {
		v := int(usagePerCustomer.Int64)
		cr.UsageLimitPerCustomer = &v
	}

	coupons, _ := s.getCoupons(ctx, id)
	cr.Coupons = coupons

	categories, _ := s.getCartRuleCategories(ctx, id)
	cr.Categories = categories

	products, _ := s.getCartRuleProducts(ctx, id)
	cr.Products = products

	groups, _ := s.getCartRuleCustomerGroups(ctx, id)
	cr.CustomerGroups = groups

	return &cr, nil
}

func (s *PricingService) CreateCartRule(ctx context.Context, req CreateCartRuleRequest) (*CartRuleResponse, error) {
	var cr CartRuleResponse
	var startOn, endOn sql.NullTime
	var maxDiscount sql.NullFloat64
	var discountStep, usagePerCoupon, usagePerCustomer sql.NullInt64
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO pricing_cart_rules (name, description, is_active, start_on, end_on, is_coupon_required,
			rule_to_apply, discount_amount, max_discount_amount, discount_step, usage_limit_per_coupon,
			usage_limit_per_customer, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,NOW(),NOW())
		RETURNING id, name, description, is_active, start_on, end_on, is_coupon_required,
			rule_to_apply, discount_amount, max_discount_amount, discount_step, usage_limit_per_coupon,
			usage_limit_per_customer, created_at, updated_at
	`, req.Name, req.Description, req.IsActive, req.StartOn, req.EndOn, req.IsCouponRequired,
		req.RuleToApply, req.DiscountAmount, req.MaxDiscountAmount, req.DiscountStep,
		req.UsageLimitPerCoupon, req.UsageLimitPerCustomer).Scan(
		&cr.ID, &cr.Name, &cr.Description, &cr.IsActive, &startOn, &endOn,
		&cr.IsCouponRequired, &cr.RuleToApply, &cr.DiscountAmount, &maxDiscount, &discountStep,
		&usagePerCoupon, &usagePerCustomer, &cr.CreatedAt, &cr.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create cart rule: %w", err)
	}
	if startOn.Valid {
		cr.StartOn = &startOn.Time
	}
	if endOn.Valid {
		cr.EndOn = &endOn.Time
	}
	if maxDiscount.Valid {
		cr.MaxDiscountAmount = &maxDiscount.Float64
	}
	if discountStep.Valid {
		v := int(discountStep.Int64)
		cr.DiscountStep = &v
	}
	if usagePerCoupon.Valid {
		v := int(usagePerCoupon.Int64)
		cr.UsageLimitPerCoupon = &v
	}
	if usagePerCustomer.Valid {
		v := int(usagePerCustomer.Int64)
		cr.UsageLimitPerCustomer = &v
	}
	cr.Coupons = []CouponResponse{}
	cr.Categories = []uint{}
	cr.Products = []uint{}
	cr.CustomerGroups = []uint{}
	return &cr, nil
}

func (s *PricingService) UpdateCartRule(ctx context.Context, id uint, req CreateCartRuleRequest) (*CartRuleResponse, error) {
	var cr CartRuleResponse
	var startOn, endOn sql.NullTime
	var maxDiscount sql.NullFloat64
	var discountStep, usagePerCoupon, usagePerCustomer sql.NullInt64
	err := s.db.QueryRowContext(ctx, `
		UPDATE pricing_cart_rules SET name=$1, description=$2, is_active=$3, start_on=$4, end_on=$5,
			is_coupon_required=$6, rule_to_apply=$7, discount_amount=$8, max_discount_amount=$9,
			discount_step=$10, usage_limit_per_coupon=$11, usage_limit_per_customer=$12, updated_at=NOW()
		WHERE id=$13
		RETURNING id, name, description, is_active, start_on, end_on, is_coupon_required,
			rule_to_apply, discount_amount, max_discount_amount, discount_step, usage_limit_per_coupon,
			usage_limit_per_customer, created_at, updated_at
	`, req.Name, req.Description, req.IsActive, req.StartOn, req.EndOn, req.IsCouponRequired,
		req.RuleToApply, req.DiscountAmount, req.MaxDiscountAmount, req.DiscountStep,
		req.UsageLimitPerCoupon, req.UsageLimitPerCustomer, id).Scan(
		&cr.ID, &cr.Name, &cr.Description, &cr.IsActive, &startOn, &endOn,
		&cr.IsCouponRequired, &cr.RuleToApply, &cr.DiscountAmount, &maxDiscount, &discountStep,
		&usagePerCoupon, &usagePerCustomer, &cr.CreatedAt, &cr.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("cart rule not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update cart rule: %w", err)
	}
	if startOn.Valid {
		cr.StartOn = &startOn.Time
	}
	if endOn.Valid {
		cr.EndOn = &endOn.Time
	}
	if maxDiscount.Valid {
		cr.MaxDiscountAmount = &maxDiscount.Float64
	}
	if discountStep.Valid {
		v := int(discountStep.Int64)
		cr.DiscountStep = &v
	}
	if usagePerCoupon.Valid {
		v := int(usagePerCoupon.Int64)
		cr.UsageLimitPerCoupon = &v
	}
	if usagePerCustomer.Valid {
		v := int(usagePerCustomer.Int64)
		cr.UsageLimitPerCustomer = &v
	}
	return &cr, nil
}

func (s *PricingService) DeleteCartRule(ctx context.Context, id uint) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `DELETE FROM pricing_cart_rule_categories WHERE cart_rule_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `DELETE FROM pricing_cart_rule_products WHERE cart_rule_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `DELETE FROM pricing_cart_rule_customer_groups WHERE cart_rule_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `DELETE FROM pricing_coupons WHERE cart_rule_id = $1`, id)
	if err != nil {
		return err
	}
	result, err := tx.ExecContext(ctx, `DELETE FROM pricing_cart_rules WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete cart rule: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("cart rule not found")
	}
	return tx.Commit()
}

func (s *PricingService) ValidateCoupon(ctx context.Context, code string) (*CouponValidationResult, error) {
	var coupon CouponResponse
	err := s.db.QueryRowContext(ctx, `
		SELECT id, cart_rule_id, code, created_at FROM pricing_coupons WHERE code = $1
	`, code).Scan(&coupon.ID, &coupon.CartRuleID, &coupon.Code, &coupon.CreatedAt)
	if err == sql.ErrNoRows {
		return &CouponValidationResult{Valid: false, Error: "coupon not found"}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find coupon: %w", err)
	}

	var cr CartRuleResponse
	var startOn, endOn sql.NullTime
	var maxDiscount sql.NullFloat64
	var discountStep, usagePerCoupon, usagePerCustomer sql.NullInt64
	err = s.db.QueryRowContext(ctx, `
		SELECT id, name, COALESCE(description,''), is_active, start_on, end_on, is_coupon_required,
			rule_to_apply, discount_amount, max_discount_amount, discount_step, usage_limit_per_coupon,
			usage_limit_per_customer, created_at, updated_at
		FROM pricing_cart_rules WHERE id = $1
	`, coupon.CartRuleID).Scan(&cr.ID, &cr.Name, &cr.Description, &cr.IsActive, &startOn, &endOn,
		&cr.IsCouponRequired, &cr.RuleToApply, &cr.DiscountAmount, &maxDiscount, &discountStep,
		&usagePerCoupon, &usagePerCustomer, &cr.CreatedAt, &cr.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart rule: %w", err)
	}
	if startOn.Valid {
		cr.StartOn = &startOn.Time
	}
	if endOn.Valid {
		cr.EndOn = &endOn.Time
	}
	if maxDiscount.Valid {
		cr.MaxDiscountAmount = &maxDiscount.Float64
	}
	if discountStep.Valid {
		v := int(discountStep.Int64)
		cr.DiscountStep = &v
	}
	if usagePerCoupon.Valid {
		v := int(usagePerCoupon.Int64)
		cr.UsageLimitPerCoupon = &v
	}
	if usagePerCustomer.Valid {
		v := int(usagePerCustomer.Int64)
		cr.UsageLimitPerCustomer = &v
	}

	now := time.Now()

	if !cr.IsActive {
		return &CouponValidationResult{Valid: false, Error: "coupon is not active"}, nil
	}
	if cr.StartOn != nil && now.Before(*cr.StartOn) {
		return &CouponValidationResult{Valid: false, Error: "coupon is not yet valid"}, nil
	}
	if cr.EndOn != nil && now.After(*cr.EndOn) {
		return &CouponValidationResult{Valid: false, Error: "coupon has expired"}, nil
	}
	if cr.UsageLimitPerCoupon != nil {
		var usageCount int
		s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM pricing_cart_rule_usages WHERE coupon_id = $1`, coupon.ID).Scan(&usageCount)
		if usageCount >= *cr.UsageLimitPerCoupon {
			return &CouponValidationResult{Valid: false, Error: "coupon usage limit reached"}, nil
		}
	}

	return &CouponValidationResult{Valid: true, CartRule: &cr}, nil
}

func (s *PricingService) GetCatalogRules(ctx context.Context) ([]CatalogRuleResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, COALESCE(description,''), is_active, start_on, end_on,
			rule_to_apply, discount_amount, max_discount_amount, created_at, updated_at
		FROM pricing_catalog_rules ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to list catalog rules: %w", err)
	}
	defer rows.Close()

	var rules []CatalogRuleResponse
	for rows.Next() {
		var cr CatalogRuleResponse
		var startOn, endOn sql.NullTime
		var maxDiscount sql.NullFloat64
		if err := rows.Scan(&cr.ID, &cr.Name, &cr.Description, &cr.IsActive, &startOn, &endOn,
			&cr.RuleToApply, &cr.DiscountAmount, &maxDiscount, &cr.CreatedAt, &cr.UpdatedAt); err != nil {
			return nil, err
		}
		if startOn.Valid {
			cr.StartOn = &startOn.Time
		}
		if endOn.Valid {
			cr.EndOn = &endOn.Time
		}
		if maxDiscount.Valid {
			cr.MaxDiscountAmount = &maxDiscount.Float64
		}
		rules = append(rules, cr)
	}
	if rules == nil {
		rules = []CatalogRuleResponse{}
	}
	return rules, nil
}

func (s *PricingService) CreateCatalogRule(ctx context.Context, req CreateCatalogRuleRequest) (*CatalogRuleResponse, error) {
	var cr CatalogRuleResponse
	var startOn, endOn sql.NullTime
	var maxDiscount sql.NullFloat64
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO pricing_catalog_rules (name, description, is_active, start_on, end_on,
			rule_to_apply, discount_amount, max_discount_amount, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW(),NOW())
		RETURNING id, name, description, is_active, start_on, end_on,
			rule_to_apply, discount_amount, max_discount_amount, created_at, updated_at
	`, req.Name, req.Description, req.IsActive, req.StartOn, req.EndOn,
		req.RuleToApply, req.DiscountAmount, req.MaxDiscountAmount).Scan(
		&cr.ID, &cr.Name, &cr.Description, &cr.IsActive, &startOn, &endOn,
		&cr.RuleToApply, &cr.DiscountAmount, &maxDiscount, &cr.CreatedAt, &cr.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create catalog rule: %w", err)
	}
	if startOn.Valid {
		cr.StartOn = &startOn.Time
	}
	if endOn.Valid {
		cr.EndOn = &endOn.Time
	}
	if maxDiscount.Valid {
		cr.MaxDiscountAmount = &maxDiscount.Float64
	}
	cr.CustomerGroups = []uint{}
	return &cr, nil
}

func (s *PricingService) UpdateCatalogRule(ctx context.Context, id uint, req CreateCatalogRuleRequest) (*CatalogRuleResponse, error) {
	var cr CatalogRuleResponse
	var startOn, endOn sql.NullTime
	var maxDiscount sql.NullFloat64
	err := s.db.QueryRowContext(ctx, `
		UPDATE pricing_catalog_rules SET name=$1, description=$2, is_active=$3, start_on=$4, end_on=$5,
			rule_to_apply=$6, discount_amount=$7, max_discount_amount=$8, updated_at=NOW()
		WHERE id=$9
		RETURNING id, name, description, is_active, start_on, end_on,
			rule_to_apply, discount_amount, max_discount_amount, created_at, updated_at
	`, req.Name, req.Description, req.IsActive, req.StartOn, req.EndOn,
		req.RuleToApply, req.DiscountAmount, req.MaxDiscountAmount, id).Scan(
		&cr.ID, &cr.Name, &cr.Description, &cr.IsActive, &startOn, &endOn,
		&cr.RuleToApply, &cr.DiscountAmount, &maxDiscount, &cr.CreatedAt, &cr.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("catalog rule not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update catalog rule: %w", err)
	}
	if startOn.Valid {
		cr.StartOn = &startOn.Time
	}
	if endOn.Valid {
		cr.EndOn = &endOn.Time
	}
	if maxDiscount.Valid {
		cr.MaxDiscountAmount = &maxDiscount.Float64
	}
	return &cr, nil
}

func (s *PricingService) DeleteCatalogRule(ctx context.Context, id uint) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM pricing_catalog_rules WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete catalog rule: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("catalog rule not found")
	}
	return nil
}

func (s *PricingService) getCoupons(ctx context.Context, cartRuleID uint) ([]CouponResponse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, cart_rule_id, code, created_at FROM pricing_coupons WHERE cart_rule_id = $1
	`, cartRuleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coupons []CouponResponse
	for rows.Next() {
		var c CouponResponse
		if err := rows.Scan(&c.ID, &c.CartRuleID, &c.Code, &c.CreatedAt); err != nil {
			return nil, err
		}
		coupons = append(coupons, c)
	}
	if coupons == nil {
		coupons = []CouponResponse{}
	}
	return coupons, nil
}

func (s *PricingService) getCartRuleCategories(ctx context.Context, cartRuleID uint) ([]uint, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT category_id FROM pricing_cart_rule_categories WHERE cart_rule_id = $1
	`, cartRuleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uint
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if ids == nil {
		ids = []uint{}
	}
	return ids, nil
}

func (s *PricingService) getCartRuleProducts(ctx context.Context, cartRuleID uint) ([]uint, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT product_id FROM pricing_cart_rule_products WHERE cart_rule_id = $1
	`, cartRuleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uint
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if ids == nil {
		ids = []uint{}
	}
	return ids, nil
}

func (s *PricingService) getCartRuleCustomerGroups(ctx context.Context, cartRuleID uint) ([]uint, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT customer_group_id FROM pricing_cart_rule_customer_groups WHERE cart_rule_id = $1
	`, cartRuleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uint
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if ids == nil {
		ids = []uint{}
	}
	return ids, nil
}

type CouponValidationResult struct {
	Valid    bool           `json:"valid"`
	CartRule *CartRuleResponse `json:"cartRule,omitempty"`
	Error    string         `json:"error,omitempty"`
}

type CreateCartRuleRequest struct {
	Name                 string     `json:"name"`
	Description          string     `json:"description"`
	IsActive             bool       `json:"isActive"`
	StartOn              *time.Time `json:"startOn,omitempty"`
	EndOn                *time.Time `json:"endOn,omitempty"`
	IsCouponRequired     bool       `json:"isCouponRequired"`
	RuleToApply          string     `json:"ruleToApply"`
	DiscountAmount       float64    `json:"discountAmount"`
	MaxDiscountAmount    *float64   `json:"maxDiscountAmount,omitempty"`
	DiscountStep         *int       `json:"discountStep,omitempty"`
	UsageLimitPerCoupon  *int       `json:"usageLimitPerCoupon,omitempty"`
	UsageLimitPerCustomer *int      `json:"usageLimitPerCustomer,omitempty"`
}

type CreateCatalogRuleRequest struct {
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	IsActive          bool       `json:"isActive"`
	StartOn           *time.Time `json:"startOn,omitempty"`
	EndOn             *time.Time `json:"endOn,omitempty"`
	RuleToApply       string     `json:"ruleToApply"`
	DiscountAmount    float64    `json:"discountAmount"`
	MaxDiscountAmount *float64   `json:"maxDiscountAmount,omitempty"`
}
