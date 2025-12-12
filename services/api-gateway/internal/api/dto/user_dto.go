package dto

import (
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/domain"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/userpb"
)

type GetUserRequestDTO struct {
	UserId string `json:"user_id" binding:"required"`
}

type GetUserResponseDTO struct {
	User domain.User `json:"user"`
}

func (gr *GetUserRequestDTO) ToProto() *userpb.GetUserRequest {
	return &userpb.GetUserRequest{
		UserId: gr.UserId,
	}
}

func toDomainAddresses(protoAddresses []*userpb.Address) []domain.Address {
	domainAddresses := make([]domain.Address, len(protoAddresses))
	for i, pa := range protoAddresses {
		domainAddresses[i] = domain.Address{
			ID:         pa.AddressId,
			Street:     pa.Street,
			City:       pa.City,
			State:      pa.State,
			PostalCode: pa.PostalCode,
			Country:    pa.Country,
			Latitude:   pa.Latitude,
			Longitude:  pa.Longitude,
		}
	}
	return domainAddresses
}

func GetUserResponseFromProto(protoRes *userpb.GetUserResponse) *GetUserResponseDTO {
	return &GetUserResponseDTO{
		User: domain.User{
			ID:      protoRes.User.UserId,
			Email:       protoRes.User.Email,
			Username:    protoRes.User.Username,
			PhoneNumber: protoRes.User.PhoneNumber,
			Password:   "********",
			Addresses:   toDomainAddresses(protoRes.User.Addresses),
			CreatedAt:   protoRes.User.CreatedAt.AsTime(),
		},
	}
}

type GetPhoneNumberRequestDTO struct {
	UserId string `json:"user_id" binding:"required"`
}

func (gpnr *GetPhoneNumberRequestDTO) ToProto() *userpb.GetPhoneNumberRequest {
	return &userpb.GetPhoneNumberRequest{
		UserId: gpnr.UserId,
	}
}

type GetPhoneNumberResponseDTO struct {
	PhoneNumber string `json:"phone_number"`
}

func GetPhoneNumberResponseFromProto(protoRes *userpb.GetPhoneNumberResponse) *GetPhoneNumberResponseDTO {
	return &GetPhoneNumberResponseDTO{
		PhoneNumber: protoRes.GetPhoneNumber(),
	}
}

type AddPhoneNumberRequestDTO struct {
	UserId      string `json:"user_id" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

func (apnr *AddPhoneNumberRequestDTO) ToProto() *userpb.AddPhoneNumberRequest {
	return &userpb.AddPhoneNumberRequest{
		UserId:      apnr.UserId,
		PhoneNumber: apnr.PhoneNumber,
	}
}

type UpdatePhoneNumberRequestDTO struct {
	UserId      string `json:"user_id" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

func (upnr *UpdatePhoneNumberRequestDTO) ToProto() *userpb.UpdatePhoneNumberRequest {
	return &userpb.UpdatePhoneNumberRequest{
		UserId:      upnr.UserId,
		PhoneNumber: upnr.PhoneNumber,
	}
}

type RemovePhoneNumberRequestDTO struct {
	UserId string `json:"user_id" binding:"required"`
}

func (rpnr *RemovePhoneNumberRequestDTO) ToProto() *userpb.RemovePhoneNumberRequest {
	return &userpb.RemovePhoneNumberRequest{
		UserId: rpnr.UserId,
	}
}

type GetAddressesRequestDTO struct {
	UserId string `json:"user_id" binding:"required"`
}

func (gar *GetAddressesRequestDTO) ToProto() *userpb.GetAddressesRequest {
	return &userpb.GetAddressesRequest{
		UserId: gar.UserId,
	}
}

type GetAddressesResponseDTO struct {
	Addresses []domain.Address `json:"addresses"`
}

func toDomainAddress(protoAddress *userpb.Address) domain.Address {
	return domain.Address{
		ID:         protoAddress.AddressId,
		Street:     protoAddress.Street,
		City:       protoAddress.City,
		State:      protoAddress.State,
		PostalCode: protoAddress.PostalCode,
		Country:    protoAddress.Country,
		Latitude:   protoAddress.Latitude,
		Longitude:  protoAddress.Longitude,
		CreatedAt:  protoAddress.CreatedAt.AsTime(),
	}
}

func GetAddressesResponseFromProto(protoRes *userpb.GetAddressesResponse) *GetAddressesResponseDTO {
	domainAddresses := make([]domain.Address, len(protoRes.Addresses))
	for i, pa := range protoRes.Addresses {
		domainAddresses[i] = toDomainAddress(pa)
	}
	return &GetAddressesResponseDTO{
		Addresses: domainAddresses,
	}
}

type AddAddressRequestDTO struct {
	UserId     string  `json:"user_id" binding:"required"`
	Street     string  `json:"street" binding:"required"`
	City       string  `json:"city" binding:"required"`
	State      string  `json:"state" binding:"required"`
	PostalCode uint32  `json:"postal_code" binding:"required"`
	Country    string  `json:"country" binding:"required"`
	Latitude   float32 `json:"latitude" binding:"required"`
	Longitude  float32 `json:"longitude" binding:"required"`
}

func (aar *AddAddressRequestDTO) ToProto() *userpb.AddAddressRequest {
	return &userpb.AddAddressRequest{
		UserId:     aar.UserId,
		Street:     aar.Street,
		City:       aar.City,
		State:      aar.State,
		PostalCode: aar.PostalCode,
		Country:    aar.Country,
		Latitude:   aar.Latitude,
		Longitude:  aar.Longitude,
	}
}

type RemoveAddressRequestDTO struct {
	UserId    string `json:"user_id" binding:"required"`
	AddressId string `json:"address_id" binding:"required"`
}

func (rar *RemoveAddressRequestDTO) ToProto() *userpb.RemoveAddressRequest {
	return &userpb.RemoveAddressRequest{
		UserId:    rar.UserId,
		AddressId: rar.AddressId,
	}
}

func AddAddressResponseFromProto(protoRes *userpb.AddAddressResponse) domain.Address {
	return toDomainAddress(protoRes.GetAddress())
}