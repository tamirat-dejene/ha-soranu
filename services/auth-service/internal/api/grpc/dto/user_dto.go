package dto

import (
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/userpb"
)

func UserResponseToProto(user domain.User) *userpb.User {
	addressesProto := make([]*userpb.Address, len(user.Addressess))
	for i, addr := range user.Addressess {
		addressesProto[i] = &userpb.Address{
			Value: addr,
		}
	}
	return &userpb.User{
		UserId:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Addresses: addressesProto,
	}
}