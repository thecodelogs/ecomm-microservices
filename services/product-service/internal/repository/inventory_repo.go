package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryRepo struct {
	db *pgxpool.Pool
}

func NewInventoryRepo(db *pgxpool.Pool) *InventoryRepo {
	return &InventoryRepo{db: db}
}

func (r *InventoryRepo) Create(ctx context.Context, inv *models.Inventory) error {
	query := `INSERT INTO inventory (id, variant_id, quantity_on_hand, quantity_reserved, reorder_point)
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, inv.ID, inv.VariantID, inv.QuantityOnHand, inv.QuantityReserved, inv.ReorderPoint)
	return err
}

func (r *InventoryRepo) GetByVariantID(ctx context.Context, variantID uuid.UUID) (*models.Inventory, error) {
	query := `SELECT id, variant_id, quantity_on_hand, quantity_reserved, quantity_available, reorder_point, updated_at
	          FROM inventory WHERE variant_id = $1`
	row := r.db.QueryRow(ctx, query, variantID)
	var i models.Inventory
	err := row.Scan(&i.ID, &i.VariantID, &i.QuantityOnHand, &i.QuantityReserved, &i.QuantityAvailable, &i.ReorderPoint, &i.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

// ReserveStock uses SELECT FOR UPDATE to prevent overselling
func (r *InventoryRepo) ReserveStock(ctx context.Context, variantID uuid.UUID, quantity int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Lock the row
	var available int
	query := `SELECT quantity_available FROM inventory WHERE variant_id = $1 FOR UPDATE`
	err = tx.QueryRow(ctx, query, variantID).Scan(&available)
	if err != nil {
		return fmt.Errorf("variant not found: %w", err)
	}

	if available < quantity {
		return errors.New("insufficient stock")
	}

	// Reserve
	update := `UPDATE inventory SET quantity_reserved = quantity_reserved + $1, updated_at = NOW() WHERE variant_id = $2`
	_, err = tx.Exec(ctx, update, quantity, variantID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *InventoryRepo) CommitStock(ctx context.Context, variantID uuid.UUID, quantity int) error {
	query := `UPDATE inventory 
	          SET quantity_on_hand = quantity_on_hand - $1,
	              quantity_reserved = quantity_reserved - $1,
	              updated_at = NOW()
	          WHERE variant_id = $2`
	_, err := r.db.Exec(ctx, query, quantity, variantID)
	return err
}

func (r *InventoryRepo) ReleaseStock(ctx context.Context, variantID uuid.UUID, quantity int) error {
	query := `UPDATE inventory 
	          SET quantity_reserved = quantity_reserved - $1,
	              updated_at = NOW()
	          WHERE variant_id = $2`
	_, err := r.db.Exec(ctx, query, quantity, variantID)
	return err
}
func (r *InventoryRepo) Update(ctx context.Context, inv *models.Inventory) error {
	query := `UPDATE inventory SET 
				quantity_on_hand = $1, 
				reorder_point = $2, 
				updated_at = NOW() 
			  WHERE variant_id = $3`
	_, err := r.db.Exec(ctx, query, inv.QuantityOnHand, inv.ReorderPoint, inv.VariantID)
	return err
}

func (r *InventoryRepo) DeleteByVariantID(ctx context.Context, variantID uuid.UUID) error {
	query := `DELETE FROM inventory WHERE variant_id = $1`
	_, err := r.db.Exec(ctx, query, variantID)
	return err
}
