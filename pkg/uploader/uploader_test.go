package uploader

import (
	"os"
	"testing"

	"Kavka/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const BUCKET_NAME = "profile-photos"

type MyTestSuite struct {
	suite.Suite
	uploadService      *Service
	uploadedObjectName string
}

func (s *MyTestSuite) SetupSuite() {
	// Load configs
	configs := config.Read()

	s.uploadService = NewUploaderService(configs)
}

func (s *MyTestSuite) TestA_UploadFile() {
	// Creating a sample txt file
	fileName := "sample_file.txt"
	filePath := s.uploadService.GenerateTMPFilePath(fileName)

	file, fileErr := os.Create(filePath)
	assert.NoError(s.T(), fileErr)
	defer file.Close()
	_, err := file.WriteString("Hello Bucket!\n")
	assert.NoError(s.T(), err)

	// Uploading File
	uploaded, err := s.uploadService.UploadFile(BUCKET_NAME, filePath, nil)
	assert.NoError(s.T(), err)

	// Store object name for next test
	s.uploadedObjectName = uploaded.Name

	// Remove sample txt file
	os.Remove(filePath)
}

func (s *MyTestSuite) TestB_DeleteFile() {
	err := s.uploadService.DeleteFile(BUCKET_NAME, s.uploadedObjectName)
	assert.NoError(s.T(), err)
}

func TestMySuite(t *testing.T) {
	suite.Run(t, new(MyTestSuite))
}
