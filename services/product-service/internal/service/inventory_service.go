package service

import (
	"context"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/repository"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/models"

	"github.com/google/uuid"
)

type InventoryService struct {
	invRepo *repository.InventoryRepo
}

func NewInventoryService(invRepo *repository.InventoryRepo) *InventoryService {
	return &InventoryService{invRepo: invRepo}
}

func (s *InventoryService) GetInventory(ctx context.Context, variantID uuid.UUID) (*models.Inventory, error) {
	return s.invRepo.GetByVariantID(ctx, variantID)
}

func (s *InventoryService) ReserveStock(ctx context.Context, variantID uuid.UUID, quantity int) error {
	return s.invRepo.ReserveStock(ctx, variantID, quantity)
}

func (s *InventoryService) CommitStock(ctx context.Context, variantID uuid.UUID, quantity int) error {
	return s.invRepo.CommitStock(ctx, variantID, quantity)
}

func (s *InventoryService) ReleaseStock(ctx context.Context, variantID uuid.UUID, quantity int) error {
	return s.invRepo.ReleaseStock(ctx, variantID, quantity)
}

func (s *InventoryService) UpdateInventory(ctx context.Context, variantID uuid.UUID, quantityOnHand int, reorderPoint int) (*models.Inventory, error) {
	inv, err := s.invRepo.GetByVariantID(ctx, variantID)
	if err != nil {
		// If it doesn't exist, create it
		newInv := &models.Inventory{
			ID:               uuid.New(),
			VariantID:        variantID,
			QuantityOnHand:   quantityOnHand,
			QuantityReserved: 0,
			ReorderPoint:     reorderPoint,
		}
		if err := s.invRepo.Create(ctx, newInv); err != nil {
			return nil, err
		}
		return s.invRepo.GetByVariantID(ctx, variantID)
	}

	inv.QuantityOnHand = quantityOnHand
	inv.ReorderPoint = reorderPoint

	if err := s.invRepo.Update(ctx, inv); err != nil {
		return nil, err
	}

	return s.invRepo.GetByVariantID(ctx, variantID)
}
