package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/greenfield0000/microcore/bussines/service"
)

var emailVerifierCheckerIsNull = fmt.Errorf("Необходимо передать кинфигурацию для проверки")

type EmailVerifyCheckConfig struct {
	VerificationService service.EmailVerifierService
	CheckFunction       func(c *fiber.Ctx, srv service.EmailVerifierService) error
}

func New(config EmailVerifyCheckConfig) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		verificator := config.VerificationService
		if verificator == nil {
			return emailVerifierCheckerIsNull
		}

		checkFN := config.CheckFunction

		if config.CheckFunction == nil {
			return emailVerifierCheckerIsNull
		}

		return checkFN(c, config.VerificationService)
	}
}
