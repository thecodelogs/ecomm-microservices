package models

import "github.com/google/uuid"

type UserRole struct {
	UserID     uuid.UUID `db:"user_id"                json:"user_id"`
	RoleID     uuid.UUID `db:"role_id"             json:"role_id"`
	AssignedAT string    `db:"assigned_at"             json:"assigned_at"`
}
