package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepo struct {
	db *pgxpool.Pool
}

func NewTokenRepo(db *pgxpool.Pool) *TokenRepo {
	return &TokenRepo{db: db}
}

func (r *TokenRepo) hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func (r *TokenRepo) Create(ctx context.Context, userID uuid.UUID, rawToken string, deviceInfo []byte, expiresAt time.Time) error {
	query := `
        INSERT INTO refresh_tokens (id, user_id, token_hash, device_info, expires_at, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := r.db.Exec(ctx, query,
		uuid.New(), userID, r.hashToken(rawToken), deviceInfo, expiresAt, time.Now().UTC(),
	)
	return err
}

func (r *TokenRepo) GetByHash(ctx context.Context, rawToken string) (*models.RefreshToken, error) {
	query := `
        SELECT id, user_id, token_hash, device_info, expires_at, revoked_at, created_at
        FROM refresh_tokens WHERE token_hash = $1
    `
	row := r.db.QueryRow(ctx, query, r.hashToken(rawToken))

	var t models.RefreshToken
	err := row.Scan(&t.ID, &t.UserID, &t.TokenHash, &t.DeviceInfo, &t.ExpiresAt, &t.RevokedAt, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TokenRepo) Revoke(ctx context.Context, tokenHash string) error {
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE token_hash = $2`
	_, err := r.db.Exec(ctx, query, time.Now().UTC(), tokenHash)
	return err
}

func (r *TokenRepo) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE user_id = $2 AND revoked_at IS NULL`
	_, err := r.db.Exec(ctx, query, time.Now().UTC(), userID)
	return err
}
