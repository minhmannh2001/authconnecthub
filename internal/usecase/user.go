package usecase

import (
	"github.com/minhmannh2001/authconnecthub/internal/entity"
)

type UserUseCase struct {
	userRepo UserRepo
}

func NewUserUseCase(ur UserRepo) *UserUseCase {
	return &UserUseCase{userRepo: ur}
}

func (uc *UserUseCase) Create(u entity.User) (entity.User, error) {
	return uc.userRepo.Create(u)

}

func (uc *UserUseCase) FindByUsernameOrEmail(username, email string) (*entity.User, error) {
	return uc.userRepo.FindByUsernameOrEmail(username, email)
}
