package email

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/flosch/pongo2"
	"github.com/kavkaco/Kavka-Core/config"
	"go.uber.org/zap"
)

type EmailManager interface {
	sendEmail(template string, receiver []string, args interface{}) error
	readTemplate(template string) *pongo2.Template
	SendWelcomeEmail(recipientEmail, name string) error
	SendResetPasswordEmail(recipientEmail, url, name, exp string) error
	SendVerificationEmail(recipientEmail, otp string) error
}

const TemplateFormat = "html"

type EmailOtp struct {
	Logger        *zap.Logger
	Configs       *config.Email
	TemplatesPath string
}

func NewEmailService(logger *zap.Logger, Configs *config.Email, templatesPath string) EmailManager {
	return &EmailOtp{logger, Configs, templatesPath}
}
func (s *EmailOtp) readTemplate(template string) *pongo2.Template {
	tpl := pongo2.Must(pongo2.FromFile(s.TemplatesPath + "/" + template))
	return tpl
}
func (s *EmailOtp) sendEmail(template string, receiver []string, args interface{}) error {
	if config.CurrentEnv == config.Development {
		log.Println("Email sent")
		log.Println(args)
		return nil
	}
	pongoTemplate := s.readTemplate(template)

	body, err := pongoTemplate.Execute(args.(pongo2.Context))
	if err != nil {
		return err
	}
	message := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Verification Code\r\n\r\n%s", "Kafka", receiver, body))
	auth := smtp.PlainAuth("", s.Configs.SenderEmail, s.Configs.Password, s.Configs.Host)
	err = smtp.SendMail(s.Configs.Host+":"+s.Configs.Port, auth, s.Configs.SenderEmail, receiver, message)
	if err != nil {
		return err
	}
	log.Println("Verification code email sent successfully!")
	return nil

}
func (s *EmailOtp) SendWelcomeEmail(recipientEmail, name string) error {
	err := s.sendEmail("verification_code.html", []string{recipientEmail}, struct{ Name string }{Name: name})
	if err != nil {
		return err
	}
	return nil
}
func (s *EmailOtp) SendVerificationEmail(recipientEmail, otp string) error {
	err := s.sendEmail("verification_code.html", []string{recipientEmail}, struct{ code string }{code: otp})
	if err != nil {
		return err
	}
	return nil
}
func (s *EmailOtp) SendResetPasswordEmail(recipientEmail, url, name, exp string) error {
	err := s.sendEmail("submit_reset_password.html", []string{recipientEmail}, struct{ name, url, exp string }{name: name, url: url, exp: exp})
	if err != nil {
		return err
	}
	return nil
}
