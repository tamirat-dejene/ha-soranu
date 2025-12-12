package handler

import (
	"context"

	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/api/grpc/dto"
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
	errs "github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain/err"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/userpb"
	"google.golang.org/grpc"
)

type userHandler struct {
	userpb.UnimplementedUserServiceServer
	userUsecase domain.UserUseCase
}

// AddAddress implements userpb.UserServiceServer.
func (u *userHandler) AddAddress(ctx context.Context, req *userpb.AddAddressRequest) (*userpb.AddAddressResponse, error) {
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	address, err := u.userUsecase.AddAddress(ctx, req.UserId, dto.ToDomainAddress(req))
	if err != nil {
		return nil, err
	}

	return &userpb.AddAddressResponse{
		Address: dto.ToProtoAddress(*address),
	}, nil	

}

// AddPhoneNumber implements userpb.UserServiceServer.
func (u *userHandler) AddPhoneNumber(ctx context.Context, req *userpb.AddPhoneNumberRequest) (*userpb.MessageResponse, error) {
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	err := u.userUsecase.AddPhoneNumber(ctx, req.UserId, req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	return &userpb.MessageResponse{
		Message: "Phone number added successfully",
	}, nil
}

// GetAddresses implements userpb.UserServiceServer.
func (u *userHandler) GetAddresses(ctx context.Context, req *userpb.GetAddressesRequest) (*userpb.GetAddressesResponse, error) {
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	addresses, err := u.userUsecase.GetAddresses(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	var protoAddresses []*userpb.Address
	for _, addr := range addresses {
		protoAddresses = append(protoAddresses, dto.ToProtoAddress(addr))
	}

	return &userpb.GetAddressesResponse{
		Addresses: protoAddresses,
	}, nil
}

// GetPhoneNumber implements userpb.UserServiceServer.
func (u *userHandler) GetPhoneNumber(ctx context.Context, req *userpb.GetPhoneNumberRequest) (*userpb.GetPhoneNumberResponse, error) {
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	user, err := u.userUsecase.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &userpb.GetPhoneNumberResponse{
		PhoneNumber: user.PhoneNumber,
	}, nil
}

// GetUser implements userpb.UserServiceServer.
func (u *userHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	user, err := u.userUsecase.GetUser(ctx, req.UserId)
	if err != nil || user == nil {
		if err == nil {
			err = errs.ErrUserNotFound
		}
		return nil, err
	}

	return &userpb.GetUserResponse{
		User: dto.ToProtoUser(user),
	}, nil
}

// RemoveAddress implements userpb.UserServiceServer.
func (u *userHandler) RemoveAddress(ctx context.Context, req *userpb.RemoveAddressRequest) (*userpb.MessageResponse, error) {
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	err := u.userUsecase.RemoveAddress(ctx, req.UserId, req.AddressId)
	if err != nil {
		return nil, err
	}

	return &userpb.MessageResponse{
		Message: "Address removed successfully",
	}, nil
}

// RemovePhoneNumber implements userpb.UserServiceServer.
func (u *userHandler) RemovePhoneNumber(ctx context.Context, req *userpb.RemovePhoneNumberRequest) (*userpb.MessageResponse, error) {
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	err := u.userUsecase.RemovePhoneNumber(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &userpb.MessageResponse{
		Message: "Phone number removed successfully",
	}, nil
}

// UpdatePhoneNumber implements userpb.UserServiceServer.
func (u *userHandler) UpdatePhoneNumber(ctx context.Context, req *userpb.UpdatePhoneNumberRequest) (*userpb.MessageResponse, error) {
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	err := u.userUsecase.UpdatePhoneNumber(ctx, req.UserId, req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	return &userpb.MessageResponse{
		Message: "Phone number updated successfully",
	}, nil
}

func NewGrpcUserHandler(s *grpc.Server, usecase domain.UserUseCase) {
	handler := &userHandler{userUsecase: usecase}
	userpb.RegisterUserServiceServer(s, handler)
}
