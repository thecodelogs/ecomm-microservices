package service

import (
	"context"
	"time"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/models"
	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/repository"

	"github.com/google/uuid"
)

type ReviewService struct {
	revRepo  *repository.ReviewRepo
	prodRepo *repository.ProductRepo
}

func NewReviewService(revRepo *repository.ReviewRepo, prodRepo *repository.ProductRepo) *ReviewService {
	return &ReviewService{revRepo: revRepo, prodRepo: prodRepo}
}

func (s *ReviewService) CreateReview(ctx context.Context, productID, userID, orderID uuid.UUID, rating int16, title, body string) error {
	rev := &models.Review{
		ID:        uuid.New(),
		ProductID: productID,
		UserID:    userID,
		OrderID:   orderID,
		Rating:    rating,
		Title:     title,
		Body:      body,
		Status:    "pending", // requires moderation
		CreatedAt: time.Now().UTC(),
	}

	if err := s.revRepo.Create(ctx, rev); err != nil {
		return err
	}

	// Update product rating
	avg, count, _ := s.revRepo.GetAverageRating(ctx, productID)
	return s.prodRepo.UpdateRating(ctx, productID, avg, count)
}

func (s *ReviewService) ListReviews(ctx context.Context, productID uuid.UUID, page, pageSize int32) ([]models.Review, int32, error) {
	return s.revRepo.ListByProduct(ctx, productID, page, pageSize)
}
