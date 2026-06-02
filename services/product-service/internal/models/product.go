package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID      `db:"id"           json:"id"`
	ParentID    uuid.NullUUID  `db:"parent_id"    json:"parent_id,omitempty"`
	Slug        string         `db:"slug"         json:"slug"`
	Name        string         `db:"name"         json:"name"`
	Description sql.NullString `db:"description"  json:"description,omitempty"`
	ImageURL    sql.NullString `db:"image_url"    json:"image_url,omitempty"`
	SortOrder   int            `db:"sort_order"   json:"sort_order"`
	IsActive    bool           `db:"is_active"    json:"is_active"`
	CreatedAt   time.Time      `db:"created_at"   json:"created_at"`
}

type Product struct {
	ID               uuid.UUID       `db:"id"                json:"id"`
	CategoryID       uuid.UUID       `db:"category_id"       json:"category_id"`
	Slug             string          `db:"slug"              json:"slug"`
	Name             string          `db:"name"              json:"name"`
	Description      string          `db:"description"       json:"description"`
	ShortDescription sql.NullString  `db:"short_description" json:"short_description,omitempty"`
	Brand            sql.NullString  `db:"brand"             json:"brand,omitempty"`
	Tags             []string        `db:"tags"              json:"tags,omitempty"`
	Attributes       json.RawMessage `db:"attributes"        json:"attributes,omitempty"`
	Status           string          `db:"status"            json:"status"`
	VendorID         uuid.UUID       `db:"vendor_id"         json:"vendor_id,omitempty"`
	AvgRating        float32         `db:"avg_rating"        json:"avg_rating,omitempty"`
	ReviewCount      int             `db:"review_count"      json:"review_count,omitempty"`
	CreatedAt        time.Time       `db:"created_at"        json:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at"        json:"updated_at"`
	ImageUrl         []string        `db:"image_url"          json:"image_url,omitempty"`
}

type ProductImage struct {
	ID        uuid.UUID `db:"id"         json:"id"`
	ProductID uuid.UUID `db:"product_id" json:"product_id"`
	URL       string    `db:"url"        json:"url"`
	AltText   string    `db:"alt_text"   json:"alt_text"`
	SortOrder int       `db:"sort_order" json:"sort_order"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Variant struct {
	ID             uuid.UUID       `db:"id"               json:"id"`
	ProductID      uuid.UUID       `db:"product_id"       json:"product_id"`
	SKU            string          `db:"sku"              json:"sku"`
	Name           string          `db:"name"             json:"name"`
	Options        json.RawMessage `db:"options"          json:"options"`
	Price          int64           `db:"price"            json:"price"` // paise/cents
	CompareAtPrice sql.NullInt64   `db:"compare_at_price" json:"compare_at_price,omitempty"`
	CostPrice      sql.NullInt64   `db:"cost_price"       json:"cost_price,omitempty"`
	WeightGrams    int             `db:"weight_grams"     json:"weight_grams"`
	ImageURL       sql.NullString  `db:"image_url"        json:"image_url,omitempty"`
	IsActive       bool            `db:"is_active"        json:"is_active"`
	CreatedAt      time.Time       `db:"created_at"       json:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"       json:"updated_at"`
	Images         []VariantImage  `json:"images"`
	InitialStock   int             `json:"initial_stock"`
}

type VariantImage struct {
	ID        uuid.UUID `db:"id"         json:"id"`
	VariantID uuid.UUID `db:"variant_id" json:"variant_id"`
	URL       string    `db:"url"        json:"url"`
	AltText   string    `db:"alt_text"   json:"alt_text"`
	SortOrder int       `db:"sort_order" json:"sort_order"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Inventory struct {
	ID                uuid.UUID `db:"id"                   json:"id"`
	VariantID         uuid.UUID `db:"variant_id"           json:"variant_id"`
	QuantityOnHand    int       `db:"quantity_on_hand"     json:"quantity_on_hand"`
	QuantityReserved  int       `db:"quantity_reserved"    json:"quantity_reserved"`
	QuantityAvailable int       `db:"quantity_available"   json:"quantity_available"`
	ReorderPoint      int       `db:"reorder_point"        json:"reorder_point"`
	UpdatedAt         time.Time `db:"updated_at"           json:"updated_at"`
}

type Review struct {
	ID        uuid.UUID `db:"id"         json:"id"`
	ProductID uuid.UUID `db:"product_id" json:"product_id"`
	UserID    uuid.UUID `db:"user_id"    json:"user_id"`
	OrderID   uuid.UUID `db:"order_id"   json:"order_id,omitempty"`
	Rating    int16     `db:"rating"     json:"rating"`
	Title     string    `db:"title"      json:"title"`
	Body      string    `db:"body"       json:"body"`
	Status    string    `db:"status"     json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
