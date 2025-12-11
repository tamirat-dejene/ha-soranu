package handler

import (
	"context"

	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/api/grpc/dto"
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
	constants "github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain/const"
	errs "github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain/err"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/authpb"
	"go.uber.org/zap"
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

// LoginWithEmailAndPassword implements authpb.AuthServiceServer.
func (a *authHandler) LoginWithEmailAndPassword(ctx context.Context, req *authpb.EPLoginRequest) (*authpb.EPLoginResponse, error) {
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}
	// Convert proto request to domain credentials
	creds := dto.CreateLoginCredentialsFromProto(req)

	// Call usecase to login user
	token, err := a.usecase.LoginWithEP(ctx, creds)
	if err != nil {
		logger.Error("Login failed", zap.String("email", creds.Email), zap.Error(err))
		return nil, errs.ToGRPCError(err)
	}
	logger.Info("User logged in successfully", zap.String("email", creds.Email))

	// Convert domain response to proto response
	resp := dto.AuthTokenToProto(token)

	return resp, nil
}

// LoginWithGoogle implements authpb.AuthServiceServer.
func (a *authHandler) LoginWithGoogle(ctx context.Context, req *authpb.GLoginRequest) (*authpb.GLoginResponse, error) {
	if req == nil {
		return nil, errs.ErrInvalidRequest
	}

	// Convert proto request to domain Google token
	googleToken := req.GetGoogleToken()

	// Call usecase to login user with Google
	userInfo, token, err := a.usecase.LoginWithGoogle(ctx, googleToken)
	if err != nil {
		logger.Error("Google login failed", zap.Error(err))
		return nil, errs.ToGRPCError(err)
	}
	logger.Info("User logged in with Google successfully", zap.String("email", userInfo.Email))

	// Convert domain response to proto response
	resp := &authpb.GLoginResponse{
		Username:     userInfo.Username,
		Email:        userInfo.Email,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

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
		logger.Error("Logout failed", zap.Error(err))
		return nil, errs.ToGRPCError(err)
	}
	logger.Info("User logged out successfully")

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
	new_auth, err := a.usecase.RefreshTokens(ctx, refreshToken)
	if err != nil {
		logger.Error("Token refresh failed", zap.Error(err))
		return nil, errs.ToGRPCError(err)
	}

	// Convert domain response to proto response
	return dto.AuthTokenToProtoRefresh(new_auth.AccessToken, new_auth.RefreshToken), nil
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
		logger.Error("Registration failed", zap.String("email", domainReq.Email), zap.Error(err))
		return nil, errs.ToGRPCError(err)
	}
	logger.Info("User registered successfully", zap.String("user_id", userID), zap.String("email", domainReq.Email))

	// Convert domain response to proto response
	return dto.RegisterResponseToProto(userID, token), nil
}
