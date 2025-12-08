package dto

import (
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/authpb"
)

func CreateUserRequestFromProto(protoReq *authpb.RegisterRequest) domain.CreateUserRequest {
	return domain.CreateUserRequest{
		Username: protoReq.GetUsername(),
		Email:    protoReq.GetEmail(),
		Password: protoReq.GetPassword(),
	}
}

func RegisterResponseToProto(userID string, token domain.AuthToken) *authpb.RegisterResponse {
	return &authpb.RegisterResponse{
		UserId: userID,
		Tokens: &authpb.EPLoginResponse{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		},
	}
}

func CreateLoginCredentialsFromProto(protoReq *authpb.EPLoginRequest) domain.LoginCredentials {
	return domain.LoginCredentials{
		Email:    protoReq.GetEmail(),
		Password: protoReq.GetPassword(),
	}
}

func AuthTokenToProto(token domain.AuthToken) *authpb.EPLoginResponse {
	return &authpb.EPLoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
}

func NewLogoutResponseProto(message string) *authpb.LogoutResponse {
	return &authpb.LogoutResponse{
		Message: message,
	}
}

func AuthTokenToProtoRefresh(accesstoken string, refreshtoken string) *authpb.RefreshResponse {
	return &authpb.RefreshResponse{
		AccessToken:  accesstoken,
		RefreshToken: refreshtoken,
	}
}