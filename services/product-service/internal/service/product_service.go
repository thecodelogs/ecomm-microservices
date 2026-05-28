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
	imgRepo  *repository.VariantImageRepo
}

func NewProductService(prodRepo *repository.ProductRepo, varRepo *repository.VariantRepo, catRepo *repository.CategoryRepo, invRepo *repository.InventoryRepo, imgRepo *repository.VariantImageRepo) *ProductService {
	return &ProductService{prodRepo: prodRepo, varRepo: varRepo, catRepo: catRepo, invRepo: invRepo, imgRepo: imgRepo}
}

func (s *ProductService) CreateProduct(ctx context.Context, p *models.Product, variants []models.Variant) error {
	p.ID = uuid.New()
	p.Slug = generateSlug(p.Name)
	p.CreatedAt = time.Now().UTC()
	p.UpdatedAt = time.Now().UTC()

	if err := s.prodRepo.Create(ctx, p); err != nil {
		return fmt.Errorf("create product: %w", err)
	}

	for i := range variants {
		if variants[i].ID == uuid.Nil {
			variants[i].ID = uuid.New()
		}
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

		for j := range variants[i].Images {
			if variants[i].Images[j].ID == uuid.Nil {
				variants[i].Images[j].ID = uuid.New()
			}
			variants[i].Images[j].VariantID = variants[i].ID
			variants[i].Images[j].CreatedAt = time.Now().UTC()
			if err := s.imgRepo.Create(ctx, &variants[i].Images[j]); err != nil {
				return fmt.Errorf("create variant image: %w", err)
			}
		}
	}

	return nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, p *models.Product, variants []models.Variant) error {
	p.UpdatedAt = time.Now().UTC()
	if err := s.prodRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("update product: %w", err)
	}

	for i := range variants {
		if variants[i].ID == uuid.Nil {
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
		} else {
			variants[i].ProductID = p.ID
			variants[i].UpdatedAt = time.Now().UTC()
			if err := s.varRepo.Update(ctx, &variants[i]); err != nil {
				return fmt.Errorf("update variant: %w", err)
			}
			// Update inventory
			inv, err := s.invRepo.GetByVariantID(ctx, variants[i].ID)
			if err == nil {
				inv.QuantityOnHand = int(variants[i].WeightGrams)
				s.invRepo.Update(ctx, inv)
			} else {
				// If inventory didn't exist for some reason, create it
				inv = &models.Inventory{
					ID:             uuid.New(),
					VariantID:      variants[i].ID,
					QuantityOnHand: int(variants[i].WeightGrams),
				}
				s.invRepo.Create(ctx, inv)
			}
		}

		// Clear existing images for the updated variant
		if variants[i].ID != uuid.Nil {
			if err := s.imgRepo.DeleteByVariantID(ctx, variants[i].ID); err != nil {
				return fmt.Errorf("delete old variant images: %w", err)
			}
		}

		for j := range variants[i].Images {
			if variants[i].Images[j].ID == uuid.Nil {
				variants[i].Images[j].ID = uuid.New()
			}
			variants[i].Images[j].VariantID = variants[i].ID
			variants[i].Images[j].CreatedAt = time.Now().UTC()
			if err := s.imgRepo.Create(ctx, &variants[i].Images[j]); err != nil {
				return fmt.Errorf("create variant image: %w", err)
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
		return product, nil, nil
	}

	for i := range variants {
		imgs, err := s.imgRepo.GetByVariantID(ctx, variants[i].ID)
		if err == nil {
			variants[i].Images = imgs
		}
	}

	return product, variants, nil
}

func (s *ProductService) GetVariantsByProductID(ctx context.Context, productID uuid.UUID) ([]models.Variant, error) {
	variants, err := s.varRepo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}
	for i := range variants {
		imgs, err := s.imgRepo.GetByVariantID(ctx, variants[i].ID)
		if err == nil {
			variants[i].Images = imgs
		} else {
			fmt.Printf("DEBUG GetVariantsByProductID variant %s error: %v\n", variants[i].ID, err)
		}
	}
	return variants, nil
}

