package service

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

type AuthService interface {
	// Authenticate проверяет login/pass, активность юзера
	// и возвращает подписанный JWT токен.
	Authenticate(ctx context.Context, username, password string) (token string, err error)

	// Authorize проверяет токен и превращает его в объект UserInfo
	// (тот самый, с методом InRole).
	Authorize(ctx context.Context, token string) (*utils.UserInfo, error)
}
