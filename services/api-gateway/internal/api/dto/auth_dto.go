package dto

import (
	"github.com/tamirat-dejene/ha-soranu/shared/protos/authpb"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/userpb"
)

// UserRegisterRequestDTO represents data needed to register a new user.
type UserRegisterRequestDTO struct {
	Email       string  `json:"email" binding:"required,email"`
	Username    string  `json:"username" binding:"required"`
	PhoneNumber string  `json:"phone_number" binding:"required"`
	Password    string  `json:"password" binding:"required,min=4"`
}

// UserRegisterResponseDTO represents the response after registration.
type UserRegisterResponseDTO struct {
	User   *userpb.User `json:"user"`
	Tokens *authpb.AuthTokens `json:"tokens"`
}

// EPLoginRequestDTO for email/password login
type EPLoginRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LogoutRequestDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// GoogleLoginRequestDTO for Google OAuth login
type GoogleLoginRequestDTO struct {
	IdToken string `json:"id_token" binding:"required"`
}

// LoginResponseDTO
type LoginResponseDTO struct {
	User   *userpb.User
	Tokens *authpb.AuthTokens
}

type RefreshRequestDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (rr *RefreshRequestDTO) ToProto() *authpb.RefreshRequest {
	return &authpb.RefreshRequest{
		RefreshToken: rr.RefreshToken,
	}
}

func RefreshRequestFromProto(protoReq *authpb.RefreshRequest) *RefreshRequestDTO {
	return &RefreshRequestDTO{
		RefreshToken: protoReq.GetRefreshToken(),
	}
}

// RefreshResponseDTO
type RefreshResponseDTO struct {
	Tokens *authpb.AuthTokens
}

func RefreshResponseFromProto(protoResp *authpb.RefreshResponse) *RefreshResponseDTO {
	return &RefreshResponseDTO{
		Tokens: protoResp.GetTokens(),
	}
}

// AddressDTO represents a user address in DTO layer
type AddressDTO struct {
	ID         string
	UserID     string
	Street     string
	City       string
	State      string
	PostalCode uint32
	Country    string
	Latitude   float32
	Longitude  float32
}

func (ur UserRegisterRequestDTO) ToProto() *authpb.UserRegisterRequest {
    return &authpb.UserRegisterRequest{
        Email:       ur.Email,
        Username:    ur.Username,
        PhoneNumber: ur.PhoneNumber,
        Password:    ur.Password,
    }
}

func UserRegisterResponseFromProto(protoResp *authpb.UserRegisterResponse) *UserRegisterResponseDTO {
	return &UserRegisterResponseDTO{
		User:   protoResp.GetUser(),
		Tokens: protoResp.GetTokens(),
	}
}

func (lr *EPLoginRequestDTO) ToProto() *authpb.EPLoginRequest {
	return &authpb.EPLoginRequest{
		Email:    lr.Email,
		Password: lr.Password,
	}
}

func EPLoginRequestFromProto(protoReq *authpb.EPLoginRequest) *EPLoginRequestDTO {
	return &EPLoginRequestDTO{
		Email:    protoReq.GetEmail(),
		Password: protoReq.GetPassword(),
	}
}

func LoginResponseFromProto(protoResp *authpb.LoginResponse) *LoginResponseDTO {
	return &LoginResponseDTO{
		User:   protoResp.GetUser(),
		Tokens: protoResp.GetTokens(),
	}
}

func (gr *GoogleLoginRequestDTO) ToProto() *authpb.GLoginRequest {
	return &authpb.GLoginRequest{
		IdToken: gr.IdToken,
	}
}

func (lr *LogoutRequestDTO) ToProto() *authpb.LogoutRequest {
	return &authpb.LogoutRequest{
		RefreshToken: lr.RefreshToken,
	}
}