package handler

import (
	"context"
	"database/sql"
	"encoding/json"
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
func (h *ProductHandler) CreateProduct(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.CreateProductResponse, error) {
	// ── Auth Check ──
	_, err := h.checkAdminAuth(ctx)
	if err != nil {
		return nil, err
	}

	categoryID, _ := uuid.Parse(req.CategoryId)
	vendorID, _ := uuid.Parse(req.VendorId)

	var attributes json.RawMessage
	if req.Attributes != "" {
		attributes = json.RawMessage(req.Attributes)
	}

	p := &models.Product{
		CategoryID:  categoryID,
		Name:        req.Name,
		Description: req.Description,
		Slug:        req.Slug,
		ShortDescription: sql.NullString{
			String: req.ShortDescription,
			Valid:  req.ShortDescription != "",
		},
		Brand: sql.NullString{
			String: req.Brand,
			Valid:  req.Brand != "",
		},
		Tags:       req.Tags,
		Attributes: attributes,
		Status:     req.Status,
		VendorID:   vendorID,
	}

	// Create default variant
	v := models.Variant{
		SKU:      req.Slug + "-default",
		Name:     req.Name,
		Price:    int64(req.Price * 100), // Convert to cents
		IsActive: true,
		// Using weight_grams as placeholder for stock in req if needed
		WeightGrams: int(req.Stock),
	}

	if err := h.prodSvc.CreateProduct(ctx, p, []models.Variant{v}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	return &productpb.CreateProductResponse{
		Id:      p.ID.String(),
		Product: toProtoProduct(p),
	}, nil
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *productpb.GetProductRequest) (*productpb.GetProductResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product id")
	}

	product, variants, err := h.prodSvc.GetProduct(ctx, id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	pbProd := toProtoProduct(product)
	if len(variants) > 0 {
		pbProd.Price = float64(variants[0].Price) / 100.0
		// Get stock from inventory for the first variant
		inv, err := h.invSvc.GetInventory(ctx, variants[0].ID)
		if err == nil {
			pbProd.Stock = int32(inv.QuantityAvailable)
		}
	}

	return &productpb.GetProductResponse{
		Product: pbProd,
	}, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *productpb.UpdateProductRequest) (*productpb.UpdateProductResponse, error) {
	// ── Auth Check ──
	_, err := h.checkAdminAuth(ctx)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product id")
	}

	categoryID, _ := uuid.Parse(req.CategoryId)

	var attributes json.RawMessage
	if req.Attributes != "" {
		attributes = json.RawMessage(req.Attributes)
	}

	p := &models.Product{
		ID:          id,
		CategoryID:  categoryID,
		Name:        req.Name,
		Description: req.Description,
		Slug:        req.Slug,
		ShortDescription: sql.NullString{
			String: req.ShortDescription,
			Valid:  req.ShortDescription != "",
		},
		Brand: sql.NullString{
			String: req.Brand,
			Valid:  req.Brand != "",
		},
		Tags:       req.Tags,
		Attributes: attributes,
		Status:     req.Status,
	}

	// Update default variant
	v := models.Variant{
		ProductID: id,
		SKU:       req.Slug + "-default",
		Name:      req.Name,
		Price:     int64(req.Price * 100),
		IsActive:  true,
		// Using weight_grams as stock placeholder
		WeightGrams: int(req.Stock),
	}

	if err := h.prodSvc.UpdateProduct(ctx, p, []models.Variant{v}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}

	return &productpb.UpdateProductResponse{
		Product: toProtoProduct(p),
	}, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *productpb.DeleteProductRequest) (*productpb.DeleteProductResponse, error) {
	// ── Auth Check ──
	_, err := h.checkAdminAuth(ctx)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product id")
	}

	if err := h.prodSvc.DeleteProduct(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}

	return &productpb.DeleteProductResponse{Success: true}, nil
}

func (h *ProductHandler) ListProducts(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ProductListResponse, error) {

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

func (h *ProductHandler) checkAdminAuth(ctx context.Context) (*auth.AccessTokenClaims, error) {
	claims, err := auth.ExtractClaims(ctx, h.cfg.PasetoSecret)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: %v", err)
	}

	if !strings.EqualFold(claims.Role, "admin") {
		return nil, status.Error(codes.PermissionDenied, "access denied: admin only")
	}
	return claims, nil
}


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
	imageURL := ""
	if len(p.ImageUrl) > 0 {
		imageURL = p.ImageUrl[0]
	}
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
		ImageUrl:    imageURL,
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
