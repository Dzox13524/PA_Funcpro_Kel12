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

type UserServiceInterface interface {
	GetByID(ctx context.Context, id string) (domain.User, error)
	Create(ctx context.Context, name, email, password string) (domain.User, error)
}

type userService struct {
	userRepo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) UserServiceInterface {
	return &userService {
		userRepo: repo,
	}
}

func (s *userService) GetByID(ctx context.Context, id string) (domain.User, error) {
	middleware.HandleLog("Service: Meminta user dengan ID " + id)
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	middleware.HandleLog("Service: Meminta user dengan Email " + email)
	return  s.userRepo.GetUserByEmail(ctx, email)
}

func (s *userService) Create(ctx context.Context, name, email, password string) (domain.User, error) {
	if len(password) < 8 {
		return domain.User{}, errors.New("Password minimal 8 karakter")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}

	newUser := domain.User{
		ID : uuid.New().String(),
		Name: name,
		Email: email,
		Password: string(hashPassword),
		Role: "User",
	}

	createUser, err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		return domain.User{}, err
	}
	createUser.Password = ""
	return createUser, nil

}