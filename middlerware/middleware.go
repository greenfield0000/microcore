package middlerware

import (
	"encoding/json"
	jwt_go "github.com/dgrijalva/jwt-go"
	"github.com/greenfield0000/microcore/http-common"
	middleware_auth "github.com/greenfield0000/microcore/middlerware/middleware-auth"
	"github.com/greenfield0000/microcore/security"
	"github.com/greenfield0000/microcore/security/jwt"
	"github.com/greenfield0000/microcore/security/jwt/storage"
	"github.com/valyala/fasthttp"
	"microcore/middlerware/domains"
	"time"
)

const (
	contentTypeJson = "application/json; charset=utf-8"
	accessTokenKey  = "access-token"
	refreshTokenKey = "refresh-token"
	// cors
	corsAllowHeaders     = "*"
	corsAllowMethods     = "*"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "false"
	corsExposeHeaders    = "Access-Token"
)

type Service interface {
}

type IMiddleWare interface {
	wareAll(ctx *fasthttp.RequestCtx)
	WareLogin(next fasthttp.RequestHandler, security *jwt.Security) fasthttp.RequestHandler
	WareLogout(next fasthttp.RequestHandler, security *jwt.Security) fasthttp.RequestHandler
	WareRegistry(next fasthttp.RequestHandler, security *jwt.Security) fasthttp.RequestHandler
	WareCommon(next fasthttp.RequestHandler) fasthttp.RequestHandler
	WareSecurity(next fasthttp.RequestHandler, security *jwt.Security) fasthttp.RequestHandler
}

type authMiddleWare struct {
	IMiddleWare
	storage      *storage.JwtStorage
	service      *Service
	passwordHash *security.PasswordHash
}

func CORS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)
		ctx.Response.Header.Set("Access-Control-Expose-Headers", corsExposeHeaders)
		next(ctx)
	}
}

func NewMiddleWare(storage *storage.JwtStorage, service *Service) *authMiddleWare {
	return &authMiddleWare{storage: storage, passwordHash: &security.PasswordHash{}, service: service}
}

func (m authMiddleWare) WareSecurity(next fasthttp.RequestHandler, security *jwt.Security) fasthttp.RequestHandler {
	// Если токен валидный, то делаем, что делали
	jwtStorage := security.JwtStorage()
	return func(ctx *fasthttp.RequestCtx) {
		defer m.wareAll(ctx)
		// получаем аксес токен
		request := &ctx.Request
		// проверяем его на валидность
		tokenPairDetails, err := security.ExtractTokenPairDetails(request)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}
		// если он валидный, то просто делаем, что делали дальше
		accessToken := tokenPairDetails.AccessToken
		claims := security.GetTokenClaims(accessToken).(jwt_go.MapClaims)
		if claims != nil {
			if accountId, ok := claims[security.ContextKey().AccountIdKey]; ok {
				if accountId != nil {
					if isValidAccessToken := security.TokenValid(accessToken); isValidAccessToken {
						ctx.SetUserValue(security.ContextKey().AccountIdKey, accountId)
						next(ctx)
						return
					} else {
						// если он не валидный, то получаем рефреш токен

						// если он валидный, генерируем новую пару -> устанавливаем в хедерсы новую пару, при этом делаем, что делали
						if isValidRefreshToken := security.TokenValid(tokenPairDetails.RefreshToken); isValidRefreshToken {
							claims := accessToken.Claims.(jwt_go.MapClaims)
							accountId := int64(claims[security.ContextKey().AccountIdKey].(float64))

							// если он в списке протухших, то считаем его тоже невалидным
							isExpired, err := jwtStorage.IsExistExpiredToken(tokenPairDetails.RefreshToken.Raw)
							if err != nil || isExpired {
								ctx.SetStatusCode(fasthttp.StatusUnauthorized)
								return
							}

							tokenPair, err := security.CreateTokenPair(accountId)
							err = jwtStorage.PutExpiredToken(accessToken.Raw, tokenPairDetails.RefreshToken.Raw)
							if err != nil {
								return
							}
							if err != nil {
								ctx.SetStatusCode(fasthttp.StatusUnauthorized)
								return
							}
							ctx.Response.Header.Set(accessTokenKey, tokenPair.AccessToken)
							ctx.Response.Header.Set(refreshTokenKey, tokenPair.RefreshToken)
							ctx.SetUserValue(security.ContextKey().AccountIdKey, accountId)
							next(ctx)
						}
					}
				}
			}
		}

		// если он не валидный, то прекращаем работу
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return
	}
}

