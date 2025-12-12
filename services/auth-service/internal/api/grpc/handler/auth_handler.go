package handler

import (
	"context"

	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/api/grpc/dto"
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
	errs "github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain/err"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/authpb"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/userpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type authHandler struct {
	authpb.UnimplementedAuthServiceServer
	usecase domain.AuthUseCase
}

// LoginWithEmailAndPassword implements authpb.AuthServiceServer.
func (a *authHandler) LoginWithEmailAndPassword(ctx context.Context, req *authpb.EPLoginRequest) (*authpb.LoginResponse, error) {
	logger.Info("Received email-password login request")
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	user, tokens, err := a.usecase.LoginWithEmailAndPassword(ctx, dto.ToDomainLoginWithEmail(req))
	if err != nil {
		logger.Error("Failed to login with email and password", zap.String("email", req.Email), zap.Error(err))
		return nil, err
	}

	logger.Info("User logged in successfully", zap.String("user_id", user.UserID), zap.String("email", user.Email))

	return &authpb.LoginResponse{
		User:   dto.ToProtoUser(user),
		Tokens: dto.ToProtoTokens(tokens),
	}, nil
}

// LoginWithGoogle implements authpb.AuthServiceServer.
func (a *authHandler) LoginWithGoogle(ctx context.Context, req *authpb.GLoginRequest) (*authpb.LoginResponse, error) {
	logger.Info("Received Google login request")
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	user, tokens, err := a.usecase.LoginWithGoogle(ctx, dto.ToDomainLoginWithGoogle(req))
	if err != nil {
		logger.Error("Failed to login with Google", zap.String("id_token", req.IdToken), zap.Error(err))
		return nil, err
	}

	logger.Info("User logged in with Google successfully", zap.String("user_id", user.UserID), zap.String("email", user.Email))

	return &authpb.LoginResponse{
		User:   dto.ToProtoUser(user),
		Tokens: dto.ToProtoTokens(tokens),
	}, nil
}

// Logout implements authpb.AuthServiceServer.
func (a *authHandler) Logout(ctx context.Context, req *authpb.LogoutRequest) (*userpb.MessageResponse, error) {
	logger.Info("Received logout request")
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	err := a.usecase.Logout(ctx, req.RefreshToken)
	if err != nil {
		logger.Error("Failed to logout user", zap.String("refresh_token", req.RefreshToken), zap.Error(err))
		return nil, err
	}

	logger.Info("User logged out successfully", zap.String("refresh_token", req.RefreshToken))
	return &userpb.MessageResponse{
		Message: "Logout successful",
	}, nil
	
}

// Refresh implements authpb.AuthServiceServer.
func (a *authHandler) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.RefreshResponse, error) {
	logger.Info("Received token refresh request")
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	tokens, err := a.usecase.RefreshTokens(ctx, req.RefreshToken)
	if  err != nil {
		logger.Error("Failed to refresh auth tokens", zap.Error(err))
		return nil, err
	}

	logger.Info("Auth tokens refreshed successfully")
	return &authpb.RefreshResponse{
		Tokens: dto.ToProtoTokens(tokens),
	}, nil
	
}

// Register implements authpb.AuthServiceServer.
func (a *authHandler) Register(ctx context.Context, req *authpb.UserRegisterRequest) (*authpb.UserRegisterResponse, error) {
	logger.Info("Received user registration request")
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	req_dto := dto.ToDomainUserRegister(req)
	user, authtoken, err :=  a.usecase.Register(ctx, req_dto)

	if err != nil {
		logger.Error("Failed to register user", zap.String("email", req_dto.Email), zap.Error(err))
		return nil, err
	}

	logger.Info("User registered successfully", zap.String("user_id", user.UserID), zap.String("email", req_dto.Email))

	return &authpb.UserRegisterResponse{
		User:   dto.ToProtoUser(user),
		Tokens: dto.ToProtoTokens(authtoken),
	}, nil
}

func NewGrpcAuthHandler(s *grpc.Server, usecase domain.AuthUseCase) {
	handler := &authHandler{usecase: usecase}
	authpb.RegisterAuthServiceServer(s, handler)
}
