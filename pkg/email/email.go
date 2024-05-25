package email

import (
	"fmt"

	"github.com/kavkaco/Kavka-Core/config"
	"go.uber.org/zap"
)

const TEMPLATE_FORMAT = "txt"

type EmailOtp struct {
	Logger        *zap.Logger
	Config        *config.Email
	TemplatesPath string
}

func NewEmailService(logger *zap.Logger, configs *config.Email, templatesPath string) *EmailOtp {
	return &EmailOtp{logger, configs, templatesPath}
}

func (s *EmailOtp) SendEmail(template string, receivers []string, args interface{}) error {
	if config.CurrentEnv == config.Development {
		fmt.Println("------ Email Sent ------")
		fmt.Println(args)
		fmt.Println("-----------------------")
	}

	return nil
}

func (s *EmailOtp) Template(name string, args interface{}) (string, error) {
	panic("not implemented")
}
