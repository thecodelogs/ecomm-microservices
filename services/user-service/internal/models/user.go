package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID      `db:"id"                json:"id"`
	Email            string         `db:"email"             json:"email"`
	PasswordHash     string         `db:"password_hash"     json:"-"`
	FirstName        string         `db:"first_name"        json:"first_name"`
	LastName         string         `db:"last_name"         json:"last_name"`
	Phone            sql.NullString `db:"phone"             json:"phone,omitempty"`
	AvatarURL        sql.NullString `db:"avatar_url"        json:"avatar_url,omitempty"`
	Status           string         `db:"status"            json:"status"`
	IsEmailVerified  bool           `db:"is_email_verified" json:"is_email_verified"`
	EmailVerifiedAt  sql.NullTime   `db:"email_verified_at" json:"email_verified_at,omitempty"`
	LastLoginAt      sql.NullTime   `db:"last_login_at"     json:"last_login_at,omitempty"`
	FailedLoginCount int16          `db:"failed_login_count" json:"-"`
	LockedUntil      sql.NullTime   `db:"locked_until"      json:"-"`
	CreatedAt        time.Time      `db:"created_at"        json:"created_at"`
	UpdatedAt        time.Time      `db:"updated_at"        json:"updated_at"`
	DeletedAt        sql.NullTime   `db:"deleted_at"        json:"-"`
}

type Address struct {
	ID          uuid.UUID      `db:"id"           json:"id"`
	UserID      uuid.UUID      `db:"user_id"      json:"user_id"`
	Label       string         `db:"label"        json:"label"`
	FullName    string         `db:"full_name"    json:"full_name"`
	Phone       string         `db:"phone"        json:"phone"`
	Line1       string         `db:"line1"        json:"line1"`
	Line2       sql.NullString `db:"line2"        json:"line2,omitempty"`
	City        string         `db:"city"         json:"city"`
	State       string         `db:"state"        json:"state"`
	PostalCode  string         `db:"postal_code"  json:"postal_code"`
	CountryCode string         `db:"country_code" json:"country_code"`
	IsDefault   bool           `db:"is_default"   json:"is_default"`
	CreatedAt   time.Time      `db:"created_at"   json:"created_at"`
}

type RefreshToken struct {
	ID         uuid.UUID      `db:"id"          json:"id"`
	UserID     uuid.UUID      `db:"user_id"     json:"user_id"`
	TokenHash  string         `db:"token_hash"  json:"-"`
	DeviceInfo sql.NullString `db:"device_info" json:"device_info,omitempty"`
	ExpiresAt  time.Time      `db:"expires_at"  json:"expires_at"`
	RevokedAt  sql.NullTime   `db:"revoked_at"  json:"revoked_at,omitempty"`
	CreatedAt  time.Time      `db:"created_at"  json:"created_at"`
}

type Role struct {
	ID          uuid.UUID `db:"id"          json:"id"`
	Name        string    `db:"name"        json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at"  json:"created_at"`
}
