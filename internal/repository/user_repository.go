package repository

import (
	"context"
	"errors"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	GetUserByID(ctx context.Context, id string) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	Create(ctx context.Context, user domain.User) (domain.User, error)
}


type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &userRepository{
		db:db,
	}
}

func (r*userRepository) GetUserByID(ctx context.Context, ID string) (domain.User, error) {
	var user domain.User
	result := r.db.WithContext(ctx).First(&user, "ID = ?", ID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.User{}, result.Error
		}
	}
	return user, nil
}

func (r*userRepository) GetUserByEmail(ctx context.Context, Email string) (domain.User, error) {
	var user domain.User
	result := r.db.WithContext(ctx).First(&user, "email = ?", Email)
	return user, result.Error
}

func (r *userRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	result := r.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
        return domain.User{}, result.Error
    }
	return  user, nil
}