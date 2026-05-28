package repository

import (
	"context"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VariantRepo struct {
	db *pgxpool.Pool
}

func NewVariantRepo(db *pgxpool.Pool) *VariantRepo {
	return &VariantRepo{db: db}
}

func (r *VariantRepo) Create(ctx context.Context, v *models.Variant) error {
	query := `INSERT INTO variants (id, product_id, sku, name, options, price, compare_at_price, cost_price, weight_grams, image_url, is_active)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.Exec(ctx, query, v.ID, v.ProductID, v.SKU, v.Name, v.Options, v.Price, v.CompareAtPrice, v.CostPrice, v.WeightGrams, v.ImageURL, v.IsActive)
	return err
}

func (r *VariantRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Variant, error) {
	query := `SELECT id, product_id, sku, name, options, price, compare_at_price, cost_price, weight_grams, image_url, is_active, created_at, updated_at
	          FROM variants WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	return r.scanVariant(row)
}

func (r *VariantRepo) GetByProductID(ctx context.Context, productID uuid.UUID) ([]models.Variant, error) {
	query := `SELECT id, product_id, sku, name, options, price, compare_at_price, cost_price, weight_grams, image_url, is_active, created_at, updated_at
	          FROM variants WHERE product_id = $1 AND is_active = true`
	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants []models.Variant
	for rows.Next() {
		v, err := r.scanVariant(rows)
		if err != nil {
			return nil, err
		}
		variants = append(variants, *v)
	}
	return variants, rows.Err()
}

func (r *VariantRepo) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Variant, error) {
	query := `SELECT id, product_id, sku, name, options, price, compare_at_price, cost_price, weight_grams, image_url, is_active, created_at, updated_at
	          FROM variants WHERE id = ANY($1)`
	rows, err := r.db.Query(ctx, query, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variants []models.Variant
	for rows.Next() {
		v, err := r.scanVariant(rows)
		if err != nil {
			return nil, err
		}
		variants = append(variants, *v)
	}
	return variants, rows.Err()
}

func (r *VariantRepo) Update(ctx context.Context, v *models.Variant) error {
	query := `UPDATE variants SET 
				sku = $1, 
				name = $2, 
				options = $3, 
				price = $4, 
				compare_at_price = $5, 
				cost_price = $6, 
				weight_grams = $7, 
				image_url = $8, 
				is_active = $9, 
				updated_at = NOW() 
			  WHERE id = $10`
	_, err := r.db.Exec(ctx, query, v.SKU, v.Name, v.Options, v.Price, v.CompareAtPrice, v.CostPrice, v.WeightGrams, v.ImageURL, v.IsActive, v.ID)
	return err
}

func (r *VariantRepo) DeleteByProductID(ctx context.Context, productID uuid.UUID) error {
	query := `UPDATE variants SET is_active = false, updated_at = NOW() WHERE product_id = $1`
	_, err := r.db.Exec(ctx, query, productID)
	return err
}

func (r *VariantRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE variants SET is_active = false, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *VariantRepo) scanVariant(row pgx.Row) (*models.Variant, error) {
	var v models.Variant
	err := row.Scan(&v.ID, &v.ProductID, &v.SKU, &v.Name, &v.Options, &v.Price, &v.CompareAtPrice, &v.CostPrice, &v.WeightGrams, &v.ImageURL, &v.IsActive, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &v, nil
}
