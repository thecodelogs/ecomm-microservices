package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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

	var variantsToCreate []models.Variant
	if len(req.Variants) > 0 {
		for _, v := range req.Variants {
			var vID uuid.UUID
			if v.Id != "" {
				parsed, err := uuid.Parse(v.Id)
				if err == nil {
					vID = parsed
				}
			}
			
			options := json.RawMessage(nil)
			if v.Options != "" {
				options = json.RawMessage(v.Options)
			}

			variantsToCreate = append(variantsToCreate, models.Variant{
				ID:             vID,
				SKU:            v.Sku,
				Name:           v.Name,
				Options:        options,
				Price:          int64(v.Price * 100),
				CompareAtPrice: sql.NullInt64{Int64: int64(v.CompareAtPrice * 100), Valid: v.CompareAtPrice > 0},
				CostPrice:      sql.NullInt64{Int64: int64(v.CostPrice * 100), Valid: v.CostPrice > 0},
				WeightGrams:    int(v.WeightGrams),
				IsActive:       v.IsActive,
			})
		}
	} else {
		variantsToCreate = []models.Variant{
			{
				ID:       uuid.New(),
				SKU:      req.Slug + "-default",
				Name:     req.Name,
				IsActive: true,
			},
		}
	}

	if err := h.prodSvc.CreateProduct(ctx, p, variantsToCreate); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	return &productpb.CreateProductResponse{
		Id:      p.ID.String(),
		Product: toProtoProduct(p, variantsToCreate),
	}, nil
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *productpb.GetProductRequest) (*productpb.GetProductResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product id")
	}

	product, variants, err := h.prodSvc.GetProduct(ctx, id)
	fmt.Printf("GetProduct id=%s, variants=%+v, err=%v\n", id, variants, err)
	if err != nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return &productpb.GetProductResponse{
		Product: toProtoProduct(product, variants),
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

	var variantsToUpdate []models.Variant
	if len(req.Variants) > 0 {
		for _, v := range req.Variants {
			var vID uuid.UUID
			if v.Id != "" {
				parsed, err := uuid.Parse(v.Id)
				if err == nil {
					vID = parsed
				}
			}
			
			options := json.RawMessage(nil)
			if v.Options != "" {
				options = json.RawMessage(v.Options)
			}

			variantsToUpdate = append(variantsToUpdate, models.Variant{
				ID:             vID,
				ProductID:      id,
				SKU:            v.Sku,
				Name:           v.Name,
				Options:        options,
				Price:          int64(v.Price * 100), // Note: Since we changed DB to BIGINT, we keep storing as cents? Wait!
				CompareAtPrice: sql.NullInt64{Int64: int64(v.CompareAtPrice * 100), Valid: v.CompareAtPrice > 0},
				CostPrice:      sql.NullInt64{Int64: int64(v.CostPrice * 100), Valid: v.CostPrice > 0},
				WeightGrams:    int(v.WeightGrams),
				IsActive:       v.IsActive,
			})
		}
	} else {
		// Fallback to updating default variant if no variants provided
		variantsToUpdate = []models.Variant{
			{
				ProductID: id,
				SKU:       req.Slug + "-default",
				Name:      req.Name,
				IsActive:  true,
			},
		}
	}

	if err := h.prodSvc.UpdateProduct(ctx, p, variantsToUpdate); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}

	return &productpb.UpdateProductResponse{
		Product: toProtoProduct(p, variantsToUpdate),
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
		variants, _ := h.prodSvc.GetVariantsByProductID(ctx, p.ID)
		pbProducts = append(pbProducts, toProtoProduct(&p, variants))
	}

	return &productpb.ProductListResponse{
		Products: pbProducts,
		Total:    total,
		Page:     req.Page,
	}, nil
}

