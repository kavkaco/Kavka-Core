package sms_service

import (
	"os"
	"testing"

	"github.com/kavkaco/Kavka-Core/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSendSms(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // nolint

	// Get wd
	wd, _ := os.Getwd()
	templatesPath := wd + "/../../app/views/sms/"

	// Load configs
	configs := config.Read()

	receivers := []string{"+989368392346"}

	smsService := NewSmsService(logger, &configs.SMS, templatesPath)

	template, templateErr := smsService.Template("code_sent", struct{ Code int }{
		Code: 123456,
	})
	assert.NoError(t, templateErr)

	err := smsService.SendSMS(template, receivers)
	assert.NoError(t, err)

	t.Log(template)
}
