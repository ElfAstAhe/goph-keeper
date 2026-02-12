package facade

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/api/rest/dto"
	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/service"
	"github.com/ElfAstAhe/goph-keeper/internal/ep/mapper"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

type AuthFacadeImpl struct {
	jwtHelper   *utils.JWTHelper
	userService service.UserService
	authService service.AuthService
}

func NewAuthFacadeImpl(jwtHelper *utils.JWTHelper, userService service.UserService, authService service.AuthService) *AuthFacadeImpl {
	return &AuthFacadeImpl{
		jwtHelper:   jwtHelper,
		userService: userService,
		authService: authService,
	}
}

func (a *AuthFacadeImpl) Login(ctx context.Context, login *dto.LoginDto) (*dto.LoginResultDto, error) {
	// валидируем
	if err := a.validateLogin(login); err != nil {
		return nil, errs.NewAuthUnauthorizedError("unauthorized", err)
	}
	// проводим аутентификацию и авторизацию
	token, err := a.authService.Authenticate(ctx, login.Username, login.Password)
	if err != nil {
		return nil, err
	}
	// получаем строку токена
	tokenString, err := a.jwtHelper.BuildTokenStr(token)
	if err != nil {
		return nil, err
	}

	// результат
	res := dto.NewLoginResultDto(tokenString, "")

	return res, nil
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

func (a *AuthFacadeImpl) validateLogin(login *dto.LoginDto) error {
	return login.Validate()
}
