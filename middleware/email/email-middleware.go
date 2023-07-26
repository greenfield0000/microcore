package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/greenfield0000/microcore/bussines/service"
	httpcommon "github.com/greenfield0000/microcore/common"
)

var emailVerifierCheckerIsNull = fmt.Errorf("Необходимо передать кинфигурацию для проверки")

// EmailVerifyCheckConfig базовая конфигурация проверки на верифицированную почту1
type EmailVerifyCheckConfig struct {
	VerificationService service.EmailVerifierService
	CheckFunction       func(c *fiber.Ctx, srv service.EmailVerifierService) error
}

// MailMustBeConfirmed почта должна быть подтверждена
func MailMustBeConfirmed(config EmailVerifyCheckConfig) fiber.Handler {
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
