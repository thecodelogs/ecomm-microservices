package handler

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
	productpb "github.com/manojnegi/ecomm-microservices/gen/go/product/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/auth"
	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/config"
	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/models"
	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/service"
)

type BrandHandler struct {
	productpb.UnimplementedBrandServiceServer

	brandSvc *service.BrandService
	cfg      *config.Config
}

func NewBrandHandler(brandSvc *service.BrandService, cfg *config.Config) *BrandHandler {
	return &BrandHandler{brandSvc: brandSvc, cfg: cfg}
}

func (h *BrandHandler) CreateBrand(ctx context.Context, req *productpb.CreateBrandRequest) (*productpb.CreateBrandResponse, error) {
	claims, err := auth.ExtractClaims(ctx, h.cfg.PasetoSecret)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: %v", err)
	}
	if !strings.EqualFold(claims.Role, "admin") {
		return nil, status.Error(codes.PermissionDenied, "access denied: admin only")
	}

	brandModel := &models.Brand{
		Name: req.Name,
		Description: sql.NullString{
			String: req.Description,
			Valid:  req.Description != "",
		},
		ImageURL: sql.NullString{
			String: req.ImageUrl,
			Valid:  req.ImageUrl != "",
		},
		IsActive: req.IsActive,
	}

	if err := h.brandSvc.CreateBrand(ctx, brandModel); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create brand: %v", err)
	}

	return &productpb.CreateBrandResponse{
		Id:    brandModel.ID.String(),
		Brand: toProtoBrand(brandModel),
	}, nil
}

func (h *BrandHandler) UpdateBrand(ctx context.Context, req *productpb.UpdateBrandRequest) (*productpb.UpdateBrandResponse, error) {
	claims, err := auth.ExtractClaims(ctx, h.cfg.PasetoSecret)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: %v", err)
	}
	if !strings.EqualFold(claims.Role, "admin") {
		return nil, status.Error(codes.PermissionDenied, "access denied: admin only")
	}

	if _, err := uuid.Parse(req.Id); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid brand id")
	}

	existing, err := h.brandSvc.GetBrand(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "brand not found")
	}

	existing.Name = req.Name
	existing.Description = sql.NullString{
		String: req.Description,
		Valid:  req.Description != "",
	}
	existing.ImageURL = sql.NullString{
		String: req.ImageUrl,
		Valid:  req.ImageUrl != "",
	}
	existing.IsActive = req.IsActive

	if err := h.brandSvc.UpdateBrand(ctx, existing); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update brand: %v", err)
	}

	return &productpb.UpdateBrandResponse{
		Brand: toProtoBrand(existing),
	}, nil
}

func (h *BrandHandler) DeleteBrand(ctx context.Context, req *productpb.DeleteBrandRequest) (*productpb.DeleteBrandResponse, error) {
	claims, err := auth.ExtractClaims(ctx, h.cfg.PasetoSecret)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: %v", err)
	}
	if !strings.EqualFold(claims.Role, "admin") {
		return nil, status.Error(codes.PermissionDenied, "access denied: admin only")
	}

	if _, err := uuid.Parse(req.Id); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid brand id")
	}

	if err := h.brandSvc.DeleteBrand(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete brand: %v", err)
	}

	return &productpb.DeleteBrandResponse{
		Success: true,
	}, nil
}

func (h *BrandHandler) GetBrand(ctx context.Context, req *productpb.GetBrandRequest) (*productpb.GetBrandResponse, error) {
	brand, err := h.brandSvc.GetBrand(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "brand not found")
	}

	return &productpb.GetBrandResponse{
		Brand: toProtoBrand(brand),
	}, nil
}

func (h *BrandHandler) ListBrands(ctx context.Context, req *productpb.ListBrandsRequest) (*productpb.BrandListResponse, error) {
	brands, total, err := h.brandSvc.ListBrands(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoBrands []*productpb.Brand
	for _, b := range brands {
		bCopy := b
		protoBrands = append(protoBrands, toProtoBrand(&bCopy))
	}

	return &productpb.BrandListResponse{
		Brands: protoBrands,
		Total:  total,
		Page:   req.Page,
	}, nil
}

func toProtoBrand(b *models.Brand) *productpb.Brand {
	return &productpb.Brand{
		Id:          b.ID.String(),
		Name:        b.Name,
		Description: b.Description.String,
		ImageUrl:    b.ImageURL.String,
		IsActive:    b.IsActive,
		CreatedAt:   b.CreatedAt.Unix(),
		UpdatedAt:   b.UpdatedAt.Unix(),
	}
}
