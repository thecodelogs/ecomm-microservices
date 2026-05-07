package service

import (
	"context"
	"fmt"
	"time"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/repository"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/models"

	"github.com/google/uuid"
)

type ProductService struct {
	prodRepo *repository.ProductRepo
	varRepo  *repository.VariantRepo
	catRepo  *repository.CategoryRepo
}

func NewProductService(prodRepo *repository.ProductRepo, varRepo *repository.VariantRepo, catRepo *repository.CategoryRepo) *ProductService {
	return &ProductService{prodRepo: prodRepo, varRepo: varRepo, catRepo: catRepo}
}

func (s *ProductService) CreateProduct(ctx context.Context, p *models.Product, variants []models.Variant) error {
	p.ID = uuid.New()
	p.Slug = generateSlug(p.Name)
	p.Status = "draft"
	p.CreatedAt = time.Now().UTC()
	p.UpdatedAt = time.Now().UTC()

	if err := s.prodRepo.Create(ctx, p); err != nil {
		return fmt.Errorf("create product: %w", err)
	}

	for i := range variants {
		variants[i].ID = uuid.New()
		variants[i].ProductID = p.ID
		variants[i].CreatedAt = time.Now().UTC()
		if err := s.varRepo.Create(ctx, &variants[i]); err != nil {
			return fmt.Errorf("create variant: %w", err)
		}
	}

	return nil
}

func (s *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, []models.Variant, error) {
	product, err := s.prodRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	variants, err := s.varRepo.GetByProductID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	return product, variants, nil
}

func (s *ProductService) GetProductBySlug(ctx context.Context, slug string) (*models.Product, []models.Variant, error) {
	// This would need a slug lookup in repo — simplified
	return nil, nil, fmt.Errorf("not implemented")
}

func (s *ProductService) ListProducts(ctx context.Context, categoryID uuid.UUID, page, pageSize int32) ([]models.Product, int32, error) {
	fmt.Println("adslkjhASLKDJHAKSJLD=================", categoryID)
	return s.prodRepo.ListByCategory(ctx, categoryID, page, pageSize)
}

func (s *ProductService) GetVariantsBatch(ctx context.Context, ids []uuid.UUID) ([]models.Variant, error) {
	return s.varRepo.GetByIDs(ctx, ids)
}

func generateSlug(name string) string {
	// Simplified — use github.com/gosimple/slug in production
	return fmt.Sprintf("%s-%s", name, uuid.New().String()[:8])
}
