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
	invRepo  *repository.InventoryRepo
}

func NewProductService(prodRepo *repository.ProductRepo, varRepo *repository.VariantRepo, catRepo *repository.CategoryRepo, invRepo *repository.InventoryRepo) *ProductService {
	return &ProductService{prodRepo: prodRepo, varRepo: varRepo, catRepo: catRepo, invRepo: invRepo}
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

		// Create default inventory
		inv := &models.Inventory{
			ID:                uuid.New(),
			VariantID:         variants[i].ID,
			QuantityOnHand:    int(variants[i].WeightGrams), // Using weight_grams as placeholder for stock in req if needed, or better, pass stock
			QuantityReserved:  0,
			QuantityAvailable: 0,
			ReorderPoint:      10,
		}
		if err := s.invRepo.Create(ctx, inv); err != nil {
			return fmt.Errorf("create inventory: %w", err)
		}
	}

	return nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, p *models.Product, variants []models.Variant) error {
	p.UpdatedAt = time.Now().UTC()
	if err := s.prodRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("update product: %w", err)
	}

	// Simple logic: if variants are provided, we should ideally sync them.
	// For now, if only one variant is provided, we update it or create it if missing.
	for i := range variants {
		if variants[i].ID == uuid.Nil {
			// Try to find existing variant for this product
			existing, err := s.varRepo.GetByProductID(ctx, p.ID)
			if err == nil && len(existing) > 0 {
				variants[i].ID = existing[0].ID
				variants[i].ProductID = p.ID
				variants[i].UpdatedAt = time.Now().UTC()
				if err := s.varRepo.Update(ctx, &variants[i]); err != nil {
					return fmt.Errorf("update variant: %w", err)
				}
				// Update inventory
				inv, err := s.invRepo.GetByVariantID(ctx, variants[i].ID)
				if err == nil {
					inv.QuantityOnHand = int(variants[i].WeightGrams) // Using weight_grams as stock placeholder
					if err := s.invRepo.Update(ctx, inv); err != nil {
						return fmt.Errorf("update inventory: %w", err)
					}
				}
			} else {
				// Create new variant and inventory
				variants[i].ID = uuid.New()
				variants[i].ProductID = p.ID
				variants[i].CreatedAt = time.Now().UTC()
				if err := s.varRepo.Create(ctx, &variants[i]); err != nil {
					return fmt.Errorf("create variant: %w", err)
				}
				inv := &models.Inventory{
					ID:             uuid.New(),
					VariantID:      variants[i].ID,
					QuantityOnHand: int(variants[i].WeightGrams),
				}
				s.invRepo.Create(ctx, inv)
			}
		} else {
			variants[i].UpdatedAt = time.Now().UTC()
			if err := s.varRepo.Update(ctx, &variants[i]); err != nil {
				return fmt.Errorf("update variant: %w", err)
			}
			// Update inventory
			inv, err := s.invRepo.GetByVariantID(ctx, variants[i].ID)
			if err == nil {
				inv.QuantityOnHand = int(variants[i].WeightGrams)
				s.invRepo.Update(ctx, inv)
			}
		}
	}

	return nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return s.prodRepo.Delete(ctx, id)
}

func (s *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, []models.Variant, error) {
	product, err := s.prodRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	variants, err := s.varRepo.GetByProductID(ctx, id)
	if err != nil {
		// It's okay if there are no variants, though unusual
		return product, nil, nil
	}

	return product, variants, nil
}

func (s *ProductService) GetProductBySlug(ctx context.Context, slug string) (*models.Product, []models.Variant, error) {
	// This would need a slug lookup in repo — simplified
	return nil, nil, fmt.Errorf("not implemented")
}

func (s *ProductService) ListProducts(ctx context.Context, categoryID uuid.UUID, page, pageSize int32) ([]models.Product, int32, error) {
	return s.prodRepo.List(ctx, categoryID, page, pageSize)
}

func (s *ProductService) GetVariantsBatch(ctx context.Context, ids []uuid.UUID) ([]models.Variant, error) {
	return s.varRepo.GetByIDs(ctx, ids)
}

func generateSlug(name string) string {
	// Simplified — use github.com/gosimple/slug in production
	return fmt.Sprintf("%s-%s", name, uuid.New().String()[:8])
}
