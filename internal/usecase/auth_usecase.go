package usecase

import (
	"golang-redis/internal/entity"
	"golang-redis/internal/repository"

	"github.com/google/uuid"
)

type AuthUsecase interface {
	Login(user entity.User) (string, error)
	Logout(token string) error
	ValidateToken(token string) (*entity.Session, error)
}

type authUsecase struct {
	sessionRepo repository.SessionRepository
}

func NewAuthUsecase(sessionRepo repository.SessionRepository) AuthUsecase {
	return &authUsecase{sessionRepo: sessionRepo}
}

func (u *authUsecase) ValidateToken(token string) (*entity.Session, error) {
	return u.sessionRepo.GetByToken(token)
}

func (u *authUsecase) Login(user entity.User) (string, error) {
	token := uuid.NewString()

	session := entity.Session{
		UserID: user.Id,
		Token:  token,
	}

	if err := u.sessionRepo.Save(session); err != nil {
		return "", err
	}

	return token, nil
}

func (u *authUsecase) Logout(token string) error {
	return u.sessionRepo.Delete(token)
}
