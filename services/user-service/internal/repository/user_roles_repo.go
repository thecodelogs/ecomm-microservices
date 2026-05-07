package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRolesRepo struct {
	db *pgxpool.Pool
}

func NewUserRolesRepo(db *pgxpool.Pool) *UserRolesRepo {
	return &UserRolesRepo{db: db}
}

func (r *UserRolesRepo) Create(ctx context.Context, user_role *models.UserRole) error {
	query := `
	INSERT INTO user_roles (user_id, role_id)
	VALUES ($1, $2)
`
	_, err := r.db.Exec(ctx, query,
		user_role.UserID, user_role.RoleID,
	)
	return err
}

func (r *UserRepo) GetRoleIDByName(ctx context.Context, roleName string) (uuid.UUID, error) {
	query := `
		SELECT id
		FROM roles
		WHERE name = $1
		LIMIT 1
	`

	var roleID uuid.UUID

	err := r.db.QueryRow(ctx, query, roleName).Scan(&roleID)
	if err != nil {
		return uuid.Nil, err
	}

	return roleID, nil
}
