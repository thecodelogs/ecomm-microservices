package main

import (
	"context"
	"log"
	"net"

	// categorypb "github.com/manojnegi/ecomm-microservices/gen/go/category/v1"
	productpb "github.com/manojnegi/ecomm-microservices/gen/go/product/v1"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/service"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/repository"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/handler"

	"github.com/manojnegi/ecomm-microservices/services/product-service/internal/config"

	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect to product_db: %v", err)
	}
	defer pool.Close()

	// Run migrations
	if err := runMigrations(pool); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	// Repository layer
	catRepo := repository.NewCategoryRepo(pool)
	prodRepo := repository.NewProductRepo(pool)
	varRepo := repository.NewVariantRepo(pool)
	invRepo := repository.NewInventoryRepo(pool)
	revRepo := repository.NewReviewRepo(pool)

	// Service layer
	prodSvc := service.NewProductService(prodRepo, varRepo, catRepo)
	invSvc := service.NewInventoryService(invRepo)
	revSvc := service.NewReviewService(revRepo, prodRepo)
	catSvc := service.NewCategoryService(prodRepo, varRepo, catRepo)

	// Handler layer
	prodHandler := handler.NewProductHandler(prodSvc, invSvc, revSvc)
	catHandler := handler.NewCategoryHandler(catSvc)

	// gRPC server
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	productpb.RegisterProductServiceServer(srv, prodHandler)
	productpb.RegisterCategoryServiceServer(srv, catHandler)

	reflection.Register(srv)

	log.Printf("Product Service running on :%s", cfg.Port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//go:embed migrations/006_seed_data.sql
var seedSQL string

func runMigrations(pool *pgxpool.Pool) error {
	_, err := pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS categories (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			parent_id UUID REFERENCES categories(id),
			slug VARCHAR(200) UNIQUE NOT NULL,
			name VARCHAR(200) NOT NULL,
			description TEXT,
			image_url TEXT,
			sort_order INTEGER DEFAULT 0,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMPTZ DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS products (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			category_id UUID REFERENCES categories(id),
			slug VARCHAR(300) UNIQUE NOT NULL,
			name VARCHAR(500) NOT NULL,
			description TEXT,
			short_description TEXT,
			brand VARCHAR(200),
			tags TEXT[],
			attributes JSONB,
			status VARCHAR(20) DEFAULT 'draft',
			vendor_id UUID,
			avg_rating DECIMAL(3,2) DEFAULT 0,
			review_count INTEGER DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS product_images (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			product_id UUID REFERENCES products(id) ON DELETE CASCADE,
			url TEXT NOT NULL,
			alt_text VARCHAR(300),
			sort_order INTEGER DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS variants (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			product_id UUID REFERENCES products(id) ON DELETE CASCADE,
			sku VARCHAR(100) UNIQUE NOT NULL,
			name VARCHAR(200) NOT NULL,
			options JSONB,
			price DECIMAL(12,2) NOT NULL,
			compare_at_price DECIMAL(12,2),
			cost_price DECIMAL(12,2),
			weight_grams INTEGER,
			image_url TEXT,
			is_active BOOLEAN DEFAULT true,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS inventory (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			variant_id UUID UNIQUE REFERENCES variants(id) ON DELETE CASCADE,
			quantity_on_hand INTEGER NOT NULL DEFAULT 0 CHECK (quantity_on_hand >= 0),
			quantity_reserved INTEGER NOT NULL DEFAULT 0 CHECK (quantity_reserved >= 0),
			quantity_available INTEGER GENERATED ALWAYS AS (quantity_on_hand - quantity_reserved) STORED,
			reorder_point INTEGER DEFAULT 10,
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS reviews (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			product_id UUID REFERENCES products(id) ON DELETE CASCADE,
			user_id UUID NOT NULL,
			order_id UUID,
			rating SMALLINT CHECK (rating BETWEEN 1 AND 5),
			title VARCHAR(300),
			body TEXT,
			status VARCHAR(20) DEFAULT 'pending',
			created_at TIMESTAMPTZ DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id);
		CREATE INDEX IF NOT EXISTS idx_products_status ON products(status);
		CREATE INDEX IF NOT EXISTS idx_products_vendor ON products(vendor_id);
		CREATE INDEX IF NOT EXISTS idx_variants_product ON variants(product_id);
		CREATE INDEX IF NOT EXISTS idx_reviews_product ON reviews(product_id);
		CREATE INDEX IF NOT EXISTS idx_inventory_available ON inventory(quantity_available);
	`)
	if err != nil {
		return err
	}

	// Run seed data
	_, err = pool.Exec(context.Background(), seedSQL)
	return err
}
