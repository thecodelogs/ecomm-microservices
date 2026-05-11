package handler

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	categorypb "github.com/manojnegi/ecomm-microservices/gen/go/product/v1"
	productpb "github.com/manojnegi/ecomm-microservices/gen/go/product/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/models"
	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/service"
)

type CategoryHandler struct {
	categorypb.UnimplementedCategoryServiceServer

	// prodSvc *service.ProductService
	catSvc *service.CategoryService
	// invSvc  *service.InventoryService
	// revSvc  *service.ReviewService
}

func NewCategoryHandler(catSvc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{catSvc: catSvc}
}

func (h *CategoryHandler) CreateCategory(ctx context.Context, req *categorypb.CreateCategoryRequest) (*productpb.CreateCategoryResponse, error) {
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
