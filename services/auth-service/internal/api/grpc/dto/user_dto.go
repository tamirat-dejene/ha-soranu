package dto

import (
	"github.com/tamirat-dejene/ha-soranu/services/auth-service/internal/domain"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/userpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToDomainAddress(p *userpb.AddAddressRequest) domain.Address {
	return domain.Address{
		Street:     p.Street,
		City:       p.City,
		State:      p.State,
		PostalCode: p.PostalCode,
		Country:    p.Country,
		Latitude:   p.Latitude,
		Longitude:  p.Longitude,
	}
}

func ToProtoAddress(a domain.Address) *userpb.Address {
	return &userpb.Address{
		AddressId:  a.AddressID,
		Street:     a.Street,
		City:       a.City,
		State:      a.State,
		PostalCode: a.PostalCode,
		Country:    a.Country,
		Latitude:   a.Latitude,
		Longitude:  a.Longitude,
		CreatedAt:  timestamppb.New(a.CreatedAt),
	}
}
