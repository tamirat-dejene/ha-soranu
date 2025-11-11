package usecase

import (
	"context"

	"github.com/tamirat-dejene/ha-soranu/user-service/internal/domain/entity"
	"github.com/tamirat-dejene/ha-soranu/user-service/internal/domain/errs"
)

type userUsecase struct {
	repo entity.UserRepository
}

func NewUserUsecase(repo entity.UserRepository) entity.UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) CreateUser(ctx context.Context, user *entity.User) error {
	err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return errs.ErrFailedToCreateUser
	}
	return nil
}

func (u *userUsecase) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errs.ErrUserNotFound
	}
	return user, nil
}

func (u *userUsecase) UpdateUser(ctx context.Context, user *entity.User) error {
	err := u.repo.UpdateUser(ctx, user)
	if err != nil {
		return errs.ErrFailedToUpdateUser
	}
	return nil
}

func (u *userUsecase) DeleteUser(ctx context.Context, id string) error {
	err := u.repo.DeleteUser(ctx, id)
	if err != nil {
		return errs.ErrFailedToDeleteUser
	}
	return nil
}