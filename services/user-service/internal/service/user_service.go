package service

import (
	"context"
	"time"

	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/repository"

	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/models"

	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repository.UserRepo
}

func NewUserService(userRepo *repository.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, firstName, lastName, phone string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Phone.String = phone
	user.Phone.Valid = phone != ""
	user.UpdatedAt = time.Now().UTC()

	return s.userRepo.Update(ctx, user)
}

func (s *UserService) ListUsers(ctx context.Context, page, pageSize int32, status, search string) ([]models.User, int32, error) {
	return s.userRepo.List(ctx, page, pageSize, status, search)
}

func (s *UserService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return s.userRepo.GetRoles(ctx, userID)
}

func (s *UserService) UpdateStatus(ctx context.Context, userID uuid.UUID, status string) error {
	return s.userRepo.UpdateStatus(ctx, userID, status)
}

func (s *UserService) SoftDelete(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.SoftDelete(ctx, userID)
}
