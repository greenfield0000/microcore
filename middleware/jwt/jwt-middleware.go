package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	httpcommon "github.com/greenfield0000/microcore/common"
	"os"
)

const (
	authorizationHeaderName = fiber.HeaderAuthorization
	refreshTokenHeaderName  = "refresh-token"
	bearer                  = "Bearer "
)

type AuthHeader struct {
	Authorization string `reqHeader:"Authorization"`
	RefreshToken  string `reqHeader:"refresh-token"`
}

// AuthRequired требуется авторизация при вызове ручки
func AuthRequired() fiber.Handler {
	jwtAccessSecret := []byte(os.Getenv("JWT_ACCESS_SECRET"))
	jwtRefreshSecret := []byte(os.Getenv("JWT_REFRESH_SECRET"))
	return jwtware.New(jwtware.Config{
		SigningKey: jwtAccessSecret,
		ErrorHandler: func(c *fiber.Ctx, err error) error {

			// пробуем сначала обновить токены
			// смотрим на рефреш токен

			header := new(AuthHeader)
			if parseErr := c.ReqHeaderParser(header); parseErr != nil {
				return parseErr
			}

			refreshToken, parseErr := jwt.Parse(header.RefreshToken, func(token *jwt.Token) (interface{}, error) {
				return jwtRefreshSecret, nil
			})

			if parseErr != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(httpcommon.CreateErrorMessage("Сессия пользователя истекла"))
			}

			_, ok := refreshToken.Claims.(jwt.MapClaims)
			if ok && refreshToken.Valid {
				c.Append(authorizationHeaderName, fmt.Sprintf("%s %s", bearer, "new Access Token"))
				c.Append(refreshTokenHeaderName, fmt.Sprintf("%s %s", bearer, "new Refresh Token"))
				return nil
			}

			if err.Error() == "Missing or malformed JWT" {
				return c.Status(fiber.StatusUnauthorized).JSON(httpcommon.CreateErrorMessage("Пользователь не авторизован"))
			}
			return c.Status(fiber.StatusUnauthorized).JSON(httpcommon.CreateErrorMessage("Сессия пользователя истекла"))
		},
	})
}
