package handler

import (
	"context"
	"database/sql"
	"strings"

	productpb "github.com/manojnegi/ecomm-microservices/gen/go/product/v1"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/auth"
	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/config"
	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/models"
	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/service"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductHandler struct {
	productpb.UnimplementedProductServiceServer
	// productpb.UnimplementedInventoryServiceServer

	prodSvc *service.ProductService
	invSvc  *service.InventoryService
	revSvc  *service.ReviewService
	cfg     *config.Config
}

func NewProductHandler(prodSvc *service.ProductService, invSvc *service.InventoryService, revSvc *service.ReviewService, cfg *config.Config) *ProductHandler {
	return &ProductHandler{prodSvc: prodSvc, invSvc: invSvc, revSvc: revSvc, cfg: cfg}
}

// ── ProductService RPCs ──
// func (h *ProductHandler) CreateProduct(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.CreateProductResponse, error) {
// 	categoryID, err := uuid.Parse(req.CategoryId)
// 	if err != nil {
// 		return nil, status.Error(codes.InvalidArgument, "invalid category id")
// 	}

// 	vendorID, err := uuid.Parse(req.VendorId)
// 	if err != nil {
// 		return nil, status.Error(codes.InvalidArgument, "invalid vendor id")
// 	}

// 	productModel := &models.Product{
// 		Name:        req.Name,
// 		Description: req.Description,
// 		CategoryID:  categoryID,
// 		ImageUrl:    req.ImageUrl,
// 		Slug:        req.Slug,
// 		ShortDescription: sql.NullString{
// 			String: req.ShortDescription,
// 			Valid:  req.ShortDescription != "",
// 		},
// 		Brand: sql.NullString{
// 			String: req.Brand,
// 			Valid:  req.Brand != "",
// 		},
// 		Tags:       req.Tags,
// 		Attributes: json.RawMessage(req.Attributes),
// 		Status:     req.Status,
// 		VendorID:   vendorID,
// 	}

// 	// 2. Call your service layer to persist the data
// 	// Assuming your service updates productModel with a new UUID/ID
// 	if err := h.prodSvc.CreateProduct(ctx, productModel); err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
// 	}

// 	// 3. Return the EXACT response type requested by the .proto file
// 	return &productpb.CreateProductResponse{
// 		Id:          productModel.ID.String(),
// 		Name:        productModel.Name,
// 		Description: productModel.Description,
// 		Price:       req.Price, // Usually you'd return the saved value
// 		Stock:       req.Stock, // Usually you'd return the saved value
// 		Category:    req.Category,
// 		ImageUrl:    req.ImageUrl,
// 	}, nil
// }

// func (h *ProductHandler) GetProduct(ctx context.Context, req *productpb.GetProductRequest) (*productpb.ProductDetail, error) {
// 	id, err := uuid.Parse(req.ProductId)
// 	if err != nil {
// 		return nil, status.Error(codes.InvalidArgument, "invalid product id")
// 	}

// 	product, variants, err := h.prodSvc.GetProduct(ctx, id)
// 	if err != nil {
// 		return nil, status.Error(codes.NotFound, err.Error())
// 	}

// 	return &productpb.ProductDetail{
// 		Product:  toProtoProduct(product),
// 		Variants: toProtoVariants(variants),
// 	}, nil
// }

func (h *ProductHandler) ListProducts(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ProductListResponse, error) {
	// ── Auth Check ──
	claims, err := auth.ExtractClaims(ctx, h.cfg.PasetoSecret)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: %v", err)
	}

	if !strings.EqualFold(claims.Role, "admin") {
		return nil, status.Error(codes.PermissionDenied, "access denied: admin only")
	}

	catID, _ := uuid.Parse(req.CategoryId)
	products, total, err := h.prodSvc.ListProducts(ctx, catID, req.Page, req.PageSize)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var pbProducts []*productpb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, toProtoProduct(&p))
	}

	return &productpb.ProductListResponse{
		Products: pbProducts,
		Total:    total,
		Page:     req.Page,
	}, nil
}

// func (h *ProductHandler) GetVariantsBatch(ctx context.Context, req *productpb.GetVariantsBatchRequest) (*productpb.VariantList, error) {
// 	var ids []uuid.UUID
// 	for _, id := range req.VariantIds {
// 		ids = append(ids, uuid.MustParse(id))
// 	}

// 	variants, err := h.prodSvc.GetVariantsBatch(ctx, ids)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &productpb.VariantList{Variants: toProtoVariants(variants)}, nil
// }

// ── InventoryService RPCs ──

// func (h *ProductHandler) ReserveStock(ctx context.Context, req *productpb.ReserveStockRequest) (*productpb.ReserveStockResponse, error) {
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

// func (h *ProductHandler) CommitStock(ctx context.Context, req *productpb.CommitStockRequest) (*productpb.Empty, error) {
// 	variantID, err := uuid.Parse(req.VariantId)
// 	if err != nil {
// 		return nil, status.Error(codes.InvalidArgument, "invalid variant id")
// 	}

// 	if err := h.invSvc.CommitStock(ctx, variantID, int(req.Quantity)); err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &productpb.Empty{}, nil
// }

// func (h *ProductHandler) ReleaseStock(ctx context.Context, req *productpb.ReleaseStockRequest) (*productpb.Empty, error) {
// 	variantID, err := uuid.Parse(req.VariantId)
// 	if err != nil {
// 		return nil, status.Error(codes.InvalidArgument, "invalid variant id")
// 	}

// 	if err := h.invSvc.ReleaseStock(ctx, variantID, int(req.Quantity)); err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &productpb.Empty{}, nil
// }

// func (h *ProductHandler) GetInventory(ctx context.Context, req *productpb.GetInventoryRequest) (*productpb.Inventory, error) {
// 	variantID, err := uuid.Parse(req.VariantId)
// 	if err != nil {
// 		return nil, status.Error(codes.InvalidArgument, "invalid variant id")
// 	}

// 	inv, err := h.invSvc.GetInventory(ctx, variantID)
// 	if err != nil {
// 		return nil, status.Error(codes.NotFound, err.Error())
// 	}

// 	return toProtoInventory(inv), nil
// }

// ── Helpers ──

func toProtoProduct(p *models.Product) *productpb.Product {
	return &productpb.Product{
		Id:          p.ID.String(),
		CategoryId:  p.CategoryID.String(),
		Slug:        p.Slug,
		Name:        p.Name,
		Description: p.Description,
		Brand:       p.Brand.String,
		Status:      p.Status,
		VendorId:    p.VendorID.String(),
		AvgRating:   float64(p.AvgRating),
		ReviewCount: int32(p.ReviewCount),
	}
}

// func toProtoVariants(vs []models.Variant) []*productpb.Variant {
// 	var result []*productpb.Variant
// 	for _, v := range vs {
// 		result = append(result, &productpb.Variant{
// 			Id:        v.ID.String(),
// 			ProductId: v.ProductID.String(),
// 			Sku:       v.SKU,
// 			Name:      v.Name,
// 			Price:     v.Price,
// 			IsActive:  v.IsActive,
// 		})
// 	}
// 	return result
// }

func toProtoInventory(i *models.Inventory) *productpb.Inventory {
	return &productpb.Inventory{
		VariantId:         i.VariantID.String(),
		QuantityOnHand:    int32(i.QuantityOnHand),
		QuantityReserved:  int32(i.QuantityReserved),
		QuantityAvailable: int32(i.QuantityAvailable),
	}
}

func toNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
