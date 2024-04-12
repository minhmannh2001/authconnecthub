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

func (uc *UserUseCase) Update(u *entity.User) error {
	return uc.userRepo.Update(u)
}

func (uc *UserUseCase) FindByUsernameOrEmail(username, email string) (*entity.User, error) {
	return uc.userRepo.FindByUsernameOrEmail(username, email)
}

func (uc *UserUseCase) GetUserSocialAccounts(username string) (map[string]entity.SocialAccount, error) {
	return uc.userRepo.GetUserSocialAccounts(username)
}

func (uc *UserUseCase) AddUserSocialAccounts(username string, socialAccounts map[string]string) (bool, error) {
	return uc.userRepo.AddUserSocialAccounts(username, socialAccounts)
}

func (uc *UserUseCase) RemoveUserSocialAccount(username string, accountType string) (bool, error) {
	return uc.userRepo.RemoveUserSocialAccount(username, accountType)
}
