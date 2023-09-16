package jwtutl

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

type TokenType int

const (
	ACCESS_TOKEN_TYPE TokenType = iota
	REFRESH_TOKEN_TYPE
)

const UserLocalNameKey = "user"

type TokenPair struct {
	AccessToken  string `reqHeader:"Authorization"`
	RefreshToken string `reqHeader:"Refresh-Token"`
}

type JwtTokenPair struct {
	AccessToken     *jwt.Token
	AccessTokenStr  string
	RefreshToken    *jwt.Token
	RefreshTokenStr string
}

type TokenPairProperty struct {
	Email     string
	AccountId uint64
}

// JwtManager ...
type JwtManager interface {
	CreateTokenPair(property TokenPairProperty) (JwtTokenPair, error)
	RefreshTokenPair(pair JwtTokenPair) (JwtTokenPair, error)
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
func (c CommonJwtManager) CreateTokenPair(property TokenPairProperty) (JwtTokenPair, error) {
	defaultAccessExpTime := time.Now().Add(time.Minute * 1).Unix()
	DefaultRefreshExpTime := time.Now().Add(time.Hour * 2).Unix()

	accessClaims := jwt.MapClaims{"email": property.Email, "account_id": property.AccountId, "exp": defaultAccessExpTime}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	refreshClaims := jwt.MapClaims{"exp": DefaultRefreshExpTime}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessTokenStr, err := accessToken.SignedString(c.jwtAccessSecret)
	if err != nil {
		return JwtTokenPair{}, err
	}

	refreshTokenStr, err := refreshToken.SignedString(c.jwtRefreshSecret)
	if err != nil {
		return JwtTokenPair{}, err
	}

	return JwtTokenPair{
		AccessToken:     accessToken,
		AccessTokenStr:  accessTokenStr,
		RefreshToken:    refreshToken,
		RefreshTokenStr: refreshTokenStr,
	}, nil
}

// RefreshTokenPair ...
func (c CommonJwtManager) RefreshTokenPair(pair JwtTokenPair) (JwtTokenPair, error) {
	oldAccessTokenClaims := pair.AccessToken
	oldRefreshToken := pair.RefreshToken

	_, okR := oldRefreshToken.Claims.(jwt.MapClaims)
	oldAccessClaims, okA := oldAccessTokenClaims.Claims.(jwt.MapClaims)

	defaultErr := errors.New("Не удалось обновить токены")

	if !(okR && okA && oldRefreshToken.Valid) {
		return JwtTokenPair{}, defaultErr
	}

	return c.CreateTokenPair(TokenPairProperty{
		Email:     oldAccessClaims["email"].(string),
		AccountId: uint64(oldAccessClaims["account_id"].(float64)),
	})
}

// ParseToken ...
func (c CommonJwtManager) ParseToken(tokenType TokenType, token string) (*jwt.Token, error) {
	// Парсим access
	parseFunc := func(token *jwt.Token) (interface{}, error) {
		return c.jwtAccessSecret, nil
	}
	// Парсим refresh
	if tokenType == REFRESH_TOKEN_TYPE {
		parseFunc = func(token *jwt.Token) (interface{}, error) {
			return c.jwtRefreshSecret, nil
		}
	}
	return jwt.Parse(token, parseFunc)
}
