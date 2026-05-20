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

func (r *ProductRepo) Update(ctx context.Context, p *models.Product) error {
	query := `UPDATE products SET 
				category_id = $1, 
				slug = $2, 
				name = $3, 
				description = $4, 
				short_description = $5, 
				brand = $6, 
				tags = $7, 
				attributes = $8, 
				status = $9, 
				updated_at = NOW() 
			  WHERE id = $10`
	_, err := r.db.Exec(ctx, query, p.CategoryID, p.Slug, p.Name, p.Description, p.ShortDescription, p.Brand, p.Tags, p.Attributes, p.Status, p.ID)
	return err
}

func (r *ProductRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *ProductRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	query := `SELECT id, category_id, slug, name, description, short_description, brand, tags, attributes, status, vendor_id, avg_rating, review_count, created_at, updated_at
	          FROM products WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	return r.scanProduct(row)
}

func (r *ProductRepo) GetBySlug(ctx context.Context, slug string) (*models.Product, error) {
	query := `SELECT id, category_id, slug, name, description, short_description, brand, tags, attributes, status, vendor_id, avg_rating, review_count, created_at, updated_at
	          FROM products WHERE slug = $1 AND status = 'active'`
	row := r.db.QueryRow(ctx, query, slug)
	return r.scanProduct(row)
}

func (r *ProductRepo) List(ctx context.Context, categoryID uuid.UUID, page, pageSize int32) ([]models.Product, int32, error) {
	countWhere := "WHERE status = 'active'"
	queryWhere := "WHERE status = 'active'"
	
	var countArgs []interface{}
	queryArgs := []interface{}{pageSize, (page-1)*pageSize}
	
	if categoryID != uuid.Nil {
		countWhere += " AND category_id = $1"
		countArgs = append(countArgs, categoryID)
		
		queryWhere += " AND category_id = $3"
		queryArgs = append(queryArgs, categoryID)
	}

	countQuery := `SELECT COUNT(*) FROM products ` + countWhere
	var total int32
	err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, category_id, slug, name, description, short_description, brand, tags, attributes, status, vendor_id, avg_rating, review_count, created_at, updated_at
	          FROM products ` + queryWhere + `
	          ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	
	rows, err := r.db.Query(ctx, query, queryArgs...)
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
