package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	httpcommon "github.com/greenfield0000/microcore/common"
	"os"
)

// AuthRequired требуется авторизация при вызове ручки
func AuthRequired() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_ACCESS_SECRET")),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			m := make(map[interface{}]interface{})

			if err := c.ReqHeaderParser(m); err != nil {
				return err
			}

			_ = fmt.Errorf("map is %s", m)

			if err.Error() == "Missing or malformed JWT" {
				return c.Status(fiber.StatusUnauthorized).JSON(httpcommon.CreateErrorMessage("Пользователь не авторизован"))
			}
			return c.Status(fiber.StatusUnauthorized).JSON(httpcommon.CreateErrorMessage("Сессия пользователя истекла"))
		},
	})
}
