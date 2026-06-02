package model

import "time"

type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Username  string     `json:"username"`
	Phone     *string    `json:"phone,omitempty"`
	Role      Role       `json:"role"`
	Status    UserStatus `json:"status"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type Product struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Sku         string     `json:"sku"`
	Brand       *string    `json:"brand,omitempty"`
	BrandID     *string    `json:"brandId,omitempty"`
	Stock       int        `json:"stock"`
	CategoryIDs []string   `json:"categoryIds"`
	Images      []string   `json:"images"`
	Variants    []*Variant `json:"variants"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type Category struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description,omitempty"`
	ImageURL    *string `json:"imageUrl,omitempty"`
	SortOrder   int     `json:"sortOrder"`
	IsActive    bool    `json:"isActive"`
	ParentID    *string `json:"parentId,omitempty"`
}

func (User) IsNode()     {}
func (Product) IsNode()  {}
func (Category) IsNode() {}

func (u User) GetID() string { return u.ID }
func (p Product) GetID() string { return p.ID }
func (c Category) GetID() string { return c.ID }
