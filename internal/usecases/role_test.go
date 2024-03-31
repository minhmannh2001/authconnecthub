package usecases_test

import (
	"errors"
	"testing"

	"github.com/minhmannh2001/authconnecthub/internal/entity"
	"github.com/minhmannh2001/authconnecthub/internal/usecases"
	"github.com/minhmannh2001/authconnecthub/internal/usecases/repos/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetRoleIDByName_Success(t *testing.T) {
	mockRoleRepo := mocks.NewIRoleRepo(t)

	expectedRoleID := uint(1)
	expectedName := "admin"

	mockRoleRepo.On("GetRoleIDByName", expectedName).Return(expectedRoleID, nil)

	uc := usecases.NewRoleUseCase(mockRoleRepo)

	roleID, err := uc.GetRoleIDByName(expectedName)

	assert.Nil(t, err)
	assert.Equal(t, expectedRoleID, roleID)

	mockRoleRepo.AssertExpectations(t)
}

func TestGetRoleIDByName_RoleNotFound(t *testing.T) {
	mockRoleRepo := mocks.NewIRoleRepo(t)
	nonExistentName := "unknown-role"

	mockRoleRepo.On("GetRoleIDByName", nonExistentName).Return(
		uint(0),
		&entity.RoleNotFoundError{Name: nonExistentName},
	)

	uc := usecases.NewRoleUseCase(mockRoleRepo)

	roleID, err := uc.GetRoleIDByName(nonExistentName)

	assert.NotNil(t, err)
	assert.Equal(t, uint(0), roleID) // Expect 0 as default value for non-existent roles
	assert.Equal(
		t,
		err,
		&entity.RoleNotFoundError{Name: nonExistentName},
	)

	mockRoleRepo.AssertExpectations(t)
}

func TestGetRoleIDByName_DatabaseError(t *testing.T) {
	mockRoleRepo := mocks.NewIRoleRepo(t)

	expectedName := "admin"
	expectedError := errors.New("database error")

	mockRoleRepo.On("GetRoleIDByName", expectedName).Return(uint(0), expectedError)

	uc := usecases.NewRoleUseCase(mockRoleRepo)

	roleID, err := uc.GetRoleIDByName(expectedName)

	// Assertions
	assert.Equal(t, expectedError, err)
	assert.Equal(t, uint(0), roleID) // No role ID returned

	mockRoleRepo.AssertExpectations(t)
}
