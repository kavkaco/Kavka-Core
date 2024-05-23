package email

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/kavkaco/Kavka-Core/config"
	"go.uber.org/zap"
)

const TEMPLATE_FORMAT = "txt"

type EmailOtp struct {
	logger        *zap.Logger
	config        *config.Email
	templatesPath string
}

func NewEmailService(logger *zap.Logger, configs *config.Email, templatesPath string) *EmailOtp {
	return &EmailOtp{logger, configs, templatesPath}
}

func (s *EmailOtp) SendEmail(body string, receivers []string) error {
	if config.CurrentEnv == config.Development {
		fmt.Println("------ Email Sent ------")
		fmt.Printf("%s\n", strings.TrimSpace(body))
		fmt.Println("-----------------------")
	}

	return nil
}

// Parses and returns the template.
func (s *EmailOtp) Template(name string, args interface{}) (string, error) {
	filename := fmt.Sprintf("%s/%s.%s", s.templatesPath, name, TEMPLATE_FORMAT)

	fileData, readErr := os.ReadFile(filename)
	if readErr != nil {
		return "", readErr
	}

	renderedFile := new(bytes.Buffer)

	t := template.Must(template.New(name).Parse(string(fileData)))
	err := t.Execute(renderedFile, args)
	if err != nil {
		return "", err
	}

	return renderedFile.String(), nil
}
