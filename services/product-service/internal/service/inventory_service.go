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
