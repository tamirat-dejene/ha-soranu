package handler

import (
	"context"

	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/api/grpc/dto"
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
	constants "github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain/const"
	errs "github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain/err"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/authpb"
	"google.golang.org/grpc"
)

type authHandler struct {
	authpb.UnimplementedAuthServiceServer
	usecase domain.AuthUsecase
}

func NewGrpcAuthHandler(s *grpc.Server, usecase domain.AuthUsecase) {
	handler := &authHandler{usecase: usecase}
	authpb.RegisterAuthServiceServer(s, handler)
}

func (a *authHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}
	// Convert proto request to domain credentials
	creds := dto.CreateLoginCredentialsFromProto(req)

	// Call usecase to login user
	token, err := a.usecase.Login(ctx, creds)
	if err != nil {
		return nil, err
	}

	// Convert domain response to proto response
	resp := dto.AuthTokenToProto(token)

	return resp, nil
}

// Logout implements authpb.AuthServiceServer.
func (a *authHandler) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	// Validate request
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	// Call usecase to logout user
	err := a.usecase.Logout(ctx, req.GetAccessToken())
	if err != nil {
		return nil, err
	}

	return dto.NewLogoutResponseProto(constants.LogoutSuccessMessage), nil
}

// Refresh implements authpb.AuthServiceServer.
func (a *authHandler) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.RefreshResponse, error) {
	// Validate request
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	// Convert proto request to domain refresh token
	refreshToken := req.GetRefreshToken()

	// Call usecase to refresh tokens
	newAccessToken, err := a.usecase.RefreshTokens(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	// Convert domain response to proto response
	return dto.AuthTokenToProtoRefresh(newAccessToken), nil
}

// Register implements authpb.AuthServiceServer.
func (a *authHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	// Validate request
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	// Convert proto request to domain request
	domainReq := dto.CreateUserRequestFromProto(req)

	// Call usecase to register user
	userID, token, err := a.usecase.Register(ctx, domainReq)
	if err != nil {
		return nil, err
	}

	// Convert domain response to proto response
	return dto.RegisterResponseToProto(userID, token), nil
}
