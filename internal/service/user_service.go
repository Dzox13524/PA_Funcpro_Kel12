package service

import (
	"context"
	"errors"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type GetUserByIDFunc func(ctx context.Context, id string) (domain.User, error)
type CreateUserFunc func(ctx context.Context, name, email, password string) (domain.User, error)

func NewGetUserByID(getUserRepo repository.GetUserByIDRepoFunc) GetUserByIDFunc {
	return func(ctx context.Context, id string) (domain.User, error) {
		middleware.HandleLog("Service: Meminta user dengan ID " + id)
		return getUserRepo(ctx, id)
	}
}

func NewCreateUser(createRepo repository.CreateUserRepoFunc) CreateUserFunc {
	return func(ctx context.Context, name, email, password string) (domain.User, error) {
		if len(password) < 8 {
			return domain.User{}, errors.New("Password minimal 8 karakter")
		}

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return domain.User{}, err
		}

		newUser := domain.User{
			ID:       uuid.New().String(),
			Name:     name,
			Email:    email,
			Password: string(hashPassword),
			Role:     "User",
		}

		createdUser, err := createRepo(ctx, newUser)
		if err != nil {
			return domain.User{}, err
		}

		safeUser := createdUser
		safeUser.Password = ""
		
		return safeUser, nil
	}
}