package session

import (
	"Kavka/config"
	"Kavka/database"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const PHONE = "sample_phone_number"

var STATIC_ID = primitive.NewObjectID()

type MyTestSuite struct {
	suite.Suite
	session      *Session
	generatedOTP int
}

func (s *MyTestSuite) SetupSuite() {
	// Load configs
	configs := config.Read()

	// Init Redis
	var redisClient = database.GetRedisDBInstance(configs.Redis)

	s.session = NewSession(redisClient, configs.App.Auth)
}

func (s *MyTestSuite) TestLogin() {
	otp, loginErr := s.session.Login(PHONE)

	s.generatedOTP = otp

	s.NoError(loginErr)
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
			tokens, ok := s.session.VerifyOTP(PHONE, tt.otp, STATIC_ID)

			assert.Equal(s.T(), ok, tt.result, fmt.Sprintf("Invalid OTP: %d", tt.otp))

			if ok {
				assert.NotEmpty(t, tokens.AccessToken, "Token is empty")
				assert.NotEmpty(t, tokens.RefreshToken, "Token is empty")
			}
		})
	}
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}
