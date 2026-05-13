package handler

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
	categorypb "github.com/manojnegi/ecomm-microservices/gen/go/product/v1"
	productpb "github.com/manojnegi/ecomm-microservices/gen/go/product/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/auth"
	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/config"
	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/models"
	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/service"
)

type CategoryHandler struct {
	categorypb.UnimplementedCategoryServiceServer

	catSvc *service.CategoryService
	cfg    *config.Config
}

func NewCategoryHandler(catSvc *service.CategoryService, cfg *config.Config) *CategoryHandler {
	return &CategoryHandler{catSvc: catSvc, cfg: cfg}
}

func (h *CategoryHandler) CreateCategory(ctx context.Context, req *categorypb.CreateCategoryRequest) (*productpb.CreateCategoryResponse, error) {
	// ── Auth Check ──
	claims, err := auth.ExtractClaims(ctx, h.cfg.PasetoSecret)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: %v", err)
	}

	if !strings.EqualFold(claims.Role, "admin") {
		return nil, status.Error(codes.PermissionDenied, "access denied: admin only")
	}

	var parentID uuid.NullUUID

	if req.ParentId != "" {
		parsedID, err := uuid.Parse(req.ParentId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid parent id")
		}

		parentID = uuid.NullUUID{
			UUID:  parsedID,
			Valid: true,
		}
	}

	categoryModel := &models.Category{
		Name:     req.Name,
		ParentID: parentID,
		ImageURL: sql.NullString{
			String: req.ImageUrl,
			Valid:  req.ImageUrl != "",
		},
		Slug: req.Slug,
		Description: sql.NullString{
			String: req.Description,
			Valid:  req.Description != "",
		},
		SortOrder: int(req.SortOrder),
		IsActive:  req.IsActive,
	}

	if err := h.catSvc.CreateCategory(ctx, categoryModel); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create category: %v", err)
	}

	return &productpb.CreateCategoryResponse{
		Id: categoryModel.ID.String(),
	}, nil
}

func (h *CategoryHandler) ListCategories(ctx context.Context, req *categorypb.ListCategoryRequest) (*categorypb.CategoryListResponse, error) {
	// ── Auth Check ──
	claims, err := auth.ExtractClaims(ctx, h.cfg.PasetoSecret)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: %v", err)
	}

	if !strings.EqualFold(claims.Role, "admin") {
		return nil, status.Error(codes.PermissionDenied, "access denied: admin only")
	}

	categories, total, err := h.catSvc.ListCategories(ctx, req.Page, req.PageSize)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var catProducts []*productpb.Category
	for _, p := range categories {
		catProducts = append(catProducts, toProtoCategory(&p))
	}

	return &categorypb.CategoryListResponse{
		Category: catProducts,
		Total:    total,
		Page:     req.Page,
	}, nil
}

func toProtoCategory(p *models.Category) *categorypb.Category {
	return &productpb.Category{
		Id:          p.ID.String(),
		Slug:        p.Slug,
		Name:        p.Name,
		Description: p.Description.String,
		IsActive:    p.IsActive,
		ImageUrl:    p.ImageURL.String,
		SortOrder:   int32(p.SortOrder),
		ParentId:    p.ParentID.UUID.String(),
	}
}
