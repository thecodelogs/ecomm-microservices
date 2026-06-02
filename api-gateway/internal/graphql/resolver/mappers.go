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

	var variants []*model.Variant
	for _, v := range p.Variants {
		variants = append(variants, mapVariantFromProto(v, baseURL))
	}

	var brandPtr *string
	if p.Brand != "" {
		b := p.Brand
		brandPtr = &b
	}
	var brandIdPtr *string
	if p.BrandId != "" {
		bId := p.BrandId
		brandIdPtr = &bId
	}

	return &model.Product{
		ID:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Brand:       brandPtr,
		BrandID:     brandIdPtr,
		Sku:         p.Slug,
		CategoryID:  p.CategoryId,
		Variants:    variants,
		CreatedAt:   time.Unix(p.CreatedAt, 0),
		UpdatedAt:   time.Unix(p.CreatedAt, 0),
	}
}

func mapVariantFromProto(v *productpb.Variant, baseURL string) *model.Variant {
	if v == nil {
		return nil
	}

	var optionsPtr *string
	if v.Options != "" {
		opts := v.Options
		optionsPtr = &opts
	}
	var cmpPricePtr *float64
	if v.CompareAtPrice > 0 {
		cmp := v.CompareAtPrice
		cmpPricePtr = &cmp
	}
	var costPricePtr *float64
	if v.CostPrice > 0 {
		cst := v.CostPrice
		costPricePtr = &cst
	}

	vImageUrl := v.ImageUrl
	if vImageUrl != "" && !strings.HasPrefix(vImageUrl, "http") {
		vImageUrl = baseURL + vImageUrl
	}

	var images []*model.VariantImage
	for _, img := range v.Images {
		createdAt, _ := time.Parse(time.RFC3339, img.CreatedAt)
		
		imgUrl := img.Url
		if imgUrl != "" && !strings.HasPrefix(imgUrl, "http") {
			imgUrl = baseURL + imgUrl
		}

		altText := img.AltText
		var altTextPtr *string
		if altText != "" {
			altTextPtr = &altText
		}

		images = append(images, &model.VariantImage{
			ID:        img.Id,
			VariantID: img.VariantId,
			URL:       imgUrl,
			AltText:   altTextPtr,
			SortOrder: int(img.SortOrder),
			CreatedAt: createdAt,
		})
	}

	if len(images) > 0 {
		println("API GATEWAY mapVariantFromProto mapped images:", len(images), "ID:", images[0].ID, "URL:", images[0].URL)
	}

	// Fallback to first image from Images array if ImageURL is empty
	if vImageUrl == "" && len(images) > 0 {
		vImageUrl = images[0].URL
	}
	
	var vImageUrlPtr *string
	if vImageUrl != "" {
		vImageUrlPtr = &vImageUrl
	}

	return &model.Variant{
		ID:             v.Id,
		ProductID:      v.ProductId,
		Sku:            v.Sku,
		Name:           v.Name,
		Options:        optionsPtr,
		Price:          v.Price,
		CompareAtPrice: cmpPricePtr,
		CostPrice:      costPricePtr,
		WeightGrams:    int(v.WeightGrams),
		ImageURL:       vImageUrlPtr,
		IsActive:       v.IsActive,
		CreatedAt:      time.Unix(v.CreatedAt, 0),
		UpdatedAt:      time.Unix(v.UpdatedAt, 0),
		Images:         images,
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

func floatValue(f *float64) float64 {
	if f != nil {
		return *f
	}
	return 0
}

func intValue(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}
