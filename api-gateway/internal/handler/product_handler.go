package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	productpb "github.com/manojnegi/ecomm-microservices/gen/go/product/v1"
)

type ProductHandler struct {
	productClient productpb.ProductServiceClient
}

func NewProductHandler(productClient productpb.ProductServiceClient) *ProductHandler {
	return &ProductHandler{productClient: productClient}
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
