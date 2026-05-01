package main

import (
	"context"
	"log"
	"net"

	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/config"

	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/handler"

	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/repository"

	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/service"

	// "ecommerce/gen/userpb"
	userpb "github.com/manojnegi/ecomm-microservices/gen/go/user/v1"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	// Run migrations (simplified — use golang-migrate in production)
	if err := runMigrations(pool); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	// Layer wiring: repo → service → handler
	userRepo := repository.NewUserRepo(pool)
	addrRepo := repository.NewAddressRepo(pool)
	tokenRepo := repository.NewTokenRepo(pool)

	authSvc := service.NewAuthService(userRepo, tokenRepo, cfg.PASETO_SECRET)
	userSvc := service.NewUserService(userRepo)
	addrSvc := service.NewAddressService(addrRepo)

	authHandler := handler.NewAuthHandler(authSvc, userSvc)
	addrHandler := handler.NewAddressHandler(addrSvc)
	// gRPC server
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(handler.AuthInterceptor(authSvc)),
	)

	// Register all service interfaces
	userpb.RegisterAuthServiceServer(srv, authHandler)
	userpb.RegisterUserServiceServer(srv, authHandler)
	userpb.RegisterAddressServiceServer(srv, addrHandler)

	reflection.Register(srv)

	log.Printf("User Service running on :%s", cfg.Port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runMigrations(pool *pgxpool.Pool) error {
	// In production, use golang-migrate or atlas
	// For now, ensure tables exist
	_, err := pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            email VARCHAR(255) UNIQUE NOT NULL,
            password_hash TEXT NOT NULL,
            first_name VARCHAR(100) NOT NULL,
            last_name VARCHAR(100) NOT NULL,
            phone VARCHAR(20),
            avatar_url TEXT,
            status VARCHAR(20) DEFAULT 'active',
            is_email_verified BOOLEAN DEFAULT false,
            email_verified_at TIMESTAMPTZ,
            last_login_at TIMESTAMPTZ,
            failed_login_count SMALLINT DEFAULT 0,
            locked_until TIMESTAMPTZ,
            created_at TIMESTAMPTZ DEFAULT NOW(),
            updated_at TIMESTAMPTZ DEFAULT NOW(),
            deleted_at TIMESTAMPTZ
        );

        CREATE TABLE IF NOT EXISTS addresses (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            user_id UUID NOT NULL REFERENCES users(id),
            label VARCHAR(50) NOT NULL,
            full_name VARCHAR(200) NOT NULL,
            phone VARCHAR(20) NOT NULL,
            line1 VARCHAR(255) NOT NULL,
            line2 VARCHAR(255),
            city VARCHAR(100) NOT NULL,
            state VARCHAR(100) NOT NULL,
            postal_code VARCHAR(20) NOT NULL,
            country_code CHAR(2) NOT NULL,
            is_default BOOLEAN DEFAULT false,
            created_at TIMESTAMPTZ DEFAULT NOW()
        );

        CREATE TABLE IF NOT EXISTS refresh_tokens (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            user_id UUID NOT NULL REFERENCES users(id),
            token_hash VARCHAR(255) UNIQUE NOT NULL,
            device_info JSONB,
            expires_at TIMESTAMPTZ NOT NULL,
            revoked_at TIMESTAMPTZ,
            created_at TIMESTAMPTZ DEFAULT NOW()
        );

        CREATE TABLE IF NOT EXISTS roles (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            name VARCHAR(50) UNIQUE NOT NULL,
            description TEXT,
            created_at TIMESTAMPTZ DEFAULT NOW()
        );

        CREATE TABLE IF NOT EXISTS user_roles (
            user_id UUID REFERENCES users(id),
            role_id UUID REFERENCES roles(id),
            assigned_at TIMESTAMPTZ DEFAULT NOW(),
            PRIMARY KEY (user_id, role_id)
        );

        CREATE INDEX IF NOT EXISTS idx_addresses_user_id ON addresses(user_id);
        CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
        CREATE INDEX IF NOT EXISTS idx_refresh_tokens_hash ON refresh_tokens(token_hash);
    `)
	return err
}
