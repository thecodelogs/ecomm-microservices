package repository

import (
	"context"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(ctx context.Context, p *models.Product) error {
	query := `INSERT INTO products (id, category_id, slug, name, description, short_description, brand, tags, attributes, status, vendor_id)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.Exec(ctx, query, p.ID, p.CategoryID, p.Slug, p.Name, p.Description, p.ShortDescription, p.Brand, p.Tags, p.Attributes, p.Status, p.VendorID)
	return err
}

func (r *ProductRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	query := `SELECT id, category_id, slug, name, description, short_description, brand, tags, attributes, status, vendor_id, avg_rating, review_count, created_at, updated_at
	          FROM products WHERE id = $1 AND status = 'active'`
	row := r.db.QueryRow(ctx, query, id)
	return r.scanProduct(row)
}

func (r *ProductRepo) GetBySlug(ctx context.Context, slug string) (*models.Product, error) {
	query := `SELECT id, category_id, slug, name, description, short_description, brand, tags, attributes, status, vendor_id, avg_rating, review_count, created_at, updated_at
	          FROM products WHERE slug = $1 AND status = 'active'`
	row := r.db.QueryRow(ctx, query, slug)
	return r.scanProduct(row)
}

func (r *ProductRepo) ListByCategory(ctx context.Context, categoryID uuid.UUID, page, pageSize int32) ([]models.Product, int32, error) {
	countQuery := `SELECT COUNT(*) FROM products WHERE category_id = $1 AND status = 'active'`
	var total int32
	_ = r.db.QueryRow(ctx, countQuery, categoryID).Scan(&total)

	query := `SELECT id, category_id, slug, name, description, short_description, brand, tags, attributes, status, vendor_id, avg_rating, review_count, created_at, updated_at
	          FROM products WHERE category_id = $1 AND status = 'active'
	          ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, query, categoryID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		p, err := r.scanProduct(rows)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, *p)
	}
	return products, total, rows.Err()
}

func (r *ProductRepo) UpdateRating(ctx context.Context, productID uuid.UUID, avgRating float32, reviewCount int) error {
	query := `UPDATE products SET avg_rating = $1, review_count = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.Exec(ctx, query, avgRating, reviewCount, productID)
	return err
}

func (r *ProductRepo) scanProduct(row pgx.Row) (*models.Product, error) {
	var p models.Product
	err := row.Scan(&p.ID, &p.CategoryID, &p.Slug, &p.Name, &p.Description, &p.ShortDescription, &p.Brand, &p.Tags, &p.Attributes, &p.Status, &p.VendorID, &p.AvgRating, &p.ReviewCount, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
