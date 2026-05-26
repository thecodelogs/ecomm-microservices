package handler

import (
	"context"
	"strings"

	productpb "github.com/manojnegi/ecomm-microservices/gen/go/product/v1"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/auth"
	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/config"
	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryHandler struct {
	productpb.UnimplementedInventoryServiceServer

	invSvc *service.InventoryService
	cfg    *config.Config
}

func NewInventoryHandler(invSvc *service.InventoryService, cfg *config.Config) *InventoryHandler {
	return &InventoryHandler{invSvc: invSvc, cfg: cfg}
}

func (h *InventoryHandler) GetInventory(ctx context.Context, req *productpb.GetInventoryRequest) (*productpb.GetInventoryResponse, error) {
	variantID, err := uuid.Parse(req.VariantId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid variant id")
	}

	inv, err := h.invSvc.GetInventory(ctx, variantID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "inventory not found: %v", err)
	}

	return &productpb.GetInventoryResponse{
		Inventory: &productpb.Inventory{
			VariantId:         inv.VariantID.String(),
			QuantityOnHand:    int32(inv.QuantityOnHand),
			QuantityReserved:  int32(inv.QuantityReserved),
			QuantityAvailable: int32(inv.QuantityAvailable),
			ReorderPoint:      int32(inv.ReorderPoint),
		},
	}, nil
}

func (h *InventoryHandler) UpdateInventory(ctx context.Context, req *productpb.UpdateInventoryRequest) (*productpb.UpdateInventoryResponse, error) {
	_, err := h.checkAdminAuth(ctx)
	if err != nil {
		return nil, err
	}

	variantID, err := uuid.Parse(req.VariantId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid variant id")
	}

	inv, err := h.invSvc.UpdateInventory(ctx, variantID, int(req.QuantityOnHand), int(req.ReorderPoint))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update inventory: %v", err)
	}

	return &productpb.UpdateInventoryResponse{
		Inventory: &productpb.Inventory{
			VariantId:         inv.VariantID.String(),
			QuantityOnHand:    int32(inv.QuantityOnHand),
			QuantityReserved:  int32(inv.QuantityReserved),
			QuantityAvailable: int32(inv.QuantityAvailable),
			ReorderPoint:      int32(inv.ReorderPoint),
		},
	}, nil
}

func (h *InventoryHandler) checkAdminAuth(ctx context.Context) (*auth.AccessTokenClaims, error) {
	claims, err := auth.ExtractClaims(ctx, h.cfg.PasetoSecret)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: %v", err)
	}

	if !strings.EqualFold(claims.Role, "admin") {
		return nil, status.Error(codes.PermissionDenied, "access denied: admin only")
	}
	return claims, nil
}
