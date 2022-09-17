package middlerware

import (
	jwt_go "github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
	"microcore/security"
	"microcore/security/jwt"
	"microcore/security/jwt/storage"
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

type IMiddleWare interface {
	wareAll(ctx *fasthttp.RequestCtx)
	WareLogin(next fasthttp.RequestHandler, security *jwt.Security) fasthttp.RequestHandler
	WareLogout(next fasthttp.RequestHandler, security *jwt.Security) fasthttp.RequestHandler
	WareRegistry(next fasthttp.RequestHandler, security *jwt.Security) fasthttp.RequestHandler
	WareCommon(next fasthttp.RequestHandler) fasthttp.RequestHandler
	WareSecurity(next fasthttp.RequestHandler, security *jwt.Security) fasthttp.RequestHandler
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

type Service interface{}

type authMiddleWare struct {
	IMiddleWare
	storage      *storage.JwtStorage
	service      *Service
	passwordHash *security.PasswordHash
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

		claims := tokenPairDetails.AccessToken.Claims.(jwt_go.MapClaims)
		accountId := int64(claims[security.ContextKey().AccountIdKey].(float64))
		// если он валидный, то просто делаем, что делали дальше
		if isValidAccessToken := security.TokenValid(tokenPairDetails.AccessToken); isValidAccessToken {
			ctx.SetUserValue(security.ContextKey().AccountIdKey, accountId)
			next(ctx)
			return
		} else {
			// если он не валидный, то получаем рефреш токен

			// если он валидный, генерируем новую пару -> устанавливаем в хедерсы новую пару, при этом делаем, что делали
			if isValidRefreshToken := security.TokenValid(tokenPairDetails.RefreshToken); isValidRefreshToken {
				// если он в списке протухших, то считаем его тоже невалидным
				isExpired, err := jwtStorage.IsExistExpiredToken(tokenPairDetails.RefreshToken.Raw)
				if err != nil || isExpired {
					ctx.SetStatusCode(fasthttp.StatusUnauthorized)
					return
				}

				tokenPair, err := security.CreateTokenPair(accountId)
				err = jwtStorage.PutExpiredToken(tokenPairDetails.AccessToken.Raw, tokenPairDetails.RefreshToken.Raw)
				if err != nil {
					ctx.SetStatusCode(fasthttp.StatusUnauthorized)
					return
				}
				ctx.Response.Header.Set(accessTokenKey, tokenPair.AccessToken)
				ctx.Response.Header.Set(refreshTokenKey, tokenPair.RefreshToken)
				ctx.SetUserValue(security.ContextKey().AccountIdKey, accountId)
				next(ctx)
			} else {
				// если он не валидный, то прекращаем работу
				ctx.SetStatusCode(fasthttp.StatusUnauthorized)
				return
			}
		}
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
		next(ctx)
		code := ctx.Response.StatusCode()
		switch code {
		case 200:
			accountId := ctx.UserValue(security.ContextKey().AccountIdKey).(uint64)
			tokenPair, err := security.CreateTokenPair(int64(accountId))
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusBadRequest)
				return
			}
			ctx.Response.Header.Set(accessTokenKey, tokenPair.AccessToken)
			ctx.Response.Header.Set(refreshTokenKey, tokenPair.RefreshToken)
			return
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
		next(ctx)
	}
}
