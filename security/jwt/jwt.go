package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/greenfield0000/microcore/configuration"
	"github.com/greenfield0000/microcore/security/jwt/storage"
	"github.com/valyala/fasthttp"
	"os"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

const (
	accessSecretKey  = "ACCESS_SECRET"
	refreshSecretKey = "REFRESH_SECRET"
	accountIdKey     = "account_id"
)

type ContextKey struct {
	AccountIdKey string
}

type Security struct {
	jwtConfig  configuration.Jwt
	jwtStorage storage.JwtStorage
	contextKey ContextKey
}

func (s *Security) ContextKey() ContextKey {
	return s.contextKey
}

func (s *Security) JwtConfig() configuration.Jwt {
	return s.jwtConfig
}

func (s *Security) JwtStorage() storage.JwtStorage {
	return s.jwtStorage
}

func NewSecurity(jwtConfig configuration.Jwt, jwtStorage storage.JwtStorage) *Security {
	err := os.Setenv(accessSecretKey, jwtConfig.Accesssecret)
	if err != nil {
		return nil
	}
	err = os.Setenv(refreshSecretKey, jwtConfig.Refreshsecret)
	if err != nil {
		return nil
	}
	return &Security{
		jwtConfig:  jwtConfig,
		jwtStorage: jwtStorage,
		contextKey: ContextKey{AccountIdKey: accountIdKey},
	}
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    float64
	RtExpires    float64
}

type CustomTokenDetail struct {
	RefreshToken *jwt.Token
	AccessToken  *jwt.Token
}

func (s *Security) CreateTokenPair(accountId int64) (*TokenDetails, error) {
	var err error

	td := &TokenDetails{}
	td.AtExpires = float64(time.Now().Add(time.Minute * 10).Unix())
	v4, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	td.AccessUuid = v4.String()

	td.RtExpires = float64(time.Now().Add(time.Hour * 2).Unix())
	newV4, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	td.RefreshUuid = newV4.String()
	//Creating Access Token
	err = s.createAccessToken(accountId, td)
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	err = s.createRefreshToken(accountId, td)
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (s *Security) createAccessToken(accountId int64, td *TokenDetails) error {
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = td.AccessUuid
	atClaims[accountIdKey] = accountId
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	signedString, err := at.SignedString([]byte(os.Getenv(accessSecretKey)))
	if err == nil {
		td.AccessToken = signedString
	}
	return err
}

func (s *Security) createRefreshToken(accountId int64, td *TokenDetails) error {
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims[accountIdKey] = accountId
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	signedString, err := rt.SignedString([]byte(os.Getenv(refreshSecretKey)))
	if err == nil {
		td.RefreshToken = signedString
	}
	return err
}

func (s *Security) getToken(tokenString string, secretKey string) (*jwt.Token, error) {
	parser := jwt.Parser{SkipClaimsValidation: true}
	token, err := parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		//return []byte(os.Getenv(accessSecretKey)), nil
		return []byte(os.Getenv(secretKey)), nil
	})
	return token, err
}

func (s *Security) TokenValid(token *jwt.Token) bool {
	if _, ok := s.GetTokenClaims(token).(jwt.MapClaims); ok && token.Valid && token.Claims.(jwt.MapClaims).VerifyExpiresAt(time.Now().Unix(), true) {
		return true
	}
	return false
}

func (s *Security) GetTokenClaims(token *jwt.Token) jwt.Claims {
	if token == nil {
		return nil
	}
	return token.Claims
}

func (s *Security) extractTokensFromHeader(r *fasthttp.Request) (accessToken *jwt.Token, refreshToken *jwt.Token, err error) {
	bearToken := string(r.Header.Peek("Authorization"))
	accessStrArr := strings.Split(bearToken, " ")
	if len(accessStrArr) != 2 {
		return nil, nil, errors.New("Доступ запрещен")
	}
	accessToken, err = s.getToken(accessStrArr[1], accessSecretKey)
	if err != nil {
		return nil, nil, errors.New("Доступ запрещен")
	}
	if err != nil {
		return nil, nil, errors.New("Доступ запрещен")
	}
	rToken := string(r.Header.Peek("RToken"))
	if len(rToken) == 0 {
		return nil, nil, errors.New("Доступ запрещен")
	}

	refreshToken, err = s.getToken(rToken, refreshSecretKey)
	if err != nil {
		return nil, nil, errors.New("Доступ запрещен")
	}
	return accessToken, refreshToken, nil
}

func (s *Security) ExtractTokenPairDetails(r *fasthttp.Request) (*CustomTokenDetail, error) {
	accessToken, refreshToken, err := s.extractTokensFromHeader(r)
	if err != nil {
		return nil, err
	}

	return &CustomTokenDetail{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
