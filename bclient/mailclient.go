package bclient

import "github.com/gofiber/fiber/v2"

type EmailClient struct {
	options *ClientOptions
}

// NewEmailClient ...
func NewEmailClient(options *ClientOptions) EmailClient {
	return EmailClient{
		options: options,
	}
}

// SendCode отправка кода на почту
func (e EmailClient) SendCode(c *fiber.Ctx, email string) error {
	a := fiber.AcquireAgent()

	request := a.Request()
	request.SetRequestURI(e.options.RequestUri + "/send-code")
	request.SetTimeout(e.options.RequestTimeout)

	request.Header.Set(fiber.HeaderXRequestID, c.Locals(fiber.HeaderXRequestID).(string))
	request.Header.SetMethod(fiber.MethodPost)

	a.JSON(fiber.Map{
		"email": email,
	})

	if err := a.Parse(); err != nil {
		return err
	}

	return nil
}
