package email

import (
	"fmt"
)

type emailDevelopmentManager struct{}

func NewEmailDevelopmentService() EmailService {
	return &emailDevelopmentManager{}
}

func (e *emailDevelopmentManager) SendResetPasswordEmail(recipientEmail string, url string, name string, expiry string) error {
	fmt.Println("====== EMAIL SENT ====== ")
	fmt.Printf("=== Recipient Email: %s\n", recipientEmail)
	fmt.Printf("=== Redirect Url: %s\n", url)
	fmt.Printf("=== Name: %s\n", name)
	fmt.Printf("=== Expiry: %s\n", expiry)

	return nil
}

func (e *emailDevelopmentManager) SendVerificationEmail(recipientEmail string, url string) error {
	fmt.Println("====== EMAIL SENT ====== ")
	fmt.Printf("=== Recipient Email: %s\n", recipientEmail)
	fmt.Printf("=== Redirect Url: %s\n", url)

	return nil
}
