package repos

import (
	"errors"

	"github.com/minhmannh2001/authconnecthub/internal/entity"
	"github.com/minhmannh2001/authconnecthub/internal/helper"
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

func (r *UserRepo) RetrieveByID(id uint) (entity.User, error) { // coverage-ignore
	return entity.User{}, nil
}

func (r *UserRepo) Update(u *entity.User) error {
	result := r.Conn.Session(&gorm.Session{FullSaveAssociations: true}).Updates(u)
	err := result.Error

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) Delete(u entity.User) (entity.User, error) { // coverage-ignore
	return entity.User{}, nil
}

func (r *UserRepo) FindByUsernameOrEmail(username, email string) (*entity.User, error) {
	var user entity.User

	if email == "" {
		email = helper.RandStringBytes(24)
	}

	err := r.Conn.Preload("UserProfile").Where("username = ? OR email = ?", username, email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &entity.InvalidCredentialsError{} // User not found
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetUserSocialAccounts(username string) (map[string]entity.SocialAccount, error) {
	// Find user by username
	user, err := r.FindByUsernameOrEmail(username, "")
	if err != nil {
		return nil, err
	}

	// Create a map to store social accounts (account type as key)
	socialAccounts := make(map[string]entity.SocialAccount)

	// Define social account types
	accountTypes := []string{"facebook", "twitter", "github", "youtube"}

	// Loop through each account type
	for _, accountType := range accountTypes {
		// Initialize a SocialAccount with the account type
		socialAccount := entity.SocialAccount{
			UserID:      user.ID,
			AccountType: accountType,
		}

		// Fetch account link
		err := r.Conn.Where("user_id = ? AND account_type = ?", user.ID, accountType).First(&socialAccount).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				socialAccount.AccountLink = ""
			} else {
				return nil, err
			}
		}

		// Add social account to the map
		socialAccounts[accountType] = socialAccount
	}

	return socialAccounts, nil
}

func (r *UserRepo) AddUserSocialAccounts(username string, socialAccounts map[string]string) (bool, error) {
	// Find user by username
	user, err := r.FindByUsernameOrEmail(username, "")
	if err != nil {
		return false, err
	}

	// Iterate through each social account entry
	for accountType, accountLink := range socialAccounts {
		// Create a new SocialAccount instance
		socialAccount := entity.SocialAccount{
			UserID:      user.ID,
			AccountType: accountType,
			AccountLink: accountLink,
		}

		// Insert the social account into the database
		result := r.Conn.Create(&socialAccount)
		err := result.Error

		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (r *UserRepo) RemoveUserSocialAccount(username string, accountType string) (bool, error) {
	// Find user by username
	user, err := r.FindByUsernameOrEmail(username, "")
	if err != nil {
		return false, err
	}

	// Delete social account based on user ID and account type
	result := r.Conn.Where("user_id = ? AND account_type = ?", user.ID, accountType).Delete(&entity.SocialAccount{})
	err = result.Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil // Social account not found, not an error (treated as successful deletion)
		}
		return false, err
	}

	return true, nil
}
