package usecase

type RoleUseCase struct {
	roleRepo RoleRepo
}

func NewRoleUseCase(rr RoleRepo) *RoleUseCase {
	return &RoleUseCase{roleRepo: rr}
}

func (uc *RoleUseCase) GetRoleIDByName(name string) (uint, error) {
	return uc.roleRepo.GetRoleIDByName(name)
}
