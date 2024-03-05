package repo

import (
	"errors"

	"github.com/minhmannh2001/authconnecthub/internal/entity"
	"github.com/minhmannh2001/authconnecthub/pkg/postgres"
	"gorm.io/gorm"
)

type RoleRepo struct {
	*postgres.Postgres
}

func NewRoleRepo(pg *postgres.Postgres) *RoleRepo {
	return &RoleRepo{pg}
}

func (r *RoleRepo) GetRoleIDByName(name string) (uint, error) {
	var role entity.Role
	err := r.Conn.Where("name = ?", name).First(&role).Error
	if err == nil {
		return role.ID, nil
	}

	// Handle case where no role is found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, &entity.RoleNotFoundError{Name: name}
	}

	// Handle other errors
	return 0, err
}
