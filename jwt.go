package microcore

import (
	"errors"
	"fmt"
	jwt_go "github.com/dgrijalva/jwt-go"
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
	jwtConfig  jwtConfig
	jwtStorage JwtStorage
	contextKey ContextKey
}

func (s *Security) ContextKey() ContextKey {
	return s.contextKey
}

func (s *Security) JwtConfig() jwtConfig {
	return s.jwtConfig
}

func (s *Security) JwtStorage() JwtStorage {
	return s.jwtStorage
}

func NewSecurity(jwtConfig jwtConfig, jwtStorage JwtStorage) *Security {
	err := os.Setenv(accessSecretKey, jwtConfig.AccessSecret)
	if err != nil {
		return nil
	}
	err = os.Setenv(refreshSecretKey, jwtConfig.RefreshSecret)
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
	RefreshToken *jwt_go.Token
	AccessToken  *jwt_go.Token
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
	atClaims := jwt_go.MapClaims{}
	atClaims["access_uuid"] = td.AccessUuid
	atClaims[accountIdKey] = accountId
	atClaims["exp"] = td.AtExpires
	at := jwt_go.NewWithClaims(jwt_go.SigningMethodHS256, atClaims)
	signedString, err := at.SignedString([]byte(os.Getenv(accessSecretKey)))
	if err == nil {
		td.AccessToken = signedString
	}
	return err
}

func (s *Security) createRefreshToken(accountId int64, td *TokenDetails) error {
	rtClaims := jwt_go.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims[accountIdKey] = accountId
	rtClaims["exp"] = td.RtExpires
	rt := jwt_go.NewWithClaims(jwt_go.SigningMethodHS256, rtClaims)
	signedString, err := rt.SignedString([]byte(os.Getenv(refreshSecretKey)))
	if err == nil {
		td.RefreshToken = signedString
	}
	return err
}

func (s *Security) getToken(tokenString string, secretKey string) (*jwt_go.Token, error) {
	parser := jwt_go.Parser{SkipClaimsValidation: true}
	token, err := parser.Parse(tokenString, func(token *jwt_go.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt_go.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		//return []byte(os.Getenv(accessSecretKey)), nil
		return []byte(os.Getenv(secretKey)), nil
	})
	return token, err
}

func (s *Security) TokenValid(token *jwt_go.Token) bool {
	if _, ok := token.Claims.(jwt_go.MapClaims); ok && token.Valid && token.Claims.(jwt_go.MapClaims).VerifyExpiresAt(time.Now().Unix(), true) {
		return true
	}
	return false
}

func (s *Security) extractTokensFromHeader(r *fasthttp.Request) (accessToken *jwt_go.Token, refreshToken *jwt_go.Token, err error) {
	bearToken := string(r.Header.Peek("Authorization"))
	accessStrArr := strings.Split(bearToken, " ")
	if len(accessStrArr) != 2 {
		return nil, nil, errors.New("???????????? ????????????????")
	}
	accessToken, err = s.getToken(accessStrArr[1], accessSecretKey)
	if err != nil {
		return nil, nil, errors.New("???????????? ????????????????")
	}
	if err != nil {
		return nil, nil, errors.New("???????????? ????????????????")
	}
	rToken := string(r.Header.Peek("refresh-token"))
	if rToken != "" {
		refreshToken, err = s.getToken(rToken, refreshSecretKey)
		if err != nil {
			return nil, nil, errors.New("???????????? ????????????????")
		}
		return accessToken, refreshToken, nil
	}

	return accessToken, nil, nil
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
