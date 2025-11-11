package handler

import (
	"github.com/tamirat-dejene/ha-soranu/shared/proto/userpb"
	"github.com/tamirat-dejene/ha-soranu/user-service/internal/domain/entity"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	usecase entity.UserUsecase
}

func NewUserHandler(usecase entity.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: usecase,
	}
}
