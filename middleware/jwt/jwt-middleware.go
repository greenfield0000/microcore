package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	httpcommon "github.com/greenfield0000/microcore/common"
	"os"
)

type authHeader struct {
	Authorization string `reqHeader:"Authorization"`
	RefreshToken  string `reqHeader:"refresh-token"`
}

// AuthRequired требуется авторизация при вызове ручки
func AuthRequired() fiber.Handler {

	return jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_ACCESS_SECRET")),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			header := new(authHeader)
			if err := c.ReqHeaderParser(&header); err != nil {
				return err
			}

			_ = fmt.Errorf("header Authorization is %s", header.Authorization)
			_ = fmt.Errorf("map RefreshToken is %s", header.RefreshToken)

			if err.Error() == "Missing or malformed JWT" {
				return c.Status(fiber.StatusUnauthorized).JSON(httpcommon.CreateErrorMessage("Пользователь не авторизован"))
			}
			return c.Status(fiber.StatusUnauthorized).JSON(httpcommon.CreateErrorMessage("Сессия пользователя истекла"))
		},
	})
}
