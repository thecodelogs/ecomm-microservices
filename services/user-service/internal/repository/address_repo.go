package repository

import (
	"context"
	"errors"

	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AddressRepo struct {
	db *pgxpool.Pool
}

func NewAddressRepo(db *pgxpool.Pool) *AddressRepo {
	return &AddressRepo{db: db}
}

func (r *AddressRepo) Create(ctx context.Context, addr *models.Address) error {
	query := `
        INSERT INTO addresses (id, user_id, label, full_name, phone, line1, line2, city, state, postal_code, country_code, is_default)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    `
	_, err := r.db.Exec(ctx, query,
		addr.ID, addr.UserID, addr.Label, addr.FullName, addr.Phone,
		addr.Line1, addr.Line2, addr.City, addr.State, addr.PostalCode, addr.CountryCode, addr.IsDefault,
	)
	return err
}

func (r *AddressRepo) Update(ctx context.Context, addr *models.Address) error {
	query := `
        UPDATE addresses 
        SET label = $1, full_name = $2, phone = $3, line1 = $4, line2 = $5, city = $6, state = $7, postal_code = $8, country_code = $9, is_default = $10
        WHERE id = $11 AND user_id = $12
    `
	_, err := r.db.Exec(ctx, query,
		addr.Label, addr.FullName, addr.Phone, addr.Line1, addr.Line2, addr.City, addr.State, addr.PostalCode, addr.CountryCode, addr.IsDefault,
		addr.ID, addr.UserID,
	)
	return err
}

func (r *AddressRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]models.Address, error) {
	query := `
        SELECT id, user_id, label, full_name, phone, line1, line2, city, state, postal_code, country_code, is_default, created_at
        FROM addresses WHERE user_id = $1 ORDER BY is_default DESC, created_at DESC
    `
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []models.Address
	for rows.Next() {
		var a models.Address
		err := rows.Scan(
			&a.ID, &a.UserID, &a.Label, &a.FullName, &a.Phone,
			&a.Line1, &a.Line2, &a.City, &a.State, &a.PostalCode, &a.CountryCode, &a.IsDefault, &a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}
	return addresses, rows.Err()
}

func (r *AddressRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Address, error) {
	query := `
        SELECT id, user_id, label, full_name, phone, line1, line2, city, state, postal_code, country_code, is_default, created_at
        FROM addresses WHERE id = $1
    `
	row := r.db.QueryRow(ctx, query, id)

	var a models.Address
	err := row.Scan(
		&a.ID, &a.UserID, &a.Label, &a.FullName, &a.Phone,
		&a.Line1, &a.Line2, &a.City, &a.State, &a.PostalCode, &a.CountryCode, &a.IsDefault, &a.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("address not found")
		}
		return nil, err
	}
	return &a, nil
}

func (r *AddressRepo) SetDefault(ctx context.Context, userID, addressID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Unset all defaults for user
	_, err = tx.Exec(ctx, `UPDATE addresses SET is_default = false WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	// Set new default
	_, err = tx.Exec(ctx, `UPDATE addresses SET is_default = true WHERE id = $1 AND user_id = $2`, addressID, userID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *AddressRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM addresses WHERE id = $1`, id)
	return err
}
