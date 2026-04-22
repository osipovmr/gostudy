package facade

import (
	"context"
	"errors"
	"gostudy/internal/model/dto"
	"gostudy/internal/service"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type authFacade struct {
	userService  service.UserService
	tokenService service.TokenService
}

type AuthFacade interface {
	Register(context context.Context, req dto.RegisterRequest) (*dto.UserDto, error)
	Login(context context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
}

func NewAuthFacade(userService service.UserService, tokenService service.TokenService) AuthFacade {
	return &authFacade{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (f *authFacade) Register(context context.Context, req dto.RegisterRequest) (*dto.UserDto, error) {
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	userDto, err := f.userService.Create(context, dto.CreateUserInput{Email: req.Email, Name: req.Name, Password: hashedPassword})
	if err != nil {
		return nil, err
	}
	return userDto, nil
}

func (f *authFacade) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := f.userService.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if err := checkPassword(user.Password, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := f.tokenService.GenerateAccessToken(ctx, user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := f.tokenService.GenerateRefreshToken(ctx, user)
	if err != nil {
		return nil, err
	}
	err = f.tokenService.Save(ctx, user, refreshToken)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
