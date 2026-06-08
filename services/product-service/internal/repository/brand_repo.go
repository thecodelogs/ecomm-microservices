package repository

import (
	"context"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BrandRepo struct {
	db *pgxpool.Pool
}

func NewBrandRepo(db *pgxpool.Pool) *BrandRepo {
	return &BrandRepo{db: db}
}

func (r *BrandRepo) Create(ctx context.Context, b *models.Brand) error {
	query := `INSERT INTO brands (id, name, description, image_url, is_active)
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, b.ID, b.Name, b.Description, b.ImageURL, b.IsActive)
	return err
}

func (r *BrandRepo) GetByID(ctx context.Context, id string) (*models.Brand, error) {
	query := `SELECT id, name, description, image_url, is_active, created_at, updated_at
	          FROM brands WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var b models.Brand
	err := row.Scan(&b.ID, &b.Name, &b.Description, &b.ImageURL, &b.IsActive, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BrandRepo) Update(ctx context.Context, b *models.Brand) error {
	query := `UPDATE brands 
	          SET name = $1, description = $2, image_url = $3, is_active = $4, updated_at = NOW()
	          WHERE id = $5`
	_, err := r.db.Exec(ctx, query, b.Name, b.Description, b.ImageURL, b.IsActive, b.ID)
	return err
}

func (r *BrandRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE brands SET is_active = false, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *BrandRepo) BrandsList(ctx context.Context, page, pageSize int32) ([]models.Brand, int32, error) {
	countQuery := `SELECT COUNT(*) FROM brands WHERE is_active = true`

	var total int32
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, name, description, image_url, is_active, created_at, updated_at
		FROM brands
		WHERE is_active = true
		ORDER BY name
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var brands []models.Brand
	for rows.Next() {
		var b models.Brand
		err := rows.Scan(
			&b.ID,
			&b.Name,
			&b.Description,
			&b.ImageURL,
			&b.IsActive,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		brands = append(brands, b)
	}

	return brands, total, rows.Err()
}
