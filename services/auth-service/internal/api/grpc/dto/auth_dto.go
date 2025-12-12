package dto

import (
    "github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
    "github.com/tamirat-dejene/ha-soranu/shared/protos/authpb"
    "github.com/tamirat-dejene/ha-soranu/shared/protos/userpb"
    "google.golang.org/protobuf/types/known/timestamppb"
)

func ToDomainUserRegister(req *authpb.UserRegisterRequest) domain.UserRegister {
    return domain.UserRegister{
        Email:       req.Email,
        Password:    req.Password,
        Username:    req.Username,
        PhoneNumber: req.PhoneNumber,
    }
}

func ToProtoUser(user *domain.User) *userpb.User {
    addresses := make([]*userpb.Address, len(user.Addresses))
    for i, a := range user.Addresses {
        addresses[i] = &userpb.Address{
            AddressId:  a.AddressID,
            Street:     a.Street,
            City:       a.City,
            State:      a.State,
            PostalCode: a.PostalCode,
            Country:    a.Country,
            Latitude:   a.Latitude,
            Longitude:  a.Longitude,
        }
    }

    return &userpb.User{
        UserId:      user.UserID,
        Email:       user.Email,
        Username:    user.Username,
        PhoneNumber: user.PhoneNumber,
        Addresses:   addresses,
        CreatedAt:   timestamppb.New(user.CreatedAt),
    }
}

func ToProtoTokens(t *domain.AuthTokens) *authpb.AuthTokens {
    return &authpb.AuthTokens{
        AccessToken:  t.AccessToken,
        RefreshToken: t.RefreshToken,
    }
}

func ToDomainLoginWithGoogle(req *authpb.GLoginRequest) domain.LoginWithGoogle {
    return domain.LoginWithGoogle{
        IDToken: req.IdToken,
    }
}

func ToDomainLoginWithEmail(req *authpb.EPLoginRequest) domain.LoginWithEmail {
    return domain.LoginWithEmail{
        Email:    req.Email,
        Password: req.Password,
    }
}