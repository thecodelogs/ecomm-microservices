package repository

import (
	"context"
	"fmt"
	"strings"

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
	query := `INSERT INTO products (id, category_id, slug, name, description, short_description, brand, brand_id, tags, attributes, status, vendor_id)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.Exec(ctx, query, p.ID, p.CategoryID, p.Slug, p.Name, p.Description, p.ShortDescription, p.Brand, p.BrandID, p.Tags, p.Attributes, p.Status, p.VendorID)
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
				brand_id = $7,
				tags = $8, 
				attributes = $9, 
				status = $10, 
				updated_at = NOW() 
			  WHERE id = $11`
	_, err := r.db.Exec(ctx, query, p.CategoryID, p.Slug, p.Name, p.Description, p.ShortDescription, p.Brand, p.BrandID, p.Tags, p.Attributes, p.Status, p.ID)
	return err
}

func (r *ProductRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE products SET status = 'deleted', updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *ProductRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	query := `SELECT id, category_id, slug, name, description, short_description, brand, brand_id, tags, attributes, status, vendor_id, avg_rating, review_count, created_at, updated_at
	          FROM products WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	return r.scanProduct(row)
}

func (r *ProductRepo) GetBySlug(ctx context.Context, slug string) (*models.Product, error) {
	query := `SELECT id, category_id, slug, name, description, short_description, brand, brand_id, tags, attributes, status, vendor_id, avg_rating, review_count, created_at, updated_at
	          FROM products WHERE slug = $1 AND status = 'active'`
	row := r.db.QueryRow(ctx, query, slug)
	return r.scanProduct(row)
}

func (r *ProductRepo) List(ctx context.Context, categoryID uuid.UUID, page, pageSize int32, isAdmin bool, minPrice, maxPrice float64, brands []string, sortField, sortDir string) ([]models.Product, int32, error) {
	whereClauses := []string{"1=1"}
	var args []interface{}
	argID := 1

	if !isAdmin {
		whereClauses = append(whereClauses, fmt.Sprintf("p.status = $%d", argID))
		args = append(args, "active")
		argID++
	}

	if categoryID != uuid.Nil {
		whereClauses = append(whereClauses, fmt.Sprintf("p.category_id = $%d", argID))
		args = append(args, categoryID)
		argID++
	}

	if len(brands) > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("p.brand = ANY($%d)", argID))
		args = append(args, brands)
		argID++
	}

	if minPrice > 0 || maxPrice > 0 {
		priceClause := fmt.Sprintf("EXISTS(SELECT 1 FROM variants v WHERE v.product_id = p.id")
		if minPrice > 0 {
			priceClause += fmt.Sprintf(" AND v.price >= $%d", argID)
			args = append(args, int64(minPrice*100))
			argID++
		}
		if maxPrice > 0 {
			priceClause += fmt.Sprintf(" AND v.price <= $%d", argID)
			args = append(args, int64(maxPrice*100))
			argID++
		}
		priceClause += ")"
		whereClauses = append(whereClauses, priceClause)
	}

	whereSQL := "WHERE " + strings.Join(whereClauses, " AND ")

	countQuery := `SELECT COUNT(*) FROM products p ` + whereSQL
	var total int32
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	orderBy := "ORDER BY p.updated_at DESC"
	dir := "DESC"
	if strings.ToUpper(sortDir) == "ASC" {
		dir = "ASC"
	}
	
	switch strings.ToUpper(sortField) {
	case "PRICE":
		orderBy = fmt.Sprintf("ORDER BY (SELECT MIN(price) FROM variants v WHERE v.product_id = p.id) %s", dir)
	case "CREATED_AT":
		orderBy = fmt.Sprintf("ORDER BY p.created_at %s", dir)
	case "RATING":
		orderBy = fmt.Sprintf("ORDER BY p.avg_rating %s", dir)
	}

	query := fmt.Sprintf(`SELECT p.id, p.category_id, p.slug, p.name, p.description, p.short_description, p.brand, p.brand_id, p.tags, p.attributes, p.status, p.vendor_id, p.avg_rating, p.review_count, p.created_at, p.updated_at
	          FROM products p %s %s LIMIT $%d OFFSET $%d`, whereSQL, orderBy, argID, argID+1)
	
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.Query(ctx, query, args...)
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

func (r *ProductRepo) Search(ctx context.Context, queryStr string, page, pageSize int32, isAdmin bool) ([]models.Product, int32, error) {
	whereClauses := []string{"1=1"}
	var args []interface{}
	argID := 1

	if !isAdmin {
		whereClauses = append(whereClauses, fmt.Sprintf("p.status = $%d", argID))
		args = append(args, "active")
		argID++
	}

	if queryStr != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(p.name ILIKE $%d OR p.description ILIKE $%d)", argID, argID))
		searchPattern := "%" + queryStr + "%"
		args = append(args, searchPattern)
		argID++
	}

	whereSQL := "WHERE " + strings.Join(whereClauses, " AND ")

	countQuery := `SELECT COUNT(*) FROM products p ` + whereSQL
	var total int32
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`SELECT p.id, p.category_id, p.slug, p.name, p.description, p.short_description, p.brand, p.brand_id, p.tags, p.attributes, p.status, p.vendor_id, p.avg_rating, p.review_count, p.created_at, p.updated_at
	          FROM products p %s ORDER BY p.updated_at DESC LIMIT $%d OFFSET $%d`, whereSQL, argID, argID+1)
	
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.Query(ctx, query, args...)
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
	err := row.Scan(&p.ID, &p.CategoryID, &p.Slug, &p.Name, &p.Description, &p.ShortDescription, &p.Brand, &p.BrandID, &p.Tags, &p.Attributes, &p.Status, &p.VendorID, &p.AvgRating, &p.ReviewCount, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
