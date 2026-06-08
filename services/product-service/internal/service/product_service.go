package service

import (
	"context"
	"fmt"
	"time"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/models"
	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/repository"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
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

	return s.prodRepo.RunInTx(ctx, func(ctx context.Context) error {
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
				QuantityOnHand:    variants[i].InitialStock,
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
	})
}

func (s *ProductService) UpdateProduct(ctx context.Context, p *models.Product, updateVariants bool, variants []models.Variant) error {
	p.UpdatedAt = time.Now().UTC()

	return s.prodRepo.RunInTx(ctx, func(ctx context.Context) error {
		if err := s.prodRepo.Update(ctx, p); err != nil {
			return fmt.Errorf("update product: %w", err)
		}

		if updateVariants {
			existingVariants, err := s.varRepo.GetByProductID(ctx, p.ID)
			if err != nil && err.Error() != "no rows in result set" {
				// ignore "no rows" error, but handle others if needed. The repo might just return an empty slice without an error.
			}

			existingMap := make(map[uuid.UUID]models.Variant)
			for _, ev := range existingVariants {
				existingMap[ev.ID] = ev
			}

			incomingMap := make(map[uuid.UUID]bool)

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
						QuantityOnHand: variants[i].InitialStock,
						ReorderPoint:   10,
					}
					if err := s.invRepo.Create(ctx, inv); err != nil {
						return fmt.Errorf("create inventory: %w", err)
					}
					incomingMap[variants[i].ID] = true
				} else {
					incomingMap[variants[i].ID] = true
					variants[i].ProductID = p.ID
					variants[i].UpdatedAt = time.Now().UTC()
					if err := s.varRepo.Update(ctx, &variants[i]); err != nil {
						return fmt.Errorf("update variant: %w", err)
					}
					// Ensure inventory exists, but do NOT overwrite QuantityOnHand
					_, err := s.invRepo.GetByVariantID(ctx, variants[i].ID)
					if err != nil {
						// If inventory didn't exist for some reason, create it
						inv := &models.Inventory{
							ID:             uuid.New(),
							VariantID:      variants[i].ID,
							QuantityOnHand: variants[i].InitialStock,
							ReorderPoint:   10,
						}
						if err := s.invRepo.Create(ctx, inv); err != nil {
							return fmt.Errorf("create inventory: %w", err)
						}
					}
				}

				if len(variants[i].Images) > 0 {
					existingImages, _ := s.imgRepo.GetByVariantID(ctx, variants[i].ID)
					maxSortOrder := -1
					for _, img := range existingImages {
						if img.SortOrder > maxSortOrder {
							maxSortOrder = img.SortOrder
						}
					}

					for j := range variants[i].Images {
						if variants[i].Images[j].ID == uuid.Nil {
							variants[i].Images[j].ID = uuid.New()
						}
						variants[i].Images[j].VariantID = variants[i].ID
						variants[i].Images[j].SortOrder = maxSortOrder + 1 + j
						variants[i].Images[j].CreatedAt = time.Now().UTC()
						if err := s.imgRepo.Create(ctx, &variants[i].Images[j]); err != nil {
							return fmt.Errorf("create variant image: %w", err)
						}
					}
				}
			}

			// Delete variants not in the incoming list
			for id := range existingMap {
				if !incomingMap[id] {
					if err := s.varRepo.Delete(ctx, id); err != nil {
						return fmt.Errorf("delete missing variant: %w", err)
					}
				}
			}
		}

		return nil
	})
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
	product, err := s.prodRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, nil, err
	}

	variants, err := s.varRepo.GetByProductID(ctx, product.ID)
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
	baseSlug := slug.Make(name)
	return fmt.Sprintf("%s-%s", baseSlug, uuid.New().String()[:8])
}

func (s *ProductService) CreateVariant(ctx context.Context, v *models.Variant) error {
	v.ID = uuid.New()
	v.CreatedAt = time.Now().UTC()
	v.UpdatedAt = time.Now().UTC()

	return s.prodRepo.RunInTx(ctx, func(ctx context.Context) error {
		if err := s.varRepo.Create(ctx, v); err != nil {
			return fmt.Errorf("create variant: %w", err)
		}

		// Create default inventory
		inv := &models.Inventory{
			ID:                uuid.New(),
			VariantID:         v.ID,
			QuantityOnHand:    v.InitialStock,
			QuantityReserved:  0,
			QuantityAvailable: 0,
			ReorderPoint:      10,
		}
		if err := s.invRepo.Create(ctx, inv); err != nil {
			return fmt.Errorf("create inventory: %w", err)
		}

		return nil
	})
}

func (s *ProductService) UpdateVariant(ctx context.Context, v *models.Variant) error {
	v.UpdatedAt = time.Now().UTC()

	return s.prodRepo.RunInTx(ctx, func(ctx context.Context) error {
		if err := s.varRepo.Update(ctx, v); err != nil {
			return fmt.Errorf("update variant: %w", err)
		}
		return nil
	})
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
