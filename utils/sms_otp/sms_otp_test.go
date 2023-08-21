package sms_otp

import (
	"Kavka/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const CONFIG_PATH = "/../../config/configs.yml"

func TestSendSms(t *testing.T) {
	// Get wd
	var wd, _ = os.Getwd()
	var templatesPath = wd + "/../../app/views/sms/"

	// Load configs
	var configs, configsErr = config.Read(wd + CONFIG_PATH)
	if configsErr != nil {
		panic(configsErr)
	}

	receivers := []string{"+989368392346"}

	smsOTP := NewSMSOtpService(&configs.SMS, templatesPath)

	template, templateErr := smsOTP.Template("code_sent", struct{ Code int }{
		Code: 123456,
	})
	assert.NoError(t, templateErr)

	err := smsOTP.SendSMS(template, receivers)
	assert.NoError(t, err)

	t.Log(template)
}
