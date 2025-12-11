package handler

import (
	"context"

	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/api/grpc/dto"
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
	constants "github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain/const"
	errs "github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain/err"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/userpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type userHandler struct {
	userpb.UnimplementedUserServiceServer
	userUsecase domain.UserUsecase
}

// AddUserAddress implements userpb.UserServiceServer.
func (u *userHandler) AddUserAddress(ctx context.Context, req *userpb.AddUserAddressRequest) (*userpb.AddUserAddressResponse, error) {
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}
	// address := dto.AddUserAddressRequestFromProto(req)
	err := u.userUsecase.AddUserAddress(ctx, req.GetUserId(), req.GetAddress())
	if err != nil {
		logger.Error("Failed to add user address", zap.String("user_id", req.GetUserId()), zap.Error(err))
		return nil, errs.ToGRPCError(err)
	}
	logger.Info("User address added successfully", zap.String("user_id", req.GetUserId()))
	return &userpb.AddUserAddressResponse{
		Message: constants.AddressAddSuccessMessage,
	}, nil
}

// GetUser implements userpb.UserServiceServer.
func (u *userHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	user, err := u.userUsecase.GetUserByID(ctx, req.GetUserId())
	if err != nil {
		logger.Error("Failed to get user", zap.String("user_id", req.GetUserId()), zap.Error(err))
		return nil, errs.ToGRPCError(err)
	}

	return &userpb.GetUserResponse{
		User: dto.UserResponseToProto(domain.User{
			ID:         user.ID,
			Username:   user.Username,
			Email:      user.Email,
			Addressess: user.Addressess,
		}),
	}, nil
}

// GetUserAddresses implements userpb.UserServiceServer.
func (u *userHandler) GetUserAddresses(ctx context.Context, req *userpb.GetUserAddressesRequest) (*userpb.GetUserAddressesResponse, error) {
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	addresses, err := u.userUsecase.GetUserAddresses(ctx, req.GetUserId())
	if err != nil {
		logger.Error("Failed to get user addresses", zap.String("user_id", req.GetUserId()), zap.Error(err))
		return nil, errs.ToGRPCError(err)
	}

	var protoAddresses []*userpb.Address
	for _, addr := range addresses {
		protoAddresses = append(protoAddresses, &userpb.Address{Value: addr})
	}
	return &userpb.GetUserAddressesResponse{
		Addresses: protoAddresses,
	}, nil
}

// RemoveUserAddress implements userpb.UserServiceServer.
func (u *userHandler) RemoveUserAddress(ctx context.Context, req *userpb.RemoveUserAddressRequest) (*userpb.RemoveUserAddressResponse, error) {
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	err := u.userUsecase.RemoveUserAddress(ctx, req.GetUserId(), req.GetAddress())
	if err != nil {
		logger.Error("Failed to remove user address", zap.String("user_id", req.GetUserId()), zap.Error(err))
		return nil, errs.ToGRPCError(err)
	}
	logger.Info("User address removed successfully", zap.String("user_id", req.GetUserId()))

	return &userpb.RemoveUserAddressResponse{
		Message: constants.AddressRemoveSuccessMessage,
	}, nil
}

// UpdateUser implements userpb.UserServiceServer.
func (u *userHandler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	panic("unimplemented")
}

func NewGrpcUserHandler(s *grpc.Server, usecase domain.UserUsecase) {
	handler := &userHandler{userUsecase: usecase}
	userpb.RegisterUserServiceServer(s, handler)
}
