package sms_service

import (
	"os"
	"testing"

	"github.com/kavkaco/Kavka-Core/config"

	"github.com/stretchr/testify/assert"
)

func TestSendSms(t *testing.T) {
	// Get wd
	wd, _ := os.Getwd()
	templatesPath := wd + "/../../app/views/sms/"

	// Load configs
	configs := config.Read()

	receivers := []string{"+989368392346"}

	smsOTP := NewSmsService(&configs.SMS, templatesPath)

	template, templateErr := smsOTP.Template("code_sent", struct{ Code int }{
		Code: 123456,
	})
	assert.NoError(t, templateErr)

	err := smsOTP.SendSMS(template, receivers)
	assert.NoError(t, err)

	t.Log(template)
}
