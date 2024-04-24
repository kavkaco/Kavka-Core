package sms_service

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

type SmsService struct {
	logger        *zap.Logger
	configs       *config.SMS
	templatesPath string
}

func NewSmsService(logger *zap.Logger, configs *config.SMS, templatesPath string) *SmsService {
	return &SmsService{logger, configs, templatesPath}
}

// TODO - Write sms service for production.
func (s *SmsService) SendSMS(msg string, receivers []string) error {
	if config.CurrentEnv == config.Development {
		fmt.Println("------ SMS Sent ------")
		fmt.Printf("%s\n", strings.TrimSpace(msg))
		fmt.Println("-----------------------")
	}

	return nil
}

// Parses and returns the template.
func (s *SmsService) Template(name string, args interface{}) (string, error) {
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
