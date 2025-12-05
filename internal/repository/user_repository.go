package repository

import (
	"context"
	"errors"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"gorm.io/gorm"
)

type GetUserByIDRepoFunc func(ctx context.Context, id string) (domain.User, error)
type GetUserByEmailRepoFunc func(ctx context.Context, email string) (domain.User, error)
type CreateUserRepoFunc func(ctx context.Context, user domain.User) (domain.User, error)

func NewGetUserByIDRepository(db *gorm.DB) GetUserByIDRepoFunc {
	return func(ctx context.Context, id string) (domain.User, error) {
		var user domain.User
		result := db.WithContext(ctx).First(&user, "ID = ?", id)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return domain.User{}, result.Error
			}
			return domain.User{}, result.Error
		}
		return user, nil
	}
}

func NewGetUserByEmailRepository(db *gorm.DB) GetUserByEmailRepoFunc {
	return func(ctx context.Context, email string) (domain.User, error) {
		var user domain.User
		result := db.WithContext(ctx).First(&user, "email = ?", email)
		return user, result.Error
	}
}

func NewCreateUserRepository(db *gorm.DB) CreateUserRepoFunc {
	return func(ctx context.Context, user domain.User) (domain.User, error) {
		result := db.WithContext(ctx).Create(&user)
		if result.Error != nil {
			return domain.User{}, result.Error
		}
		return user, nil
	}
}