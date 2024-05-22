package email

import (
	"os"
	"testing"

	"github.com/kavkaco/Kavka-Core/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSendEmail(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // nolint

	// Get wd
	wd, _ := os.Getwd()
	templatesPath := wd + "/../../app/views/email/"

	// Load configs
	config := config.Read()

	receivers := []string{"+989368392346"}

	emailService := NewEmailService(logger, &config.Email, templatesPath)

	template, templateErr := emailService.Template("code_sent", struct{ Code int }{
		Code: 123456,
	})
	assert.NoError(t, templateErr)

	err := emailService.SendEmail(template, receivers)
	assert.NoError(t, err)
}