func (h *ProductHandler) CreateVariant(ctx context.Context, req *productpb.CreateVariantRequest) (*productpb.CreateVariantResponse, error) {
	_, err := h.checkAdminAuth(ctx)
	if err != nil {
		return nil, err
	}

	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product id")
	}

	options := []byte(req.Options)
	if req.Options == "" {
		options = nil
	}

	v := &models.Variant{
		ProductID:      productID,
		SKU:            req.Sku,
		Name:           req.Name,
		Options:        options,
		Price:          int64(req.Price * 100),
		CompareAtPrice: sql.NullInt64{Int64: int64(req.CompareAtPrice * 100), Valid: req.CompareAtPrice > 0},
		CostPrice:      sql.NullInt64{Int64: int64(req.CostPrice * 100), Valid: req.CostPrice > 0},
		WeightGrams:    int(req.WeightGrams),
		ImageURL:       toNullString(req.ImageUrl),
		IsActive:       req.IsActive,
	}

	if err := h.prodSvc.CreateVariant(ctx, v); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create variant: %v", err)
	}

	return &productpb.CreateVariantResponse{
		Variant: toProtoVariants([]models.Variant{*v})[0],
	}, nil
}

func (h *ProductHandler) UpdateVariant(ctx context.Context, req *productpb.UpdateVariantRequest) (*productpb.UpdateVariantResponse, error) {
	_, err := h.checkAdminAuth(ctx)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid variant id")
	}
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product id")
	}

	options := []byte(req.Options)
	if req.Options == "" {
		options = nil
	}

	v := &models.Variant{
		ID:             id,
		ProductID:      productID,
		SKU:            req.Sku,
		Name:           req.Name,
		Options:        options,
		Price:          int64(req.Price * 100),
		CompareAtPrice: sql.NullInt64{Int64: int64(req.CompareAtPrice * 100), Valid: req.CompareAtPrice > 0},
		CostPrice:      sql.NullInt64{Int64: int64(req.CostPrice * 100), Valid: req.CostPrice > 0},
		WeightGrams:    int(req.WeightGrams),
		ImageURL:       toNullString(req.ImageUrl),
		IsActive:       req.IsActive,
	}

	if err := h.prodSvc.UpdateVariant(ctx, v); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update variant: %v", err)
	}

	return &productpb.UpdateVariantResponse{
		Variant: toProtoVariants([]models.Variant{*v})[0],
	}, nil
}

func (h *ProductHandler) DeleteVariant(ctx context.Context, req *productpb.DeleteVariantRequest) (*productpb.DeleteVariantResponse, error) {
	_, err := h.checkAdminAuth(ctx)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid variant id")
	}

	if err := h.prodSvc.DeleteVariant(ctx, id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete variant: %v", err)
	}

	return &productpb.DeleteVariantResponse{Success: true}, nil
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

func toProtoProduct(p *models.Product, variants []models.Variant) *productpb.Product {
	pb := &productpb.Product{
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

	if len(variants) > 0 {
		pb.Variants = toProtoVariants(variants)
	}

	return pb
}

func toProtoVariants(vs []models.Variant) []*productpb.Variant {
	var result []*productpb.Variant
	for _, v := range vs {
		optionsStr := ""
		if v.Options != nil {
			optionsStr = string(v.Options)
		}

		imageURL := ""
		if v.ImageURL.Valid {
			imageURL = v.ImageURL.String
		}

		result = append(result, &productpb.Variant{
			Id:             v.ID.String(),
			ProductId:      v.ProductID.String(),
			Sku:            v.SKU,
			Name:           v.Name,
			Options:        optionsStr,
			Price:          float64(v.Price) / 100.0,
			CompareAtPrice: float64(v.CompareAtPrice.Int64) / 100.0,
			CostPrice:      float64(v.CostPrice.Int64) / 100.0,
			WeightGrams:    int32(v.WeightGrams),
			ImageUrl:       imageURL,
			IsActive:       v.IsActive,
			CreatedAt:      v.CreatedAt.Unix(),
			UpdatedAt:      v.UpdatedAt.Unix(),
		})
	}
	return result
}

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
