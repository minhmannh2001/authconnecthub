package usecases_test

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/minhmannh2001/authconnecthub/config"
	"github.com/minhmannh2001/authconnecthub/internal/dto"
	"github.com/minhmannh2001/authconnecthub/internal/entity"
	"github.com/minhmannh2001/authconnecthub/internal/helper"
	"github.com/minhmannh2001/authconnecthub/internal/usecases"
	"github.com/minhmannh2001/authconnecthub/internal/usecases/mocks"
	repoMocks "github.com/minhmannh2001/authconnecthub/internal/usecases/repos/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthUseCase_Login_Success(t *testing.T) {
	// Login request
	requestBody := dto.LoginRequestBody{Username: "testuser", Password: "secret", RememberMe: "on"}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(requestBody.Password),
		bcrypt.DefaultCost,
	)
	assert.NoError(t, err)

	// Mock dependencies
	mockUserRepo := new(mocks.IUserUC)
	mockUser := &entity.User{ID: 1, Username: "testuser", Password: string(encryptedPassword)}
	mockUserRepo.On("FindByUsernameOrEmail", "testuser", "").Return(mockUser, nil)

	reader := rand.Reader
	bitSize := 2048

	privateKey, err := rsa.GenerateKey(reader, bitSize)
	assert.NoError(t, err)

	mockConfig := &config.Config{Authen: config.Authen{
		RefreshTokenTtl: 3600,
		AccessTokenTtl:  600,
		JwtPrivateKey:   privateKey,
	}}

	mockContext := gin.Context{}
	mockContext.Set("config", mockConfig)

	// Create use case with mocks
	uc := usecases.NewAuthUseCase(nil, mockUserRepo, mockConfig)

	// Call Login
	tokens, err := uc.Login(&mockContext, requestBody)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, tokens)
	assert.NotEmpty(t, tokens.AccessToken)
	assert.NotEmpty(t, tokens.RefreshToken)
}

func TestAuthUseCase_ValidateToken_Success(t *testing.T) {
	// Generate a valid JWT token
	reader := rand.Reader
	bitSize := 2048
	privateKey, err := rsa.GenerateKey(reader, bitSize)
	assert.NoError(t, err)

	mockConfig := &config.Config{Authen: config.Authen{
		RefreshTokenTtl: 3600,
		AccessTokenTtl:  600,
		JwtPrivateKey:   privateKey,
	}}

	mockUser := &entity.User{Username: "testuser"}
	claims := jwt.MapClaims{
		"iss": "AuthConnect Hub",
		"sub": mockUser.Username,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	assert.NoError(t, err)

	// Create use case with private key
	uc := usecases.NewAuthUseCase(nil, nil, mockConfig)

	// Call ValidateToken
	username, err := uc.ValidateToken(tokenString)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, mockUser.Username, username)
}

func TestAuthUseCase_ValidateToken_InvalidToken(t *testing.T) {
	// Generate a valid JWT token
	reader := rand.Reader
	bitSize := 2048
	privateKey, err := rsa.GenerateKey(reader, bitSize)
	assert.NoError(t, err)

	mockConfig := &config.Config{Authen: config.Authen{
		RefreshTokenTtl: 3600,
		AccessTokenTtl:  600,
		JwtPrivateKey:   privateKey,
	}}

	invalidTokens := []string{
		"invalid_token_string",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0dXNlciIsInNjb3BlIjplbmQifQ==", // Tampered token
	}

	// Create use case with private key (doesn't matter for these tests)
	uc := usecases.NewAuthUseCase(nil, nil, mockConfig)

	for _, token := range invalidTokens {
		username, err := uc.ValidateToken(token)
		assert.Error(t, err)
		assert.Empty(t, username)
	}
}

