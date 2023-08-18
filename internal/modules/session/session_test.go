package session

import (
	"Kavka/config"
	"Kavka/database"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const CONFIG_PATH = "/../../../config/configs.yml"
const PHONE = "sample_phone_number"

type MyTestSuite struct {
	suite.Suite
	session      *Session
	generatedOTP int
}

func (s *MyTestSuite) SetupSuite() {
	// Get wd
	var wd, _ = os.Getwd()

	// Load configs
	var configs, configsErr = config.Read(wd + CONFIG_PATH)
	if configsErr != nil {
		panic(configsErr)
	}

	// Init Redis
	var redisClient = database.GetRedisDBInstance(configs.Redis)

	s.session = NewSession(redisClient, configs.App.Auth)
}

func (s *MyTestSuite) TestLogin() {
	otp, loginErr := s.session.Login(PHONE)

	s.generatedOTP = otp

	s.NoError(loginErr)
	s.T().Log(otp)
}

func (s *MyTestSuite) TestVerifyOTP() {
	cases := []struct {
		name   string
		otp    int
		result bool
	}{
		{
			name:   "valid",
			otp:    s.generatedOTP,
			result: true,
		},
		{
			name:   "not_valid",
			otp:    0,
			result: false,
		},
	}

	for _, tt := range cases {
		s.T().Run(tt.name, func(t *testing.T) {
			result := s.session.VerifyOTP(PHONE, tt.otp)
			assert.Equal(s.T(), result, tt.result)
		})
	}
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}
