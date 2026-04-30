package db

import (
	"context"
	"database/sql"
	"time"
)

// UserRow represents a row from the users table.
type UserRow struct {
	ID            string
	RoleID        int32
	Email         string
	PasswordHash  string
	FirstName     string
	LastName      string
	Phone         sql.NullString
	AvatarURL     sql.NullString
	EmailVerified bool
	PhoneVerified bool
	IsActive      bool
	LastLoginAt   sql.NullTime
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// RoleRow represents a row from the roles table.
type RoleRow struct {
	ID          int32
	Name        string
	Description sql.NullString
	IsDefault   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateUserParams holds the input for inserting a new user.
type CreateUserParams struct {
	RoleID       int32
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Phone        string
}

// CreateUser inserts a new user and returns the full row.
func CreateUser(ctx context.Context, db *sql.DB, p CreateUserParams) (*UserRow, error) {
	query := `
		INSERT INTO users (role_id, email, password_hash, first_name, last_name, phone)
		VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''))
		RETURNING id, role_id, email, password_hash, first_name, last_name,
		          phone, avatar_url, email_verified, phone_verified,
		          is_active, last_login_at, created_at, updated_at
	`

	row := db.QueryRowContext(ctx, query,
		p.RoleID, p.Email, p.PasswordHash, p.FirstName, p.LastName, p.Phone,
	)

	var u UserRow
	err := row.Scan(
		&u.ID, &u.RoleID, &u.Email, &u.PasswordHash,
		&u.FirstName, &u.LastName, &u.Phone, &u.AvatarURL,
		&u.EmailVerified, &u.PhoneVerified, &u.IsActive,
		&u.LastLoginAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetRoleByID fetches a role by its primary key.
func GetRoleByID(ctx context.Context, db *sql.DB, roleID int32) (*RoleRow, error) {
	query := `SELECT id, name, description, is_default, created_at, updated_at FROM roles WHERE id = $1`

	row := db.QueryRowContext(ctx, query, roleID)

	var r RoleRow
	err := row.Scan(&r.ID, &r.Name, &r.Description, &r.IsDefault, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
