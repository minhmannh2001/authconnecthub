package usecase

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/minhmannh2001/authconnecthub/config"
	"github.com/minhmannh2001/authconnecthub/internal/dto"
	"github.com/minhmannh2001/authconnecthub/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	authRepo    AuthRepo
	userUseCase UserUseCase
	privateKey  *rsa.PrivateKey
}

func NewAuthUseCase(ar AuthRepo, uu UserUseCase, pk *rsa.PrivateKey) *AuthUseCase {
	return &AuthUseCase{
		authRepo:    ar,
		userUseCase: uu,
		privateKey:  pk,
	}
}

func (au *AuthUseCase) Login(c *gin.Context, requestBody dto.LoginRequestBody) (*dto.JwtTokens, error) {
	user, err := au.userUseCase.userRepo.FindByUsernameOrEmail(requestBody.Username, "")
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password)); err != nil {
		return nil, &entity.InvalidCredentialsError{}
	}

	_cfg, ok := c.Get("config")
	if !ok {
		return nil, errors.New("config not found")
	}
	cfg := _cfg.(*config.Config)

	// Generate and return JWT tokens upon successful login
	accessToken, err := au.CreateAccessToken(*user, cfg.Authen.AccessTokenTtl)
	if err != nil {
		return nil, err
	}

	refreshToken, err := au.CreateRefreshToken(*user, accessToken, cfg.Authen.RefreshTokenTtl)
	if err != nil {
		return nil, err
	}

	return &dto.JwtTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (au *AuthUseCase) Register() {
}

func (au *AuthUseCase) CreateAccessToken(user entity.User, expireTime int) (string, error) {
	if expireTime <= 0 {
		return "", errors.New("invalid expiration time")
	}

	claims := jwt.MapClaims{
		"iss": "AuthConnect Hub",
		"sub": user.Username,
		"aud": "users",
		"exp": time.Now().Add(time.Second * time.Duration(expireTime)).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"jti": uuid.NewString(),
	}

	if au.privateKey == nil {
		return "", errors.New("missing access token private key")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	accessToken, err := token.SignedString(au.privateKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %v", err)
	}

	return accessToken, nil
}

func (au *AuthUseCase) CreateRefreshToken(user entity.User, accessToken string, expireTime int) (string, error) {
	if expireTime <= 0 {
		return "", errors.New("invalid expiration time")
	}

	accessTokenJti, err := au.RetrieveJtiFromAccessToken(accessToken, true)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"iss":              "AuthConnect Hub",
		"sub":              user.Username,
		"aud":              "users",
		"exp":              time.Now().Add(time.Second * time.Duration(expireTime)).Unix(),
		"nbf":              time.Now().Unix(),
		"iat":              time.Now().Unix(),
		"jti":              uuid.NewString(),
		"access_token_jti": accessTokenJti,
	}

	if au.privateKey == nil {
		return "", errors.New("missing refresh token private key")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	refreshToken, err := token.SignedString(au.privateKey)
	if err != nil {
		return "", fmt.Errorf("error signing refresh token. Error: %v", err)
	}

	return refreshToken, nil
}

func (au *AuthUseCase) ValidateToken(jwtToken string) (string, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return au.privateKey.Public(), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return "", errors.New("invalid token")
	}

	username, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("missing 'sub' claim in token")
	}

	return username, nil
}

func (au *AuthUseCase) RetrieveJtiFromAccessToken(jwtToken string, validate bool) (string, error) {
	var token *jwt.Token
	var err error

	if validate {
		// Parse with validation
		token, err = jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return au.privateKey.Public(), nil
		})
	} else {
		// Parse without validation
		token, _ = jwt.Parse(jwtToken, nil)
	}
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token format (not a MapClaims)")
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return "", errors.New("missing 'jti' claim in token")
	}

	// Check validity only if validation was requested
	if validate && !token.Valid {
		return "", errors.New("invalid token")
	}

	return jti, nil
}

func (au *AuthUseCase) RetrieveAccessTokenJtiFromRefreshToken(jwtToken string) (string, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return au.privateKey.Public(), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		jti, ok := claims["access_token_jti"].(string)
		if !ok {
			return "", errors.New("missing 'access_token_jti' claim in token")
		}
		return jti, nil
	} else {
		return "", errors.New("invalid token")
	}
}

func (au *AuthUseCase) IsRefreshTokenValidForAccessToken(accessToken string, refreshToken string) (bool, error) {
	accessTokenJti, err := au.RetrieveJtiFromAccessToken(accessToken, false)
	if err != nil {
		return false, err // Error retrieving JTI from access token
	}

	refreshTokenAccessTokenJti, err := au.RetrieveAccessTokenJtiFromRefreshToken(refreshToken)
	if err != nil {
		return false, err // Error retrieving access token JTI from refresh token
	}

	// Compare the JTI values
	return accessTokenJti == refreshTokenAccessTokenJti, nil
}

func (au *AuthUseCase) CheckAndRefreshTokens(oldAccessToken string, oldRefreshToken string, cfg *config.Config) (string, string, error) {
	username, err := au.ValidateToken(oldRefreshToken)
	if err != nil {
		return "", "", err // Invalid refresh token
	}

	valid, err := au.IsRefreshTokenValidForAccessToken(oldAccessToken, oldRefreshToken)
	if !valid {
		return "", "", err
	}

	user := entity.User{Username: username}

	// Create new access token
	newAccessToken, err := au.CreateAccessToken(user, cfg.Authen.AccessTokenTtl)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := au.CreateRefreshToken(user, newAccessToken, cfg.Authen.RefreshTokenTtl)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
