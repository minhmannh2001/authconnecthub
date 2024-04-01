package usecases_test

import (
	"testing"

	"github.com/minhmannh2001/authconnecthub/internal/entity"
	"github.com/minhmannh2001/authconnecthub/internal/usecases"
	"github.com/minhmannh2001/authconnecthub/internal/usecases/repos/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserUseCase_Create_Success(t *testing.T) {
	// Create a user
	user := entity.User{Username: "testuser", Email: "test@example.com"}

	// Create a mock repository
	mockRepo := new(mocks.IUserRepo)
	mockRepo.On("Create", mock.Anything).Return(user, nil)

	// Create the use case with the mock repository
	uc := usecases.NewUserUseCase(mockRepo)

	// Call Create and assert the result
	createdUser, err := uc.Create(user)
	assert.NoError(t, err)
	assert.Equal(t, user, createdUser)

	// Verify that the mock repository was called with the correct arguments
	mockRepo.AssertExpectations(t)
}

func TestUserUseCase_Create_DuplicateUser(t *testing.T) {
	// Create a mock repository with a duplicate user error
	mockRepo := new(mocks.IUserRepo)
	mockRepo.On("Create", mock.Anything).Return(entity.User{}, &entity.ErrDuplicateUser{})

	// Create the use case
	uc := usecases.NewUserUseCase(mockRepo)

	// Create a user
	user := entity.User{Username: "testuser", Email: "test@example.com"}

	// Call Create and assert the error
	_, err := uc.Create(user)
	assert.Error(t, err)
	assert.Equal(t, &entity.ErrDuplicateUser{}, err)
}

func TestUserUseCase_FindByUsernameOrEmail_Success(t *testing.T) {
	// Create a mock repository with a user
	mockRepo := new(mocks.IUserRepo)
	user := entity.User{ID: 1, Username: "testuser", Email: "test@example.com"}
	mockRepo.On("FindByUsernameOrEmail", "testuser", "test@example.com").Return(&user, nil)

	// Create the use case
	uc := usecases.NewUserUseCase(mockRepo)

	// Call FindByUsernameOrEmail and assert the result
	foundUser, err := uc.FindByUsernameOrEmail("testuser", "test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, &user, foundUser)
}

func TestUserUseCase_FindByUsernameOrEmail_NotFound(t *testing.T) {
	// Create a mock repository with a user not found error
	mockRepo := new(mocks.IUserRepo)
	mockRepo.On("FindByUsernameOrEmail", "testuser", "test@example.com").Return(nil, &entity.InvalidCredentialsError{})

	// Create the use case
	uc := usecases.NewUserUseCase(mockRepo)

	// Call FindByUsernameOrEmail and assert the error
	_, err := uc.FindByUsernameOrEmail("testuser", "test@example.com")
	assert.Error(t, err)
	assert.IsType(t, &entity.InvalidCredentialsError{}, err)
}
