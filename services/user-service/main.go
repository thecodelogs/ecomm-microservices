package main

import (
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	userpb "github.com/manojnegi/ecomm-microservices/gen/go/user/v1"
	userdb "github.com/manojnegi/ecomm-microservices/services/user-service/db"
	usergrpc "github.com/manojnegi/ecomm-microservices/services/user-service/grpc"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, using environment variables")
	}

	// Connect to database
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := userdb.Connect(databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	slog.Info("connected to database")

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register user service
	userServer := &usergrpc.UserServer{DB: db}
	userpb.RegisterUserServiceServer(grpcServer, userServer)

	// Enable server reflection (for evans / grpcurl)
	reflection.Register(grpcServer)

	// Start listening
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on port %s: %v", port, err)
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigCh
		slog.Info("received signal, shutting down", "signal", sig)
		grpcServer.GracefulStop()
	}()

	slog.Info("user-service gRPC server starting", "port", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
