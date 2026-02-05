package service

import "github.com/ElfAstAhe/goph-keeper/internal/bll/server/repository"

type AuthServiceImpl struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthServiceImpl {
	return &AuthServiceImpl{
		userRepo: userRepo,
	}
}
