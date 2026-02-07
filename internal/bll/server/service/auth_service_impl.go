package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/model"
	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/repository"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
	"github.com/golang-jwt/jwt/v4"
)

// AuthServiceImpl - реализация сервиса аутентификации и авторизации
type AuthServiceImpl struct {
	keyCipher  utils.Cipher
	authHelper *utils.AuthHelper
	userRepo   repository.UserRepository
}

func NewAuthService(keyCipher utils.Cipher, authHelper *utils.AuthHelper, userRepo repository.UserRepository) *AuthServiceImpl {
	return &AuthServiceImpl{
		keyCipher:  keyCipher,
		authHelper: authHelper,
		userRepo:   userRepo,
	}
}

func (a *AuthServiceImpl) Authenticate(ctx context.Context, username, password string) (*jwt.Token, error) {
	// валидация
	if err := a.validateAuthenticate(username, password); err != nil {
		return nil, apperrs.NewBllValidateError("income", fmt.Sprintf("username [%s], password [censored]", username), "invalid income", err)
	}
	// подгружаем пользователя
	user, err := a.userRepo.GetByKey(ctx, model.NewUserKey(username))
	if err != nil {
		return nil, apperrs.NewBllCommonError("load user", err)
	}
	// готовим hash пароля
	passwordHash, err := a.keyCipher.EncryptString(password)
	if err != nil {
		return nil, apperrs.NewBllCommonError("error hash password", err)
	}
	// валидируем
	err = a.validateUserAndPassword(user, passwordHash)
	if err != nil {
		return nil, errs.NewAuthUnauthorizedError("invalid credential", err)
	}

	// роли
	roles, err := a.Authorize(ctx, username)
	if err != nil {
		return nil, apperrs.NewBllCommonError("authorize user", err)
	}

	// user info
	ui, err := a.userInfoFromUser(user, roles)
	if err != nil {
		return nil, apperrs.NewBllCommonError("build user info", err)
	}

	// token
	token, err := a.authHelper.TokenFromUserInfo(ui)
	if err != nil {
		return nil, apperrs.NewBllCommonError("token err", err)
	}

	return token, nil
}

func (a *AuthServiceImpl) Authorize(ctx context.Context, subject string) (utils.Roles, error) {
	if err := a.validateAuthorize(subject); err != nil {
		return nil, apperrs.NewBllValidateError("subject", "", "subject required", err)
	}
	if subject == "admin" {
		return utils.AdminRoles, nil
	}

	return utils.UserRoles, nil
}

func (a *AuthServiceImpl) validateAuthenticate(username, password string) error {
	if strings.TrimSpace(username) == "" {
		return apperrs.NewBllValidateError("username", "empty", "username is required", nil)
	}
	if strings.TrimSpace(password) == "" {
		return apperrs.NewBllValidateError("password", "empty", "password is required", nil)
	}

	return nil
}

func (a *AuthServiceImpl) validateAuthorize(subject string) error {
	if strings.TrimSpace(subject) == "" {
		return apperrs.NewBllValidateError("subject", "empty", "subject required", nil)
	}

	return nil
}

func (a *AuthServiceImpl) validateUserAndPassword(user *model.User, passwordHash string) error {
	// пользователь активен
	if !user.Active {
		return apperrs.NewBllValidateError("User.Active", fmt.Sprintf("%v", user.Active), "user inactive", nil)
	}
	// пользователя удалили (в текущей реализации репозитория это недостижимо, но всё меняется)
	if user.Deleted {
		return apperrs.NewBllValidateError("User.Deleted", fmt.Sprintf("%v", user.Deleted), "user is deleted", nil)
	}
	// пароли
	if user.PasswordHash != passwordHash {
		return apperrs.NewBllValidateError("PasswordHash", "incorrect", "password mismatch", nil)
	}

	return nil
}

func (a *AuthServiceImpl) userInfoFromUser(user *model.User, roles utils.Roles) (*utils.UserInfo, error) {
	if user == nil {
		return nil, apperrs.NewBllValidateError("user", "empty", "user is required", nil)
	}
	if len(roles) == 0 {
		return nil, apperrs.NewBllValidateError("roles", "empty", "roles is required", nil)
	}

	return utils.NewUserInfo(user.ID, user.Key.Username, roles.InRole("admin"), roles), nil
}
