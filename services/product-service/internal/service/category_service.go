package service

import (
	"context"
	"fmt"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/repository"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/models"

	"github.com/google/uuid"
)

type CategoryService struct {
	prodRepo *repository.ProductRepo
	varRepo  *repository.VariantRepo
	catRepo  *repository.CategoryRepo
}

func NewCategoryService(prodRepo *repository.ProductRepo, varRepo *repository.VariantRepo, catRepo *repository.CategoryRepo) *CategoryService {
	return &CategoryService{prodRepo: prodRepo, varRepo: varRepo, catRepo: catRepo}
}

func (c *CategoryService) CreateCategory(ctx context.Context, p *models.Category) error {
	p.ID = uuid.New()
	p.Slug = generateSlug(p.Name)

	if err := c.catRepo.Create(ctx, p); err != nil {
		return fmt.Errorf("create category: %w", err)
	}

	return nil
}

func (c *CategoryService) UpdateCategory(ctx context.Context, p *models.Category) error {
	if err := c.catRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("update category: %w", err)
	}

	return nil
}

func (c *CategoryService) GetCategory(ctx context.Context, id string) (*models.Category, error) {
	return c.catRepo.GetByID(ctx, id)
}

func (c *CategoryService) ListCategories(ctx context.Context, page, pageSize int32) ([]models.Category, int32, error) {
	return c.catRepo.CategoriesList(ctx, page, pageSize)
}
