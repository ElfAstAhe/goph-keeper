package service

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
	"github.com/golang-jwt/jwt/v4"
)

// AuthService - сервис аутентификации и авторизации
type AuthService interface {
	// Authenticate проверяет login/pass
	// и возвращает подписанный JWT токен.
	Authenticate(ctx context.Context, username, encryptedPassword string) (*jwt.Token, error)

	// Authorize выдаёт набор ролей subject
	Authorize(ctx context.Context, subject string) (utils.Roles, error)
}
