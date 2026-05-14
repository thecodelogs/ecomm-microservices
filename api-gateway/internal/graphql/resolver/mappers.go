package resolver

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	productpb "github.com/manojnegi/ecomm-microservices/gen/go/product/v1"
	userpb "github.com/manojnegi/ecomm-microservices/gen/go/user/v1"
	"github.com/manojnegi/ecommerce/api-gateway/internal/graphql/model"
)

func mapUserFromProto(u *userpb.User) *model.User {
	if u == nil {
		return nil
	}

	return &model.User{
		ID:        u.Id,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Username:  u.Email, // Map Email to Username
		Phone:     &u.Phone,
		Role:      mapRoleFromProto(u.Role),
		Status:    mapStatusFromProto(u.Status),
		CreatedAt: time.Unix(u.CreatedAt, 0),
		UpdatedAt: time.Unix(u.CreatedAt, 0), // Use CreatedAt as UpdatedAt if missing
	}
}

func mapAddressFromProto(a *userpb.Address) *model.Address {
	if a == nil {
		return nil
	}

	return &model.Address{
		ID:        a.Id,
		UserID:    a.UserId,
		Label:     a.Label,
		Street:    a.Line1,
		City:      a.City,
		State:     a.State,
		ZipCode:   a.PostalCode,
		Country:   a.CountryCode,
		IsDefault: a.IsDefault,
	}
}

func mapProductFromProto(p *productpb.Product, baseURL string) *model.Product {
	if p == nil {
		return nil
	}

	imageUrl := p.ImageUrl
	if imageUrl != "" && !strings.HasPrefix(imageUrl, "http") {
		imageUrl = baseURL + imageUrl
	}

	return &model.Product{
		ID:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Sku:         p.Slug, // Use Slug as Sku since Sku is missing
		Stock:       int(p.Stock),
		CategoryID:  p.CategoryId,
		Images:      []string{imageUrl}, // Wrap ImageUrl in a slice
		CreatedAt:   time.Unix(p.CreatedAt, 0),
		UpdatedAt:   time.Unix(p.CreatedAt, 0),
	}
}

func mapCategoryFromProto(c *productpb.Category, baseURL string) *model.Category {
	if c == nil {
		return nil
	}

	imageUrl := c.ImageUrl
	if imageUrl != "" && !strings.HasPrefix(imageUrl, "http") {
		imageUrl = baseURL + imageUrl
	}

	return &model.Category{
		ID:          c.Id,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: &c.Description,
		ImageURL:    &imageUrl,
		SortOrder:   int(c.SortOrder),
		IsActive:    c.IsActive,
		ParentID:    &c.ParentId,
	}
}

func mapRoleFromProto(r string) model.Role {
	switch r {
	case "admin":
		return model.RoleAdmin
	default:
		return model.RoleCustomer
	}
}

func mapStatusFromProto(s string) model.UserStatus {
	switch s {
	case "INACTIVE":
		return model.UserStatusInactive
	case "SUSPENDED":
		return model.UserStatusSuspended
	default:
		return model.UserStatusActive
	}
}

func mapStatusToProto(s model.UserStatus) string {
	return string(s)
}

func encodeCursor(id string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor:%s", id)))
}
