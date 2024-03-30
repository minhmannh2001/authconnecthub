package repos

import (
	"errors"

	"github.com/minhmannh2001/authconnecthub/internal/entity"
	"github.com/minhmannh2001/authconnecthub/pkg/postgres"
	"gorm.io/gorm"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) Create(u entity.User) (entity.User, error) {
	// Check for existing user with same username or email
	var existingUser entity.User
	err := r.Conn.Where("username = ?", u.Username).Or("email = ?", u.Email).First(&existingUser).Error
	if err == nil { // User already exists
		return entity.User{}, &entity.ErrDuplicateUser{Username: u.Username, Email: u.Email}
	}
	// If no existing user found, create the new user
	result := r.Conn.Create(&u)
	if err := result.Error; err != nil {
		return entity.User{}, err
	}

	return u, nil
}

func (r *UserRepo) RetrieveByID(id uint) (entity.User, error) {
	return entity.User{}, nil
}

func (r *UserRepo) Update(u entity.User) (entity.User, error) {
	return entity.User{}, nil
}

func (r *UserRepo) Delete(u entity.User) (entity.User, error) {
	return entity.User{}, nil
}

func (r *UserRepo) FindByUsernameOrEmail(username, email string) (*entity.User, error) {
	var user entity.User
	err := r.Conn.Where("username = ? OR email = ?", username, email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &entity.InvalidCredentialsError{} // User not found
		}
		return nil, err
	}
	return &user, nil
}
