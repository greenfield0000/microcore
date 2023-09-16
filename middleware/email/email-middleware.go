package middleware

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/greenfield0000/microcore/bussines/service"
	httpcommon "github.com/greenfield0000/microcore/common"
	"github.com/greenfield0000/microcore/utils/jwtutl"
)

var emailVerifierCheckerIsNull = fmt.Errorf("Необходимо передать кинфигурацию для проверки")

// EmailVerifyCheckConfig базовая конфигурация проверки на верифицированную почту1
type EmailVerifyCheckConfig struct {
	VerificationService service.EmailVerifierService
	CheckFunction       func(c *fiber.Ctx, service service.EmailVerifierService) error
}

// newEmailVerifyCheckConfig Глобальная конфигурация верификатора почты
func newEmailVerifyCheckConfig(s *service.Service) EmailVerifyCheckConfig {
	return EmailVerifyCheckConfig{
		VerificationService: s.EmailVerifierService,
		CheckFunction: func(c *fiber.Ctx, service service.EmailVerifierService) error {
			user := c.Locals(jwtutl.UserLocalNameKey)
			if user == nil {
				return errors.New("Для выполнения данного действия требуется подтвердить электронную почту")
			}

			claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
			email := claims["email"].(string)
			if ok, _ := service.IsVerifyByEmail(email); !ok {
				return errors.New("Для выполнения данного действия требуется подтвердить электронную почту")
			}
			return nil
		},
	}
}

// MailMustBeConfirmed почта должна быть подтверждена
func MailMustBeConfirmed(service *service.Service) fiber.Handler {
	config := newEmailVerifyCheckConfig(service)
	return func(c *fiber.Ctx) error {
		verification := config.VerificationService
		if verification == nil {
			return emailVerifierCheckerIsNull
		}

		checkFN := config.CheckFunction

		if config.CheckFunction == nil {
			return emailVerifierCheckerIsNull
		}

		err := checkFN(c, config.VerificationService)

		if err != nil {
			return c.JSON(httpcommon.CreateErrorMessage(err.Error()))
		}

		return c.Next()
	}
}
