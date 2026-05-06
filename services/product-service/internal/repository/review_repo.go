package repository

import (
	"context"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewRepo struct {
	db *pgxpool.Pool
}

func NewReviewRepo(db *pgxpool.Pool) *ReviewRepo {
	return &ReviewRepo{db: db}
}

func (r *ReviewRepo) Create(ctx context.Context, rev *models.Review) error {
	query := `INSERT INTO reviews (id, product_id, user_id, order_id, rating, title, body, status)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query, rev.ID, rev.ProductID, rev.UserID, rev.OrderID, rev.Rating, rev.Title, rev.Body, rev.Status)
	return err
}

func (r *ReviewRepo) ListByProduct(ctx context.Context, productID uuid.UUID, page, pageSize int32) ([]models.Review, int32, error) {
	var total int32
	_ = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM reviews WHERE product_id = $1 AND status = 'approved'`, productID).Scan(&total)

	query := `SELECT id, product_id, user_id, order_id, rating, title, body, status, created_at
	          FROM reviews WHERE product_id = $1 AND status = 'approved'
	          ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, query, productID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var rev models.Review
		err := rows.Scan(&rev.ID, &rev.ProductID, &rev.UserID, &rev.OrderID, &rev.Rating, &rev.Title, &rev.Body, &rev.Status, &rev.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, total, rows.Err()
}

func (r *ReviewRepo) GetAverageRating(ctx context.Context, productID uuid.UUID) (float32, int, error) {
	query := `SELECT COALESCE(AVG(rating), 0), COUNT(*) FROM reviews WHERE product_id = $1 AND status = 'approved'`
	var avg float32
	var count int
	err := r.db.QueryRow(ctx, query, productID).Scan(&avg, &count)
	return avg, count, err
}
