package service

import (
	"context"

	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/model"
	"github.com/ElfAstAhe/goph-keeper/internal/bll/server/repository"
	apperrs "github.com/ElfAstAhe/goph-keeper/internal/err"
	"github.com/ElfAstAhe/goph-keeper/pkg/logger"
	"github.com/ElfAstAhe/goph-keeper/pkg/utils"
)

type UserServiceImpl struct {
	userRepo         repository.UserRepository
	log              logger.Logger
	keyCipher        utils.Cipher
	dataCipherHelper *utils.CipherHelper
	keysHelper       *utils.RSAKeysHelper
}

func NewUserService(userRepo repository.UserRepository, dataCipherHelper *utils.CipherHelper, keysHelper *utils.RSAKeysHelper, log logger.Logger) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo:         userRepo,
		log:              log,
		dataCipherHelper: dataCipherHelper,
		keysHelper:       keysHelper,
	}
}

func (us *UserServiceImpl) Register(ctx context.Context, username, password, person, eMail string) (*model.User, error) {
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

func (us *UserServiceImpl) ChangeKeys(ctx context.Context, userID string) (string, error) {
	// подгружаем данные
	user, err := us.userRepo.Get(ctx, userID)
	if err != nil {
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
	//TODO implement me
	panic("implement me")
}

func (us *UserServiceImpl) UpdatePassword(ctx context.Context, userID, oldPass, newPass string) error {
	//TODO implement me
	panic("implement me")
}