func (m authMiddleWare) WareCommon(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer m.wareAll(ctx)
		next(ctx)
	}
}

func (m authMiddleWare) WareLogin(next fasthttp.RequestHandler, security *jwt.Security) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer m.wareAll(ctx)

		var lReq middleware_auth.LoginParam
		if err := json.Unmarshal(ctx.Request.Body(), &lReq); err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		acc, err := m.service.AccountService.GetByEmail(lReq.Email)
		if acc != nil {
			if err == nil && m.passwordHash.ComparePassword(*acc.Password, lReq.Password) {
				tokenPair, err := security.CreateTokenPair(int64(acc.Id))
				if err != nil {
					ctx.SetStatusCode(fasthttp.StatusBadRequest)
					return
				}
				ctx.Response.Header.Set(accessTokenKey, tokenPair.AccessToken)
				ctx.Response.Header.Set(refreshTokenKey, tokenPair.RefreshToken)
				return
			}
		}
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		ctx.Response.Header.Set(accessTokenKey, "")
		ctx.Response.Header.Set(refreshTokenKey, "")
	}
}

func (m authMiddleWare) WareLogout(next fasthttp.RequestHandler, security *jwt.Security) fasthttp.RequestHandler {
	jwtStorage := security.JwtStorage()
	return func(ctx *fasthttp.RequestCtx) {
		defer m.wareAll(ctx)
		// получаем аксес токен
		request := &ctx.Request
		// проверяем его на валидность
		tokenPairDetails, err := security.ExtractTokenPairDetails(request)
		if err != nil {
			return
		}
		err = jwtStorage.PutExpiredToken(tokenPairDetails.AccessToken.Raw, tokenPairDetails.RefreshToken.Raw)
		if err != nil {
			return
		}
	}
}

func (m authMiddleWare) wareAll(outerCtx *fasthttp.RequestCtx) {
	outerCtx.SetContentType(contentTypeJson)
}

func (m authMiddleWare) WareRegistry(next fasthttp.RequestHandler, security *jwt.Security) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var accountRegistry middleware_auth.RegistryParam
		if err := json.Unmarshal(ctx.Request.Body(), &accountRegistry); err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		acc, err := m.service.GetByEmail(accountRegistry.Email)
		next(ctx)
		if acc != nil {
			ctx.SetStatusCode(fasthttp.StatusOK)
			response := http_common.CreateErrorMessage("Имя занято")
			res, err := json.Marshal(&response)
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}
			ctx.SetBody(res)
			return
		}

		hash, err := m.passwordHash.CreateHash(accountRegistry.Password)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
		accountRegistry.Password = hash

		// Создаем аккаунт
		accountId, err := m.service.AccountService.Create(domains.Account{
			Password:  &accountRegistry.Password,
			Email:     &accountRegistry.Email,
			Phone:     &accountRegistry.Phone,
			Createdat: time.Now(),
			Blocked:   false,
		})

		// Создаем пользователя
		if err == nil {
			user, err := m.service.UserService.Create(domains.User{})
			if err == nil {
				id := *user.Id
				userAccountId, err := m.service.UserAccountService.Create(domains.UserAccount{
					UserId:    &id,
					AccountId: &accountId,
				})
				if err != nil {
					go func() {
						// Не проверяем, что удалилось
						_, _ = m.service.UserAccountService.DeleteById(userAccountId)
						_, _ = m.service.UserService.DeleteById(id)
						_, _ = m.service.AccountService.DeleteById(accountId)
						// TODO Удалить баланс !!!!
					}()
				}
			}

			m.service.TeamService.AddUser(1, *user.Id)
		}
	}
}
