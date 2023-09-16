package jwtutl

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

type TokenType int

const (
	ACCESS_TOKEN TokenType = iota
	REFRESH_TOKEN
)

var (
	defaultAccessExpTime  = time.Now().Add(time.Minute * 1).Unix()
	DefaultRefreshExpTime = time.Now().Add(time.Hour * 2).Unix()
)

type TokenPair struct {
	AccessToken  string `reqHeader:"Authorization"`
	RefreshToken string `reqHeader:"Refresh-Token"`
}

type JwtTokenPair struct {
	AccessToken  *jwt.Token
	RefreshToken *jwt.Token
}

type TokenPairProperty struct {
	Email     string
	AccountId uint64
}

// JwtManager ...
type JwtManager interface {
	CreateTokenPair(property TokenPairProperty) (TokenPair, error)
	RefreshTokenPair(pair JwtTokenPair) (TokenPair, error)
	ParseToken(tokenType TokenType, token string) (*jwt.Token, error)
}

// CommonJwtManager ...
type CommonJwtManager struct {
	jwtAccessSecret  []byte
	jwtRefreshSecret []byte
}

// NewCommonJwtManager ...
func NewCommonJwtManager() JwtManager {
	return CommonJwtManager{
		jwtAccessSecret:  []byte(os.Getenv("JWT_ACCESS_SECRET")),
		jwtRefreshSecret: []byte(os.Getenv("JWT_REFRESH_SECRET")),
	}
}

// CreateTokenPair ...
func (c CommonJwtManager) CreateTokenPair(property TokenPairProperty) (TokenPair, error) {
	defaultError := fmt.Errorf("Во время авторизации произошла ошибка")
	accessClaims := jwt.MapClaims{
		"email":      property.Email,
		"account_id": property.AccountId,
		"exp":        defaultAccessExpTime,
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(c.jwtAccessSecret)
	if err != nil {
		return TokenPair{}, defaultError
	}
	refreshClaims := jwt.MapClaims{"exp": DefaultRefreshExpTime}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(c.jwtRefreshSecret)
	if err != nil {
		return TokenPair{}, defaultError
	}
	return TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// RefreshTokenPair ...
func (c CommonJwtManager) RefreshTokenPair(pair JwtTokenPair) (TokenPair, error) {
	oldAccessTokenClaims := pair.AccessToken
	oldRefreshToken := pair.RefreshToken

	_, okR := oldRefreshToken.Claims.(jwt.MapClaims)
	oldAccessClaims, okA := oldAccessTokenClaims.Claims.(jwt.MapClaims)

	if !(okR && okA && oldRefreshToken.Valid) {
		return TokenPair{}, errors.New("")
	}

	return c.CreateTokenPair(TokenPairProperty{
		Email:     oldAccessClaims["email"].(string),
		AccountId: oldAccessClaims["account_id"].(uint64),
	})
}

// ParseToken ...
func (c CommonJwtManager) ParseToken(tokenType TokenType, token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		switch tokenType {
		case ACCESS_TOKEN:
			return c.jwtAccessSecret, nil
		case REFRESH_TOKEN:
			return c.jwtRefreshSecret, nil
		}
		return nil, jwt.ValidationError{}
	})
}
