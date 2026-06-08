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
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO products (id, slug, name, description, short_description, brand, brand_id, tags, attributes, status, vendor_id)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err = tx.Exec(ctx, query, p.ID, p.Slug, p.Name, p.Description, p.ShortDescription, p.Brand, p.BrandID, p.Tags, p.Attributes, p.Status, p.VendorID)
	if err != nil {
		return err
	}

	if len(p.CategoryIDs) > 0 {
		for _, catID := range p.CategoryIDs {
			_, err = tx.Exec(ctx, `INSERT INTO product_categories (product_id, category_id) VALUES ($1, $2)`, p.ID, catID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *ProductRepo) Update(ctx context.Context, p *models.Product) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `UPDATE products SET 
				slug = $1, 
				name = $2, 
				description = $3, 
				short_description = $4, 
				brand = $5, 
				brand_id = $6,
				tags = $7, 
				attributes = $8, 
				status = $9, 
				updated_at = NOW() 
			  WHERE id = $10`
	_, err = tx.Exec(ctx, query, p.Slug, p.Name, p.Description, p.ShortDescription, p.Brand, p.BrandID, p.Tags, p.Attributes, p.Status, p.ID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM product_categories WHERE product_id = $1`, p.ID)
	if err != nil {
		return err
	}

	if len(p.CategoryIDs) > 0 {
		for _, catID := range p.CategoryIDs {
			_, err = tx.Exec(ctx, `INSERT INTO product_categories (product_id, category_id) VALUES ($1, $2)`, p.ID, catID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *ProductRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE products SET status = 'deleted', updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *ProductRepo) HardDelete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *ProductRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	query := `SELECT id, COALESCE((SELECT ARRAY_AGG(category_id) FROM product_categories WHERE product_id = products.id), '{}'), slug, name, description, short_description, brand, brand_id, tags, attributes, status, vendor_id, avg_rating, review_count, created_at, updated_at
	          FROM products WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	return r.scanProduct(row)
}

func (r *ProductRepo) GetBySlug(ctx context.Context, slug string) (*models.Product, error) {
	query := `SELECT id, COALESCE((SELECT ARRAY_AGG(category_id) FROM product_categories WHERE product_id = products.id), '{}'), slug, name, description, short_description, brand, brand_id, tags, attributes, status, vendor_id, avg_rating, review_count, created_at, updated_at
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
	} else {
		whereClauses = append(whereClauses, fmt.Sprintf("p.status != $%d", argID))
		args = append(args, "deleted")
		argID++
	}

	if categoryID != uuid.Nil {
		whereClauses = append(whereClauses, fmt.Sprintf("EXISTS(SELECT 1 FROM product_categories pc WHERE pc.product_id = p.id AND pc.category_id = $%d)", argID))
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

	query := fmt.Sprintf(`SELECT p.id, COALESCE((SELECT ARRAY_AGG(category_id) FROM product_categories WHERE product_id = p.id), '{}'), p.slug, p.name, p.description, p.short_description, p.brand, p.brand_id, p.tags, p.attributes, p.status, p.vendor_id, p.avg_rating, p.review_count, p.created_at, p.updated_at
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
	} else {
		whereClauses = append(whereClauses, fmt.Sprintf("p.status != $%d", argID))
		args = append(args, "deleted")
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

	query := fmt.Sprintf(`SELECT p.id, COALESCE((SELECT ARRAY_AGG(category_id) FROM product_categories WHERE product_id = p.id), '{}'), p.slug, p.name, p.description, p.short_description, p.brand, p.brand_id, p.tags, p.attributes, p.status, p.vendor_id, p.avg_rating, p.review_count, p.created_at, p.updated_at
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
	err := row.Scan(&p.ID, &p.CategoryIDs, &p.Slug, &p.Name, &p.Description, &p.ShortDescription, &p.Brand, &p.BrandID, &p.Tags, &p.Attributes, &p.Status, &p.VendorID, &p.AvgRating, &p.ReviewCount, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
