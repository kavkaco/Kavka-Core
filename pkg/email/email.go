package email

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/flosch/pongo2"
	"github.com/kavkaco/Kavka-Core/config"
)

type EmailManager interface {
	sendEmail(msg *emailMessage) error
	readTemplate(template string) *pongo2.Template
	SendWelcomeEmail(recipientEmail, name string) error
	SendResetPasswordEmail(recipientEmail, url, name, exp string) error
	SendVerificationEmail(recipientEmail, url string) error
}

const TemplateFormat = "html"

type emailOtp struct {
	Configs       *config.Email
	TemplatesPath string
}
type emailMessage struct {
	template string
	receiver []string
	args     map[string]interface{}
	subject  string
}

func NewEmailService(Configs *config.Email, templatesPath string) EmailManager {
	return &emailOtp{Configs, templatesPath}
}

func newEmailMessage(template, subject string, args map[string]interface{}, reciver []string) *emailMessage {
	return &emailMessage{
		template: template,
		subject:  subject,
		args:     args,
		receiver: reciver,
	}
}

func (s *emailOtp) readTemplate(template string) *pongo2.Template {
	templateFile := s.TemplatesPath + "/" + template
	tpl := pongo2.Must(pongo2.FromFile(templateFile))
	return tpl
}

func (s *emailOtp) sendEmail(msg *emailMessage) error {
	if config.CurrentEnv == config.Development {
		log.Println("Email sent")
		log.Println(msg)
		return nil
	}
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
	auth := smtp.PlainAuth("", s.Configs.SenderEmail, s.Configs.Password, s.Configs.Host)
	err = smtp.SendMail(s.Configs.Host+":"+s.Configs.Port, auth, s.Configs.SenderEmail, msg.receiver, []byte(emailMessage))
	if err != nil {
		return err
	}
	log.Println("Verification code email sent successfully!")
	return nil
}

func (s *emailOtp) SendWelcomeEmail(recipientEmail, name string) error {
	msg := newEmailMessage(
		"welcome_message.html",
		"Welcome",
		map[string]interface{}{"name": name},
		[]string{recipientEmail},
	)
	err := s.sendEmail(msg)
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

func (s *emailOtp) SendResetPasswordEmail(recipientEmail, url, name, exp string) error {
	msg := newEmailMessage(
		"submit_reset_password.html",
		"Reset Password",
		map[string]interface{}{"name": name, "url": url, "expiry": exp},
		[]string{recipientEmail},
	)
	err := s.sendEmail(msg)
	if err != nil {
		return err
	}
	return nil
}
