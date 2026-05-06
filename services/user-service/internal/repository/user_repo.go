package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (id, email, password_hash, first_name, last_name, phone, status, is_email_verified, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `
	_, err := r.db.Exec(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Phone, user.Status, user.IsEmailVerified, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
        SELECT id, email, password_hash, first_name, last_name, phone, avatar_url,
               status, is_email_verified, email_verified_at, last_login_at,
               failed_login_count, locked_until, created_at, updated_at, deleted_at
        FROM users WHERE email = $1 AND deleted_at IS NULL
    `
	row := r.db.QueryRow(ctx, query, email)

	var u models.User
	err := row.Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Phone, &u.AvatarURL,
		&u.Status, &u.IsEmailVerified, &u.EmailVerifiedAt, &u.LastLoginAt,
		&u.FailedLoginCount, &u.LockedUntil, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
			SELECT u.id, u.email, u.password_hash, u.first_name, u.last_name, u.phone, u.avatar_url,
				   u.status, u.is_email_verified, u.email_verified_at, u.last_login_at,
				   u.failed_login_count, u.locked_until, u.created_at, u.updated_at, u.deleted_at,
				   COALESCE(r.name, 'customer') as role_name
			FROM users u
			LEFT JOIN user_roles ur ON u.id = ur.user_id
			LEFT JOIN roles r ON ur.role_id = r.id
			WHERE u.id = $1 AND u.deleted_at IS NULL
		`
	row := r.db.QueryRow(ctx, query, id)

	var u models.User
	// Note: You will need to add a 'Role' field to your models.User struct
	err := row.Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Phone, &u.AvatarURL,
		&u.Status, &u.IsEmailVerified, &u.EmailVerifiedAt, &u.LastLoginAt,
		&u.FailedLoginCount, &u.LockedUntil, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
		&u.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Update(ctx context.Context, user *models.User) error {
	query := `
        UPDATE users SET first_name=$1, last_name=$2, phone=$3, avatar_url=$4,
                         status=$5, is_email_verified=$6, updated_at=$7
        WHERE id=$8 AND deleted_at IS NULL
    `
	_, err := r.db.Exec(ctx, query,
		user.FirstName, user.LastName, user.Phone, user.AvatarURL,
		user.Status, user.IsEmailVerified, time.Now().UTC(), user.ID,
	)
	return err
}

func (r *UserRepo) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET last_login_at=$1, failed_login_count=0, locked_until=NULL WHERE id=$2`
	_, err := r.db.Exec(ctx, query, time.Now().UTC(), id)
	return err
}

func (r *UserRepo) IncrementFailedLogin(ctx context.Context, id uuid.UUID) error {
	query := `
        UPDATE users SET failed_login_count = failed_login_count + 1,
                         locked_until = CASE WHEN failed_login_count >= 4 THEN $1 ELSE locked_until END
        WHERE id=$2
    `
	lockUntil := time.Now().UTC().Add(30 * time.Minute)
	_, err := r.db.Exec(ctx, query, lockUntil, id)
	return err
}

func (r *UserRepo) List(ctx context.Context, page, pageSize int32, status, search string) ([]models.User, int32, error) {
	// Build query with filters
	query := `SELECT id, email, first_name, last_name, phone, avatar_url, status, is_email_verified, created_at 
	          FROM users WHERE deleted_at IS NULL`
	args := []interface{}{}
	argCount := 1

	if status != "" {
		query += ` AND status = $` + strconv.Itoa(argCount)
		args = append(args, status)
		argCount++
	}
	if search != "" {
		query += ` AND (email ILIKE $` + strconv.Itoa(argCount) + ` OR first_name ILIKE $` + strconv.Itoa(argCount) + `)`
		args = append(args, "%"+search+"%")
		argCount++
	}

	// Count total
	countQuery := `SELECT COUNT(*) FROM (` + query + `) t`
	var total int32
	_ = r.db.QueryRow(ctx, countQuery, args...).Scan(&total)

	// Paginate
	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(argCount) + ` OFFSET $` + strconv.Itoa(argCount+1)
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Phone,
			&u.AvatarURL, &u.Status, &u.IsEmailVerified, &u.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepo) GetRoles(ctx context.Context, userID uuid.UUID) ([]string, error) {
	query := `
		SELECT r.name FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *UserRepo) UpdateStatus(ctx context.Context, userID uuid.UUID, status string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET status=$1, updated_at=$2 WHERE id=$3`, status, time.Now().UTC(), userID)
	return err
}

func (r *UserRepo) SoftDelete(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET deleted_at=$1 WHERE id=$2`, time.Now().UTC(), userID)
	return err
}
