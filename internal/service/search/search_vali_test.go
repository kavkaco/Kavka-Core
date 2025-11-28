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
		Input *searchValidation
		Valid bool
	}{
		{
			Input: &searchValidation{Input: ""},
			Valid: false,
		},
		{
			Input: &searchValidation{Input: " "},
			Valid: false,
		},
		{
			Input: &searchValidation{Input: "Sample"},
			Valid: true,
		},
	}

	for _, tc := range testCases {
		errs := s.validator.Validate(tc.Input)
		if !tc.Valid {
			require.NotEqual(s.T(), len(errs), 0, "It seems test case is valid actually")
			continue
		}

		require.Len(s.T(), errs, 0)
	}
}

func TestValiSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(ValiTestSuite))
}