func (s *ProductService) GetProductBySlug(ctx context.Context, slug string) (*models.Product, []models.Variant, error) {
	// This would need a slug lookup in repo — simplified
	return nil, nil, fmt.Errorf("not implemented")
}

func (s *ProductService) ListProducts(ctx context.Context, categoryID uuid.UUID, page, pageSize int32, isAdmin bool, minPrice, maxPrice float64, brands []string, sortField, sortDir string) ([]models.Product, int32, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.prodRepo.List(ctx, categoryID, page, pageSize, isAdmin, minPrice, maxPrice, brands, sortField, sortDir)
}

func (s *ProductService) SearchProducts(ctx context.Context, query string, page, pageSize int32, isAdmin bool) ([]models.Product, int32, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.prodRepo.Search(ctx, query, page, pageSize, isAdmin)
}

func (s *ProductService) GetVariantsBatch(ctx context.Context, ids []uuid.UUID) ([]models.Variant, error) {
	variants, err := s.varRepo.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	for i := range variants {
		imgs, err := s.imgRepo.GetByVariantID(ctx, variants[i].ID)
		if err == nil {
			variants[i].Images = imgs
		}
	}
	return variants, nil
}

func generateSlug(name string) string {
	// Simplified — use github.com/gosimple/slug in production
	return fmt.Sprintf("%s-%s", name, uuid.New().String()[:8])
}

func (s *ProductService) CreateVariant(ctx context.Context, v *models.Variant) error {
	v.ID = uuid.New()
	v.CreatedAt = time.Now().UTC()
	v.UpdatedAt = time.Now().UTC()
	if err := s.varRepo.Create(ctx, v); err != nil {
		return fmt.Errorf("create variant: %w", err)
	}

	// Create default inventory
	inv := &models.Inventory{
		ID:                uuid.New(),
		VariantID:         v.ID,
		QuantityOnHand:    v.WeightGrams, // Using weight_grams as stock for now
		QuantityReserved:  0,
		QuantityAvailable: 0,
		ReorderPoint:      10,
	}
	if err := s.invRepo.Create(ctx, inv); err != nil {
		return fmt.Errorf("create inventory: %w", err)
	}

	return nil
}

func (s *ProductService) UpdateVariant(ctx context.Context, v *models.Variant) error {
	v.UpdatedAt = time.Now().UTC()
	if err := s.varRepo.Update(ctx, v); err != nil {
		return fmt.Errorf("update variant: %w", err)
	}

	// Update inventory
	inv, err := s.invRepo.GetByVariantID(ctx, v.ID)
	if err == nil {
		inv.QuantityOnHand = v.WeightGrams // Using weight_grams as stock
		s.invRepo.Update(ctx, inv)
	}

	return nil
}

func (s *ProductService) DeleteVariant(ctx context.Context, id uuid.UUID) error {
	return s.varRepo.Delete(ctx, id)
}

func (s *ProductService) AddVariantImage(ctx context.Context, variantID uuid.UUID, url, altText string, sortOrder int32) (*models.VariantImage, error) {
	img := &models.VariantImage{
		ID:        uuid.New(),
		VariantID: variantID,
		URL:       url,
		AltText:   altText,
		SortOrder: int(sortOrder),
		CreatedAt: time.Now().UTC(),
	}

	if err := s.imgRepo.Create(ctx, img); err != nil {
		return nil, fmt.Errorf("create variant image: %w", err)
	}
	return img, nil
}

func (s *ProductService) RemoveVariantImage(ctx context.Context, id uuid.UUID) error {
	return s.imgRepo.Delete(ctx, id)
}

func (s *ProductService) ReorderVariantImages(ctx context.Context, orders map[uuid.UUID]int32) error {
	for id, sortOrder := range orders {
		if err := s.imgRepo.UpdateSortOrder(ctx, id, sortOrder); err != nil {
			return fmt.Errorf("update sort order for %s: %w", id, err)
		}
	}
	return nil
}

func (s *ProductService) GetImagesByVariantID(ctx context.Context, variantID uuid.UUID) ([]models.VariantImage, error) {
	return s.imgRepo.GetByVariantID(ctx, variantID)
}
