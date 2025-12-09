package auth

import (
	"errors"
	"go/adv-demo/internal/user"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(userRepository *user.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (service *AuthService) Login(email, password string) (string, error) {
	user, err := service.UserRepository.FindByEmail(email)
	if err != nil {
		return "", errors.New(ErrWrongCredentials)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New(ErrWrongCredentials)
	}
	return user.Email, nil
}

func (service *AuthService) Register(email, password, name string) (string, error) {
	existingUser, _ := service.UserRepository.FindByEmail(email)
	if existingUser != nil {
		return "", errors.New(ErrUserAlreadyExists)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user := &user.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}
	_, err = service.UserRepository.Create(user)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}
