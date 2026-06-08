package repository

import (
	"context"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VariantImageRepo struct {
	db *pgxpool.Pool
}

func NewVariantImageRepo(db *pgxpool.Pool) *VariantImageRepo {
	return &VariantImageRepo{db: db}
}

func (r *VariantImageRepo) Create(ctx context.Context, img *models.VariantImage) error {
	db := getDb(ctx, r.db)
	query := `INSERT INTO variant_images (id, variant_id, url, alt_text, sort_order)
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(ctx, query, img.ID, img.VariantID, img.URL, img.AltText, img.SortOrder)
	return err
}

func (r *VariantImageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	db := getDb(ctx, r.db)
	query := `DELETE FROM variant_images WHERE id = $1`
	_, err := db.Exec(ctx, query, id)
	return err
}

func (r *VariantImageRepo) DeleteByVariantID(ctx context.Context, variantID uuid.UUID) error {
	db := getDb(ctx, r.db)
	query := `DELETE FROM variant_images WHERE variant_id = $1`
	_, err := db.Exec(ctx, query, variantID)
	return err
}

func (r *VariantImageRepo) GetByVariantID(ctx context.Context, variantID uuid.UUID) ([]models.VariantImage, error) {
	db := getDb(ctx, r.db)
	query := `SELECT id, variant_id, url, alt_text, sort_order, created_at
	          FROM variant_images WHERE variant_id = $1 ORDER BY sort_order ASC`
	rows, err := db.Query(ctx, query, variantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.VariantImage
	for rows.Next() {
		var img models.VariantImage
		err := rows.Scan(&img.ID, &img.VariantID, &img.URL, &img.AltText, &img.SortOrder, &img.CreatedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, rows.Err()
}

func (r *VariantImageRepo) UpdateSortOrder(ctx context.Context, id uuid.UUID, sortOrder int32) error {
	db := getDb(ctx, r.db)
	query := `UPDATE variant_images SET sort_order = $1 WHERE id = $2`
	_, err := db.Exec(ctx, query, sortOrder, id)
	return err
}
