package facade

import (
	"context"
	"errors"
	"gostudy/internal/kafka"
	"gostudy/internal/model/dto"
	"gostudy/internal/service"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type authFacade struct {
	userService  service.UserService
	tokenService service.TokenService
	producer     kafka.Producer
}

type AuthFacade interface {
	Register(context context.Context, req dto.RegisterRequest) (*dto.UserDto, error)
	Login(context context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	GetMe(context context.Context, email string) (*dto.UserDto, error)
	Refresh(context context.Context, req dto.RefreshRequest) (*dto.LoginResponse, error)
	Logout(context context.Context, email string) error
}

func NewAuthFacade(
	userService service.UserService,
	tokenService service.TokenService,
	producer kafka.Producer,
) AuthFacade {
	return &authFacade{
		userService:  userService,
		tokenService: tokenService,
		producer:     producer,
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
	err = f.producer.SendRegistrationMessage(context, userDto.Email)
	if err != nil {
		slog.Error(err.Error())

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

func (f *authFacade) GetMe(context context.Context, email string) (*dto.UserDto, error) {
	user, err := f.userService.GetByEmail(context, email)
	if err != nil {
		return nil, err
	}
	return &dto.UserDto{
		Uuid:  user.Uuid,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (f *authFacade) Refresh(context context.Context, req dto.RefreshRequest) (*dto.LoginResponse, error) {
	claims, err := f.tokenService.ValidateRefreshToken(context, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	user, err := f.userService.GetByEmail(context, claims.Email)
	if err != nil {
		return nil, err
	}

	if err := f.tokenService.Delete(context, user); err != nil {
		return nil, err
	}

	accessToken, err := f.tokenService.GenerateAccessToken(context, user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := f.tokenService.GenerateRefreshToken(context, user)
	if err != nil {
		return nil, err
	}

	if err := f.tokenService.Save(context, user, refreshToken); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (f *authFacade) Logout(context context.Context, email string) error {
	user, err := f.userService.GetByEmail(context, email)
	if err != nil {
		return err
	}

	if err := f.tokenService.Delete(context, user); err != nil {
		return err
	}
	return nil
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