func TestAuthUseCase_ValidateToken_ExpiredToken(t *testing.T) {
	reader := rand.Reader
	bitSize := 2048
	privateKey, err := rsa.GenerateKey(reader, bitSize)
	assert.NoError(t, err)

	mockConfig := &config.Config{Authen: config.Authen{
		RefreshTokenTtl: 3600,
		AccessTokenTtl:  600,
		JwtPrivateKey:   privateKey,
	}}

	mockUser := &entity.User{Username: "testuser"}
	claims := jwt.MapClaims{
		"iss": "AuthConnect Hub",
		"sub": mockUser.Username,
		"exp": time.Now().Add(time.Second * -10).Unix(), // Expired token
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	assert.NoError(t, err)

	// Create use case with private key
	uc := usecases.NewAuthUseCase(nil, nil, mockConfig)

	// Call ValidateToken
	username, err := uc.ValidateToken(tokenString)

	// Assertions
	assert.Error(t, err)
	assert.EqualError(t, err, "token has invalid claims: token is expired")
	assert.Empty(t, username)
}

func TestAuthUseCase_ValidateToken_MissingSubClaim(t *testing.T) {
	reader := rand.Reader
	bitSize := 2048
	privateKey, err := rsa.GenerateKey(reader, bitSize)
	assert.NoError(t, err)

	mockConfig := &config.Config{Authen: config.Authen{
		RefreshTokenTtl: 3600,
		AccessTokenTtl:  600,
		JwtPrivateKey:   privateKey,
	}}

	claims := jwt.MapClaims{
		"iss": "AuthConnect Hub",
		"exp": time.Now().Add(time.Minute * 10).Unix(), // Expired token
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	assert.NoError(t, err)

	// Create use case with private key
	uc := usecases.NewAuthUseCase(nil, nil, mockConfig)

	// Call ValidateToken
	username, err := uc.ValidateToken(tokenString)

	// Assertions
	assert.Error(t, err)
	assert.EqualError(t, err, "missing 'sub' claim in token")
	assert.Empty(t, username)
}

func TestAuthUseCase_IsRefreshTokenValidForAccessToken_Valid(t *testing.T) {
	reader := rand.Reader
	bitSize := 2048

	privateKey, err := rsa.GenerateKey(reader, bitSize)
	assert.NoError(t, err)

	mockConfig := &config.Config{Authen: config.Authen{
		RefreshTokenTtl: 3600,
		AccessTokenTtl:  600,
		JwtPrivateKey:   privateKey,
	}}

	uc := usecases.NewAuthUseCase(nil, nil, mockConfig)

	mockUser := entity.User{ID: 1, Username: "testuser"}

	accessToken, err := uc.CreateAccessToken(mockUser, mockConfig.Authen.AccessTokenTtl)
	assert.NoError(t, err)

	refreshToken, err := uc.CreateRefreshToken(mockUser, accessToken, mockConfig.Authen.RefreshTokenTtl)
	assert.NoError(t, err)

	isValid, err := uc.IsRefreshTokenValidForAccessToken(accessToken, refreshToken)

	assert.NoError(t, err)
	assert.True(t, isValid)
}

func TestAuthUseCase_IsRefreshTokenValidForAccessToken_InvalidJti(t *testing.T) {
	reader := rand.Reader
	bitSize := 2048

	privateKey, err := rsa.GenerateKey(reader, bitSize)
	assert.NoError(t, err)

	mockConfig := &config.Config{Authen: config.Authen{
		RefreshTokenTtl: 3600,
		AccessTokenTtl:  600,
		JwtPrivateKey:   privateKey,
	}}

	uc := usecases.NewAuthUseCase(nil, nil, mockConfig)

	mockUser := entity.User{ID: 1, Username: "testuser"}

	accessToken, err := uc.CreateAccessToken(mockUser, mockConfig.Authen.AccessTokenTtl)
	assert.NoError(t, err)

	refreshToken, err := uc.CreateRefreshToken(mockUser, accessToken, mockConfig.Authen.RefreshTokenTtl)
	assert.NoError(t, err)

	anotherAccessToken, err := uc.CreateAccessToken(mockUser, mockConfig.Authen.AccessTokenTtl)
	assert.NoError(t, err)

	isValid, err := uc.IsRefreshTokenValidForAccessToken(anotherAccessToken, refreshToken)

	assert.NoError(t, err)
	assert.False(t, isValid)
}

func TestAuthUseCase_IsRefreshTokenValidForAccessToken_ErrorRetrievingRefreshTokenAccessTokenJti(t *testing.T) {
	accessToken := "valid_token"
	refreshToken := "invalid_token"

	mockConfig := &config.Config{Authen: config.Authen{
		RefreshTokenTtl: 3600,
		AccessTokenTtl:  600,
	}}

	uc := usecases.NewAuthUseCase(nil, nil, mockConfig)

	isValid, err := uc.IsRefreshTokenValidForAccessToken(accessToken, refreshToken)

	assert.Error(t, err)
	assert.False(t, isValid)
}

func TestAuthUseCase_CheckAndRefreshTokens_Success(t *testing.T) {
	reader := rand.Reader
	bitSize := 2048

	privateKey, err := rsa.GenerateKey(reader, bitSize)
	assert.NoError(t, err)

	mockConfig := &config.Config{Authen: config.Authen{
		RefreshTokenTtl: 3600,
		AccessTokenTtl:  1,
		JwtPrivateKey:   privateKey,
	}}

	uc := usecases.NewAuthUseCase(nil, nil, mockConfig)

	mockUser := entity.User{ID: 1, Username: "testuser"}

	accessToken, err := uc.CreateAccessToken(mockUser, mockConfig.Authen.AccessTokenTtl)
	assert.NoError(t, err)

	refreshToken, err := uc.CreateRefreshToken(mockUser, accessToken, mockConfig.Authen.RefreshTokenTtl)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	newAccessToken, newRefreshToken, err := uc.CheckAndRefreshTokens(accessToken, refreshToken, mockConfig)

	assert.NoError(t, err)
	assert.NotEmpty(t, newAccessToken)
	assert.NotEmpty(t, newRefreshToken)
}

func TestAuthUseCase_CheckAndRefreshTokens_InvalidRefreshToken(t *testing.T) {
	reader := rand.Reader
	bitSize := 2048

	privateKey, err := rsa.GenerateKey(reader, bitSize)
	assert.NoError(t, err)

	mockConfig := &config.Config{Authen: config.Authen{
		RefreshTokenTtl: 1,
		AccessTokenTtl:  1,
		JwtPrivateKey:   privateKey,
	}}

	uc := usecases.NewAuthUseCase(nil, nil, mockConfig)

	mockUser := entity.User{ID: 1, Username: "testuser"}

	accessToken, err := uc.CreateAccessToken(mockUser, mockConfig.Authen.AccessTokenTtl)
	assert.NoError(t, err)

	refreshToken, err := uc.CreateRefreshToken(mockUser, accessToken, mockConfig.Authen.RefreshTokenTtl)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	newAccessToken, newRefreshToken, err := uc.CheckAndRefreshTokens(accessToken, refreshToken, mockConfig)

	assert.Error(t, err)
	assert.Empty(t, newAccessToken)
	assert.Empty(t, newRefreshToken)
}

func TestAuthUseCase_CheckAndRefreshTokens_InvalidAccessTokenJti(t *testing.T) {
	reader := rand.Reader
	bitSize := 2048

	privateKey, err := rsa.GenerateKey(reader, bitSize)
	assert.NoError(t, err)

	mockConfig := &config.Config{Authen: config.Authen{
		RefreshTokenTtl: 100,
		AccessTokenTtl:  100,
		JwtPrivateKey:   privateKey,
	}}

	uc := usecases.NewAuthUseCase(nil, nil, mockConfig)

	mockUser := entity.User{ID: 1, Username: "testuser"}

	accessToken, err := uc.CreateAccessToken(mockUser, mockConfig.Authen.AccessTokenTtl)
	assert.NoError(t, err)

	refreshToken, err := uc.CreateRefreshToken(mockUser, accessToken, mockConfig.Authen.RefreshTokenTtl)
	assert.NoError(t, err)

	anotherAccessToken, err := uc.CreateAccessToken(mockUser, mockConfig.Authen.AccessTokenTtl)
	assert.NoError(t, err)

	newAccessToken, newRefreshToken, err := uc.CheckAndRefreshTokens(anotherAccessToken, refreshToken, mockConfig)

	assert.NoError(t, err)
	assert.Empty(t, newAccessToken)
	assert.Empty(t, newRefreshToken)
}

func TestAuthUseCase_RetrieveFieldFromJwtToken_ValidToken_RequiredValidation(t *testing.T) {
	reader := rand.Reader
	bitSize := 2048
	privateKey, err := rsa.GenerateKey(reader, bitSize)
	assert.NoError(t, err)

	mockUser := &entity.User{Username: "testuser"}
	claims := jwt.MapClaims{
		"iss":      "AuthConnect Hub",
		"sub":      mockUser.Username,
		"username": mockUser.Username, // Example field
		"exp":      time.Now().Add(time.Minute * 10).Unix(),
		"nbf":      time.Now().Unix(),
		"iat":      time.Now().Unix(),
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	assert.NoError(t, err)

	mockConfig := &config.Config{Authen: config.Authen{
		JwtPrivateKey: privateKey,
	}}
	uc := usecases.NewAuthUseCase(nil, nil, mockConfig)

	fieldValue, err := uc.RetrieveFieldFromJwtToken(tokenString, "username", true) // Required validation

	assert.NoError(t, err)
	assert.Equal(t, mockUser.Username, fieldValue)
}

func TestAuthUseCase_RetrieveFieldFromJwtToken_InvalidToken(t *testing.T) {
	invalidTokens := []string{
		"invalid_token_string",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0dXNlciIsInNjb3BlIjplbmQifQ==", // Tampered token
	}

	reader := rand.Reader
	bitSize := 2048
	privateKey, err := rsa.GenerateKey(reader, bitSize)
	assert.NoError(t, err)

	mockConfig := &config.Config{Authen: config.Authen{
		JwtPrivateKey: privateKey,
	}}
	uc := usecases.NewAuthUseCase(nil, nil, mockConfig)

	for _, token := range invalidTokens {
		fieldValue, err := uc.RetrieveFieldFromJwtToken(token, "username", true) // Required validation

		assert.Error(t, err)
		assert.Empty(t, fieldValue)
	}
}

func TestAuthUseCase_RetrieveFieldFromJwtToken_MissingField(t *testing.T) {
	reader := rand.Reader
	bitSize := 2048
	privateKey, err := rsa.GenerateKey(reader, bitSize)
	assert.NoError(t, err)

	mockUser := &entity.User{Username: "testuser"}
	claims := jwt.MapClaims{
		"iss": "AuthConnect Hub",
		"sub": mockUser.Username,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(privateKey)
	assert.NoError(t, err)

	mockConfig := &config.Config{Authen: config.Authen{
		JwtPrivateKey: privateKey,
	}}
	uc := usecases.NewAuthUseCase(nil, nil, mockConfig)

	fieldValue, err := uc.RetrieveFieldFromJwtToken(tokenString, "missing_field", true) // Required validation

	assert.Error(t, err)
	assert.EqualError(t, err, "missing 'missing_field' claim in token")
	assert.Empty(t, fieldValue)
}

func TestAuthUseCase_Logout_Success(t *testing.T) {
	reader := rand.Reader
	bitSize := 2048
	privateKey, err := rsa.GenerateKey(reader, bitSize)
	assert.NoError(t, err)

	// Mock dependencies
	mockContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	mockConfig := &config.Config{Authen: config.Authen{
		AccessTokenTtl:  3600,
		RefreshTokenTtl: 3600,
		JwtPrivateKey:   privateKey,
	}}
	mockContext.Set("config", mockConfig)

	mockAuthRepo := repoMocks.NewIAuthRepo(t)
	mockAuthRepo.On("BlacklistToken", mock.Anything, mock.Anything).Return(nil) // Successful blacklist

	uc := usecases.NewAuthUseCase(mockAuthRepo, nil, mockConfig)

	mockUser := entity.User{ID: 1, Username: "testuser"}

	accessToken, err := uc.CreateAccessToken(mockUser, mockConfig.Authen.AccessTokenTtl)
	assert.NoError(t, err)

	refreshToken, err := uc.CreateRefreshToken(mockUser, accessToken, mockConfig.Authen.RefreshTokenTtl)
	assert.NoError(t, err)

	// Set access and refresh tokens in headers (simulated)
	mockContext.Request = &http.Request{}
	mockContext.Request.Header = http.Header{}
	mockContext.Request.Header.Set(helper.AccessTokenHeader, "Bearer "+accessToken)
	mockContext.Request.Header.Set(helper.RefreshTokenHeader, "Bearer "+refreshToken)

	err = uc.Logout(mockContext)

	// Assertions
	assert.NoError(t, err)
	mockAuthRepo.AssertExpectations(t) // Ensure blacklist calls were made
}
