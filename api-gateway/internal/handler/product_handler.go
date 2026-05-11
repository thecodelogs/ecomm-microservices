package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	productpb "github.com/manojnegi/ecomm-microservices/gen/go/product/v1"
	"github.com/manojnegi/ecommerce/api-gateway/internal/storage"
)

type ProductHandler struct {
	productClient  productpb.ProductServiceClient
	categoryclient productpb.CategoryServiceClient
	s3             *storage.S3Storage
}

func NewProductHandler(
	productClient productpb.ProductServiceClient,
	categoryclient productpb.CategoryServiceClient,
	s3 *storage.S3Storage,
) *ProductHandler {
	return &ProductHandler{
		productClient:  productClient,
		categoryclient: categoryclient,
		s3:             s3,
	}
}

func (p *ProductHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	categoryIDStr := c.Query("category_id")

	var categoryID string

	if categoryIDStr != "" {
		_, err := uuid.Parse(categoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid category_id",
			})
			return
		}

		categoryID = categoryIDStr
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := p.productClient.ListProducts(ctx, &productpb.ListProductsRequest{
		CategoryId: categoryID,
		Page:       int32(page),
		PageSize:   int32(pageSize),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (p *ProductHandler) CreateCategory(c *gin.Context) {

	name := c.PostForm("name")
	description := c.PostForm("description")
	slug := c.PostForm("slug")

	sortOrderStr := c.DefaultPostForm("sort_order", "0")
	isActiveStr := c.DefaultPostForm("is_active", "true")

	parentID := c.PostForm("parent_id")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "name is required",
		})
		return
	}

	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "slug is required",
		})
		return
	}

	if parentID != "" {
		_, err := uuid.Parse(parentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid parent_id",
			})
			return
		}
	}

	sortOrder, err := strconv.Atoi(sortOrderStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid sort_order",
		})
		return
	}

	isActive := isActiveStr == "true"

	// Upload image to S3
	var imageURL string

	file, err := c.FormFile("image")

	if err == nil {

		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to open image",
			})
			return
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to read image",
			})
			return
		}

		key := fmt.Sprintf(
			"categories/%s-%s",
			uuid.New().String(),
			file.Filename,
		)

		imageURL, err = p.s3.UploadFile(
			c.Request.Context(),
			key,
			data,
			file.Header.Get("Content-Type"),
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	ctx, cancel := context.WithTimeout(
		c.Request.Context(),
		5*time.Second,
	)
	defer cancel()

	resp, err := p.categoryclient.CreateCategory(
		ctx,
		&productpb.CreateCategoryRequest{
			Name:        name,
			Description: description,
			Slug:        slug,
			ImageUrl:    imageURL,
			SortOrder:   int32(sortOrder),
			IsActive:    isActive,
			ParentId:    parentID,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        resp.Id,
		"image_url": imageURL,
	})
}
