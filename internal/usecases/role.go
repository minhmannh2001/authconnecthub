package usecases

import "github.com/minhmannh2001/authconnecthub/internal/usecases/repos"

type RoleUseCase struct {
	roleRepo repos.IRoleRepo
}

func NewRoleUseCase(rr repos.IRoleRepo) *RoleUseCase {
	return &RoleUseCase{roleRepo: rr}
}

func (uc *RoleUseCase) GetRoleIDByName(name string) (uint, error) {
	return uc.roleRepo.GetRoleIDByName(name)
}
