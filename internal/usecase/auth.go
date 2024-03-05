package usecase

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/minhmannh2001/authconnecthub/internal/entity"
)

type AuthUseCase struct {
	authRepo   AuthRepo
	privateKey *rsa.PrivateKey
}

func NewAuthUseCase(ar AuthRepo, pk *rsa.PrivateKey) *AuthUseCase {
	return &AuthUseCase{
		authRepo:   ar,
		privateKey: pk,
	}
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
		"exp": time.Now().Add(time.Hour * time.Duration(expireTime)).Unix(),
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

func (au *AuthUseCase) CreateRefreshToken(user entity.User, access_token_jti string, expireTime int) (string, error) {
	if expireTime <= 0 {
		return "", errors.New("invalid expiration time")
	}

	claims := jwt.MapClaims{
		"iss": "AuthConnect Hub",
		"sub": user.Username,
		"aud": "users",
		"exp": time.Now().Add(time.Hour * time.Duration(expireTime)).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"jti": uuid.NewString(),
	}

	if au.privateKey == nil {
		return "", errors.New("missing refresh token private key")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
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
