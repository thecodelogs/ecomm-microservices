package repository

import (
	"context"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) Create(ctx context.Context, cat *models.Category) error {
	query := `INSERT INTO categories (id, parent_id, slug, name, description, image_url, sort_order, is_active)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query, cat.ID, cat.ParentID, cat.Slug, cat.Name, cat.Description, cat.ImageURL, cat.SortOrder, cat.IsActive)
	return err
}

func (r *CategoryRepo) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	query := `SELECT id, parent_id, slug, name, description, image_url, sort_order, is_active, created_at
	          FROM categories WHERE slug = $1 AND is_active = true`
	row := r.db.QueryRow(ctx, query, slug)
	var c models.Category
	err := row.Scan(&c.ID, &c.ParentID, &c.Slug, &c.Name, &c.Description, &c.ImageURL, &c.SortOrder, &c.IsActive, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepo) CategoriesList(ctx context.Context, page, pageSize int32) ([]models.Category, int32, error) {
	countQuery := `SELECT COUNT(*) FROM categories WHERE is_active = true`

	var total int32
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, parent_id, slug, name, description, image_url, sort_order, is_active, created_at
		FROM categories
		WHERE is_active = true
		ORDER BY sort_order, name
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var cats []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(
			&c.ID,
			&c.ParentID,
			&c.Slug,
			&c.Name,
			&c.Description,
			&c.ImageURL,
			&c.SortOrder,
			&c.IsActive,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		cats = append(cats, c)
	}

	return cats, total, rows.Err()
}
