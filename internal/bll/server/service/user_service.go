package service

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/model"
)

type UserService interface {
	// Register создает пользователя, хеширует пароль (SHA256) и
	// вызывает ChangeKeys для создания начальной пары ключей.
	Register(ctx context.Context, username, password, person, eMail string) (*model.User, error)

	// ChangeKeys генерирует новую пару ключей (RSA/Ed25519),
	// шифрует их мастер-паролем и сохраняет через репозиторий.
	ChangeKeys(ctx context.Context, userID string) (string, error)

	// GetProfile просто возвращает данные пользователя (ID, Username, Роли)
	GetProfile(ctx context.Context, userID string) (*model.User, error)

	// UpdatePassword — если решишь добавить логику перешифрования ключей при смене пароля
	UpdatePassword(ctx context.Context, userID, newPassword, oldPassword string) error
}
