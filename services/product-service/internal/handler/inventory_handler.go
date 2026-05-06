package handler

import (
	"context"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/service"

	productpb "github.com/manojnegi/ecomm-microservices/gen/go/product/v1"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryHandler struct {
	productpb.UnimplementedInventoryServiceServer
	invSvc *service.InventoryService
}

func NewInventoryHandler(invSvc *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{invSvc: invSvc}
}

// func (h *InventoryHandler) ReserveStock(ctx context.Context, req *productpb.ReserveStockRequest) (*productpb.ReserveStockResponse, error) {
// 	variantID, err := uuid.Parse(req.VariantId)
// 	if err != nil {
// 		return nil, status.Error(codes.InvalidArgument, "invalid variant id")
// 	}

// 	if err := h.invSvc.ReserveStock(ctx, variantID, int(req.Quantity)); err != nil {
// 		return nil, status.Error(codes.FailedPrecondition, err.Error())
// 	}

// 	return &productpb.ReserveStockResponse{
// 		ReservationId: uuid.New().String(),
// 		Success:       true,
// 	}, nil
// }

// func (h *InventoryHandler) CommitStock(ctx context.Context, req *productpb.CommitStockRequest) (*productpb.Empty, error) {
// 	variantID, err := uuid.Parse(req.VariantId)
// 	if err != nil {
// 		return nil, status.Error(codes.InvalidArgument, "invalid variant id")
// 	}

// 	if err := h.invSvc.CommitStock(ctx, variantID, int(req.Quantity)); err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &productpb.Empty{}, nil
// }

// func (h *InventoryHandler) ReleaseStock(ctx context.Context, req *productpb.ReleaseStockRequest) (*productpb.Empty, error) {
// 	variantID, err := uuid.Parse(req.VariantId)
// 	if err != nil {
// 		return nil, status.Error(codes.InvalidArgument, "invalid variant id")
// 	}

// 	if err := h.invSvc.ReleaseStock(ctx, variantID, int(req.Quantity)); err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &productpb.Empty{}, nil
// }

func (h *InventoryHandler) GetInventory(ctx context.Context, req *productpb.GetInventoryRequest) (*productpb.Inventory, error) {
	variantID, err := uuid.Parse(req.VariantId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid variant id")
	}

	inv, err := h.invSvc.GetInventory(ctx, variantID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &productpb.Inventory{
		VariantId:         inv.VariantID.String(),
		QuantityOnHand:    int32(inv.QuantityOnHand),
		QuantityReserved:  int32(inv.QuantityReserved),
		QuantityAvailable: int32(inv.QuantityAvailable),
	}, nil
}
