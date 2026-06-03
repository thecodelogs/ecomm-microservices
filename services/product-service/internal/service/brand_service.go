package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/models"
	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/repository"
)

type BrandService struct {
	brandRepo *repository.BrandRepo
}

func NewBrandService(brandRepo *repository.BrandRepo) *BrandService {
	return &BrandService{brandRepo: brandRepo}
}

func (s *BrandService) CreateBrand(ctx context.Context, brand *models.Brand) error {
	brand.ID = uuid.New()
	return s.brandRepo.Create(ctx, brand)
}

func (s *BrandService) UpdateBrand(ctx context.Context, brand *models.Brand) error {
	return s.brandRepo.Update(ctx, brand)
}

func (s *BrandService) DeleteBrand(ctx context.Context, id string) error {
	return s.brandRepo.Delete(ctx, id)
}

func (s *BrandService) GetBrand(ctx context.Context, id string) (*models.Brand, error) {
	return s.brandRepo.GetByID(ctx, id)
}

func (s *BrandService) ListBrands(ctx context.Context, page, pageSize int32) ([]models.Brand, int32, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.brandRepo.BrandsList(ctx, page, pageSize)
}
