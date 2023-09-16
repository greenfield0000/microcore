package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	httpcommon "github.com/greenfield0000/microcore/common"
	"github.com/greenfield0000/microcore/utils/jwtutl"
	"os"
)

const (
	authorizationHeaderName = fiber.HeaderAuthorization
	refreshTokenHeaderName  = "refresh-token"
	bearer                  = "Bearer "
)

// AuthRequired требуется авторизация при вызове ручки
func AuthRequired() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(os.Getenv("JWT_ACCESS_SECRET")),
		ErrorHandler: ErrorHandler,
	})
}

// ErrorHandler ...
func ErrorHandler(c *fiber.Ctx, err error) error {
	header := &jwtutl.TokenPair{}
	if parseErr := c.ReqHeaderParser(header); parseErr != nil {
		return parseErr
	}
	jwtManager := jwtutl.NewCommonJwtManager()

	oldRefreshToken, pErr := jwtManager.ParseToken(jwtutl.REFRESH_TOKEN_TYPE, header.RefreshToken)
	if pErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(httpcommon.CreateErrorMessage("Сессия пользователя истекла"))
	}

	oldAccessToken, _ := jwtManager.ParseToken(jwtutl.ACCESS_TOKEN_TYPE, header.AccessToken[len(bearer)+1:])
	tokenPair, err := jwtManager.RefreshTokenPair(jwtutl.JwtTokenPair{AccessToken: oldAccessToken, RefreshToken: oldRefreshToken})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(httpcommon.CreateErrorMessage("Сессия пользователя истекла"))
	}

	c.Append(authorizationHeaderName, fmt.Sprintf("%s %s", bearer, tokenPair.AccessToken))
	c.Append(refreshTokenHeaderName, fmt.Sprintf("%s", tokenPair.RefreshToken))

	return c.Next()
}
