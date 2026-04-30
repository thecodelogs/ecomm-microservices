package grpc

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"

	userdb "github.com/manojnegi/ecomm-microservices/services/user-service/db"

	userpb "github.com/manojnegi/ecomm-microservices/gen/go/user/v1"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServer implements the UserService gRPC server.
type UserServer struct {
	userpb.UnimplementedUserServiceServer
	DB *sql.DB
}

func (s *UserServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	// --- Validation ---
	if strings.TrimSpace(req.GetFname()) == "" {
		return nil, status.Error(codes.InvalidArgument, "first name is required")
	}
	if strings.TrimSpace(req.GetLname()) == "" {
		return nil, status.Error(codes.InvalidArgument, "last name is required")
	}
	if strings.TrimSpace(req.GetEmail()) == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if strings.TrimSpace(req.GetPassword()) == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// --- Hash password ---
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return nil, status.Error(codes.Internal, "failed to process password")
	}

	// Default to customer role (id=2) if not provided
	roleID := req.GetRoleId()
	if roleID == 0 {
		roleID = 2
	}

	// --- Insert user ---
	userRow, err := userdb.CreateUser(ctx, s.DB, userdb.CreateUserParams{
		RoleID:       roleID,
		Email:        strings.TrimSpace(req.GetEmail()),
		PasswordHash: string(hashedPassword),
		FirstName:    strings.TrimSpace(req.GetFname()),
		LastName:     strings.TrimSpace(req.GetLname()),
		Phone:        strings.TrimSpace(req.GetPhone()),
	})
	if err != nil {
		// Check for unique constraint violation (duplicate email)
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return nil, status.Error(codes.AlreadyExists, "a user with this email already exists")
		}
		slog.Error("failed to create user", "error", err)
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	// --- Fetch role for response ---
	roleRow, err := userdb.GetRoleByID(ctx, s.DB, userRow.RoleID)
	if err != nil {
		slog.Error("failed to fetch role", "error", err, "role_id", userRow.RoleID)
		return nil, status.Error(codes.Internal, "user created but failed to fetch role")
	}

	// --- Build response ---
	return &userpb.CreateUserResponse{
		User: mapUserRowToProto(userRow, roleRow),
	}, nil
}

// mapUserRowToProto converts DB rows into the protobuf User message.
func mapUserRowToProto(u *userdb.UserRow, r *userdb.RoleRow) *userpb.User {
	user := &userpb.User{
		Id:            u.ID,
		Name:          u.FirstName + " " + u.LastName,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		PhoneVerified: u.PhoneVerified,
		IsActive:      u.IsActive,
		CreatedAt:     u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if u.Phone.Valid {
		user.Phone = u.Phone.String
	}

	if u.LastLoginAt.Valid {
		user.LastLoginAt = u.LastLoginAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	if r != nil {
		user.Role = &userpb.Role{
			Id:        r.ID,
			Name:      r.Name,
			IsDefault: r.IsDefault,
			CreatedAt: r.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: r.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
		if r.Description.Valid {
			user.Role.Description = r.Description.String
		}
	}

	return user
}
