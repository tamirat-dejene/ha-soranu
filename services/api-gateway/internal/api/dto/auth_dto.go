package dto

import "github.com/tamirat-dejene/ha-soranu/shared/protos/authpb"

type AuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterResponse struct {
	UserId string `json:"user_id"`
	Tokens AuthToken `json:"tokens"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

func RegisterResponseFromProto(res *authpb.RegisterResponse) *RegisterResponse {
	return &RegisterResponse{
		UserId: res.UserId,
		Tokens: AuthToken{
			AccessToken:  res.Tokens.AccessToken,
			RefreshToken: res.Tokens.RefreshToken,
		},
	}
}

func EPLoginResponseFromProto(res *authpb.EPLoginResponse) *LoginResponse {
	return &LoginResponse{
		AccessToken: res.AccessToken,
		RefreshToken: res.RefreshToken,
	}
}

func GLoginResponseFromProto(res *authpb.GLoginResponse) *LoginResponse {
	return &LoginResponse{
		AccessToken: res.AccessToken,
		RefreshToken: res.RefreshToken,
	}
}

func LogoutResponseFromProto(res *authpb.LogoutResponse) *LogoutResponse {
	return &LogoutResponse{
		Message: res.Message,
	}
}