package service

import (
	"context"
	"strings"

	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/model"
	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/repository"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
	"github.com/ElfAstAhe/goph-keeper/pkg/logger"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

type UserServiceImpl struct {
	userRepo   repository.UserRepository
	log        logger.Logger
	keyCipher  utils.Cipher
	keysHelper *utils.RSAKeysHelper
}

func NewUserService(userRepo repository.UserRepository, keysHelper *utils.RSAKeysHelper, log logger.Logger) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo:   userRepo,
		log:        log,
		keysHelper: keysHelper,
	}
}

func (us *UserServiceImpl) Register(ctx context.Context, username, password, person, eMail string) (*model.User, error) {
	// валидация входящих данных
	if err := us.validateRegister(username, password); err != nil {
		return nil, err
	}

	// генерируем ключи
	privateKey, publicKey, err := us.keysHelper.Generate()
	if err != nil {
		return nil, apperrs.NewBllCommonError("generate rsa keys pair", err)
	}

	// получаем hash сумму пароля
	passwordHash, err := us.keyCipher.EncryptString(password)
	if err != nil {
		return nil, apperrs.NewBllCommonError("generate password hash", err)
	}

	// заполняем данные
	inst := model.NewEmptyUser()

	inst.Key.Username = username
	inst.PasswordHash = passwordHash
	inst.PrivateKey = privateKey
	inst.PublicKey = publicKey
	inst.Person = person
	inst.EMail = eMail
	inst.Active = true

	// создаём
	user, err := us.userRepo.Create(ctx, inst)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) validateRegister(username, password string) error {
	if strings.TrimSpace(username) == "" {
		return errs.NewAppInvalidArgumentError("username", username)
	}
	if strings.TrimSpace(password) == "" {
		return errs.NewAppInvalidArgumentError("password", "password is empty")
	}

	return nil
}

func (us *UserServiceImpl) ChangeKeys(ctx context.Context, userID string) (string, error) {
	// подгружаем данные
	user, err := us.userRepo.Get(ctx, userID)
	if err != nil {
		return "", err
	}
	if err = us.validateUserExists(user); err != nil {
		return "", err
	}

	// генерируем пару ключей
	privateKey, publicKey, err := us.keysHelper.Generate()
	if err != nil {
		return "", apperrs.NewBllCommonError("generate rsa keys pair", err)
	}

	// заполняем данные
	user.PrivateKey = privateKey
	user.PublicKey = publicKey

	// сохраняем данные
	user, err = us.userRepo.Change(ctx, user)
	if err != nil {
		return "", err
	}

	return publicKey, nil
}

func (us *UserServiceImpl) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	// подгружаем данные
	user, err := us.userRepo.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err = us.validateUserExists(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) UpdatePassword(ctx context.Context, userID, newPassword, oldPassword string) error {
	// подгружаем данные
	user, err := us.userRepo.Get(ctx, userID)
	if err != nil {
		return err
	}
	if err = us.validateUserExists(user); err != nil {
		return err
	}
	// generate hash
	newPasswordHash, err := us.keyCipher.EncryptString(newPassword)
	if err != nil {
		return apperrs.NewBllCommonError("generate new password hash", err)
	}
	oldPasswordHash, err := us.keyCipher.EncryptString(oldPassword)
	if err != nil {
		return apperrs.NewBllCommonError("generate old password hash", err)
	}

	// валидация
	if err = us.validateUpdatePassword(newPassword, newPasswordHash, oldPasswordHash, user); err != nil {
		return err
	}

	// store
	user.PasswordHash = newPasswordHash

	_, err = us.userRepo.Change(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// сюда можно добавить бизнес логику на проверку пароля
//   - strong password
//   - same password
//   - history password
//   - min length
//   - etc
//
// реализовываем простые проверки
//   - empty
//   - same password
//   - old and current password match
func (us *UserServiceImpl) validateUpdatePassword(newPassword, newPasswordHash, oldPasswordHash string, user *model.User) error {
	// * empty
	if strings.TrimSpace(newPassword) == "" {
		return apperrs.NewBllCommonError("new password is empty", nil)
	}
	// * same password
	if newPasswordHash == user.PasswordHash {
		return apperrs.NewBllCommonError("new password same as current old password", nil)
	}
	// * old and current password match
	if oldPasswordHash != user.PasswordHash {
		return apperrs.NewBllCommonError("old password does not match current password", nil)
	}

	return nil
}

func (us *UserServiceImpl) validateUserExists(user *model.User) error {
	if user == nil {
		return apperrs.NewBllModelNotExistsError("User", "user not found", nil)
	}

	return nil
}
