package email

import (
	"fmt"
	"net/smtp"

	"github.com/flosch/pongo2"
	"github.com/kavkaco/Kavka-Core/config"
)

type EmailService interface {
	SendResetPasswordEmail(recipientEmail, url, name, expiry string) error
	SendVerificationEmail(recipientEmail, url string) error
}

const TemplateFormat = "html"

type emailOtp struct {
	configs       *config.Email
	templatesPath string
}
type emailMessage struct {
	template string
	receiver []string
	args     map[string]interface{}
	subject  string
}

func NewEmailService(configs *config.Email, templatesPath string) EmailService {
	return &emailOtp{configs, templatesPath}
}

func newEmailMessage(template, subject string, args map[string]interface{}, receiver []string) *emailMessage {
	return &emailMessage{
		template: template,
		subject:  subject,
		args:     args,
		receiver: receiver,
	}
}

func (s *emailOtp) readTemplate(template string) *pongo2.Template {
	templateFile := s.templatesPath + "/" + template
	return pongo2.Must(pongo2.FromFile(templateFile))
}

func (s *emailOtp) sendEmail(msg *emailMessage) error {
	pongoTemplate := s.readTemplate(msg.template)
	ctx := make(pongo2.Context)
	for key, value := range msg.args {
		ctx[key] = value
	}

	body, err := pongoTemplate.Execute(ctx)
	if err != nil {
		return err
	}

	emailMessage := fmt.Sprintf("Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+body, msg.subject)

	auth := smtp.PlainAuth("", s.configs.SenderEmail, s.configs.Password, s.configs.Host)
	err = smtp.SendMail(s.configs.Host+":"+s.configs.Port, auth, s.configs.SenderEmail, msg.receiver, []byte(emailMessage))
	if err != nil {
		return err
	}
	return nil
}

func (s *emailOtp) SendVerificationEmail(recipientEmail, url string) error {
	msg := newEmailMessage(
		"verification_email.html",
		"Verify Account",
		map[string]interface{}{"url": url},
		[]string{recipientEmail},
	)
	err := s.sendEmail(msg)
	if err != nil {
		return err
	}

	return nil
}

func (s *emailOtp) SendResetPasswordEmail(recipientEmail, url, name, expiry string) error {
	msg := newEmailMessage(
		"reset_password.html",
		"Reset Password",
		map[string]interface{}{"name": name, "url": url, "expiry": expiry},
		[]string{recipientEmail},
	)
	err := s.sendEmail(msg)
	if err != nil {
		return err
	}
	return nil
}
