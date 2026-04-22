package facade

import (
	"context"
	"errors"
	"gostudy/internal/model/dto"
	"gostudy/internal/service"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrRegistration = errors.New("Failed to register user")
	ErrUserExists   = errors.New("user already exists")
)

type authFacade struct {
	userService  service.UserService
	tokenService service.TokenService
}

type AuthFacade interface {
	Register(context context.Context, req dto.RegisterRequest) (*dto.UserDto, error)
}

func NewAuthFacade(userService service.UserService, tokenService service.TokenService) AuthFacade {
	return &authFacade{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (f *authFacade) Register(context context.Context, req dto.RegisterRequest) (*dto.UserDto, error) {
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	userDto, err := f.userService.Create(context, dto.CreateUserInput{Email: req.Email, Name: req.Name, Password: hashedPassword})
	if err != nil {
		return nil, err
	}
	return userDto, nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
