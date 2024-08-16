package search

import (
	"testing"

	"github.com/kavkaco/Kavka-Core/utils/vali"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ValiTestSuite struct {
	suite.Suite
	validator *vali.Vali
}

func (s *ValiTestSuite) SetupSuite() {
	s.validator = vali.Validator()
}

func (s *ValiTestSuite) TestSearchValidation() {
	testCases := []struct {
		Input *SearchValidation
		Valid bool
	}{
		{
			Input: &SearchValidation{Input: ""},
			Valid: false,
		},
		{
			Input: &SearchValidation{Input: " "},
			Valid: false,
		},
		{
			Input: &SearchValidation{Input: "Sample"},
			Valid: true,
		},
	}

	for _, tc := range testCases {
		varrors := s.validator.Validate(tc.Input)
		if !tc.Valid {
			require.NotEqual(s.T(), len(varrors), 0, "It seems test case is valid actually")
			continue
		}

		require.Len(s.T(), varrors, 0)
	}
}

func TestValiSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(ValiTestSuite))
}
