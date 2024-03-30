package usecases

import (
	"github.com/minhmannh2001/authconnecthub/internal/entity"
	"github.com/minhmannh2001/authconnecthub/internal/usecases/repos"
)

type UserUseCase struct {
	userRepo repos.IUserRepo
}

func NewUserUseCase(ur repos.IUserRepo) *UserUseCase {
	return &UserUseCase{userRepo: ur}
}

func (uc *UserUseCase) Create(u entity.User) (entity.User, error) {
	return uc.userRepo.Create(u)

}

func (uc *UserUseCase) FindByUsernameOrEmail(username, email string) (*entity.User, error) {
	return uc.userRepo.FindByUsernameOrEmail(username, email)
}
