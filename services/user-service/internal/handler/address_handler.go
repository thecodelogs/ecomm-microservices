package handler

import (
	"context"

	userpb "github.com/manojnegi/ecomm-microservices/gen/go/user/v1"

	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/service"

	"github.com/manojnegi/ecomm-microservices/services/user-service/internal/models"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AddressHandler struct {
	userpb.UnimplementedAddressServiceServer
	addrSvc *service.AddressService
}

func NewAddressHandler(addrSvc *service.AddressService) *AddressHandler {
	return &AddressHandler{addrSvc: addrSvc}
}

func (h *AddressHandler) CreateAddress(ctx context.Context, req *userpb.CreateAddressRequest) (*userpb.Address, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	addr, err := h.addrSvc.CreateAddress(ctx, userID, req.Label, req.FullName, req.Phone,
		req.Line1, req.Line2, req.City, req.State, req.PostalCode, req.CountryCode, req.IsDefault)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProtoAddress(addr), nil
}

func (h *AddressHandler) UpdateAddress(ctx context.Context, req *userpb.UpdateAddressRequest) (*userpb.Address, error) {
	addressID, err := uuid.Parse(req.AddressId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address id")
	}
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	addr, err := h.addrSvc.UpdateAddress(ctx, addressID, userID, req.Label, req.FullName, req.Phone,
		req.Line1, req.Line2, req.City, req.State, req.PostalCode, req.CountryCode, req.IsDefault)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProtoAddress(addr), nil
}

func (h *AddressHandler) ListAddresses(ctx context.Context, req *userpb.ListAddressesRequest) (*userpb.AddressList, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	addresses, err := h.addrSvc.ListAddresses(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var pbAddrs []*userpb.Address
	for _, a := range addresses {
		pbAddrs = append(pbAddrs, toProtoAddress(&a))
	}

	return &userpb.AddressList{Addresses: pbAddrs}, nil
}

func (h *AddressHandler) GetAddress(ctx context.Context, req *userpb.GetAddressRequest) (*userpb.Address, error) {
	addressID, err := uuid.Parse(req.AddressId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address id")
	}

	addr, err := h.addrSvc.GetAddress(ctx, addressID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return toProtoAddress(addr), nil
}

func (h *AddressHandler) SetDefaultAddress(ctx context.Context, req *userpb.SetDefaultAddressRequest) (*userpb.Empty, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}
	addressID, err := uuid.Parse(req.AddressId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address id")
	}

	if err := h.addrSvc.SetDefaultAddress(ctx, userID, addressID); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.Empty{}, nil
}

func (h *AddressHandler) DeleteAddress(ctx context.Context, req *userpb.DeleteAddressRequest) (*userpb.Empty, error) {
	addressID, err := uuid.Parse(req.AddressId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address id")
	}

	if err := h.addrSvc.DeleteAddress(ctx, addressID); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userpb.Empty{}, nil
}

func toProtoAddress(a *models.Address) *userpb.Address {
	return &userpb.Address{
		Id:          a.ID.String(),
		UserId:      a.UserID.String(),
		Label:       a.Label,
		FullName:    a.FullName,
		Phone:       a.Phone,
		Line1:       a.Line1,
		Line2:       a.Line2.String,
		City:        a.City,
		State:       a.State,
		PostalCode:  a.PostalCode,
		CountryCode: a.CountryCode,
		IsDefault:   a.IsDefault,
		CreatedAt:   a.CreatedAt.Unix(),
	}
}
