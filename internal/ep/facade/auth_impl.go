package facade

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/service"
	"github.com/ElfAstAhe/goph-keeper/internal/ep/dto"
	"github.com/ElfAstAhe/goph-keeper/internal/ep/mapper"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
)

type AuthFacadeImpl struct {
	userService service.UserService
	authService service.AuthService
}

func NewAuthFacade(userService service.UserService, authService service.AuthService) *AuthFacadeImpl {
	return &AuthFacadeImpl{
		userService: userService,
		authService: authService,
	}
}

func (a *AuthFacadeImpl) Login(ctx context.Context, login *dto.LoginDto) (*dto.LoginResultDto, error) {

}

func (a *AuthFacadeImpl) Register(ctx context.Context, register *dto.RegisterDto) (*dto.RegisterResultDto, error) {
	user, err := a.userService.Register(ctx, register.Username, register.Password, register.Person, register.EMail)
	if err != nil {
		return nil, err
	}

	// результат
	res, err := mapper.UserToRegisterResultDto(user)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a *AuthFacadeImpl) validateRegister(register *dto.RegisterDto) error {
	return register.Validate()
}
