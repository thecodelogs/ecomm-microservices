package service

import (
	"context"
	"database/sql"

	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/models"

	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/repository"

	"github.com/google/uuid"
)

type AddressService struct {
	addrRepo *repository.AddressRepo
}

func NewAddressService(addrRepo *repository.AddressRepo) *AddressService {
	return &AddressService{addrRepo: addrRepo}
}

func (s *AddressService) CreateAddress(ctx context.Context, userID uuid.UUID, label, fullName, phone, line1, line2, city, state, postalCode, countryCode string, isDefault bool) (*models.Address, error) {
	addr := &models.Address{
		ID:          uuid.New(),
		UserID:      userID,
		Label:       label,
		FullName:    fullName,
		Phone:       phone,
		Line1:       line1,
		Line2:       sql.NullString{String: line2, Valid: line2 != ""},
		City:        city,
		State:       state,
		PostalCode:  postalCode,
		CountryCode: countryCode,
		IsDefault:   isDefault,
	}

	if err := s.addrRepo.Create(ctx, addr); err != nil {
		return nil, err
	}

	if isDefault {
		_ = s.addrRepo.SetDefault(ctx, userID, addr.ID)
	}

	return addr, nil
}

func (s *AddressService) UpdateAddress(ctx context.Context, addressID, userID uuid.UUID, label, fullName, phone, line1, line2, city, state, postalCode, countryCode string, isDefault bool) (*models.Address, error) {
	addr, err := s.addrRepo.GetByID(ctx, addressID)
	if err != nil {
		return nil, err
	}

	// Ensure the user owns this address
	if addr.UserID != userID {
		return nil, sql.ErrNoRows // or a custom permission denied error
	}

	addr.Label = label
	addr.FullName = fullName
	addr.Phone = phone
	addr.Line1 = line1
	addr.Line2 = sql.NullString{String: line2, Valid: line2 != ""}
	addr.City = city
	addr.State = state
	addr.PostalCode = postalCode
	addr.CountryCode = countryCode
	addr.IsDefault = isDefault

	if err := s.addrRepo.Update(ctx, addr); err != nil {
		return nil, err
	}

	if isDefault {
		_ = s.addrRepo.SetDefault(ctx, userID, addr.ID)
	}

	return addr, nil
}

func (s *AddressService) ListAddresses(ctx context.Context, userID uuid.UUID) ([]models.Address, error) {
	return s.addrRepo.ListByUserID(ctx, userID)
}

func (s *AddressService) GetAddress(ctx context.Context, addressID uuid.UUID) (*models.Address, error) {
	return s.addrRepo.GetByID(ctx, addressID)
}

func (s *AddressService) SetDefaultAddress(ctx context.Context, userID, addressID uuid.UUID) error {
	return s.addrRepo.SetDefault(ctx, userID, addressID)
}

func (s *AddressService) DeleteAddress(ctx context.Context, addressID uuid.UUID) error {
	return s.addrRepo.Delete(ctx, addressID)
}
