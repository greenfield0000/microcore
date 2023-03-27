package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/greenfield0000/microcore/bussines/service"
)

var emailVerifierCheckerIsNull = fmt.Errorf("Необходимо передать кинфигурацию для проверки")

type EmailVerifyCheckConfig struct {
	Verificator *service.EmailVerifierService
	Check       func(srv *service.EmailVerifierService) error
}

func New(config EmailVerifyCheckConfig) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		verificator := config.Verificator
		if verificator == nil {
			return emailVerifierCheckerIsNull
		}

		checkFN := config.Check

		if config.Check == nil {
			return emailVerifierCheckerIsNull
		}

		return checkFN(config.Verificator)
	}
}
