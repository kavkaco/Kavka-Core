package uploader

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kavkaco/Kavka-Core/config"
	"github.com/kavkaco/Kavka-Core/utils/random"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	RANDOM_FILENAME_LENGTH = 25
	UPLOAD_TMP_DIR         = "/tmp/uploads"
)

var ErrMaxFileSize = errors.New("maximum file size")

type (
	Service      struct{ minioClient *minio.Client }
	FileUploaded struct {
		Name string
		Size int64
	}
)

func NewUploaderService(config *config.Config) *Service {
	minioCredentials := config.MinIO

	endpoint := minioCredentials.Endpoint
	accessKeyID := minioCredentials.AccessKey
	secretAccessKey := minioCredentials.SecretKey

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}

	return &Service{minioClient}
}

func (s *Service) UploadFile(bucketName string, filePath string, maxFileSize *int64) (*FileUploaded, error) {
	// Collect objectName, contentType and filePath
	fileInfo, statErr := os.Stat(filePath)
	if statErr != nil {
		return nil, statErr
	}

	if maxFileSize != nil {
		if fileInfo.Size() > *maxFileSize {
			return nil, ErrMaxFileSize
		}
	}

	objectName := random.GenerateRandomFileName(RANDOM_FILENAME_LENGTH)
	contentType := filepath.Ext(filePath)

	// Upload the file
	_, err := s.minioClient.FPutObject(context.Background(), bucketName,
		objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return nil, err
	}

	return &FileUploaded{
		Name: objectName,
		Size: fileInfo.Size(),
	}, err
}

func (s *Service) DeleteFile(bucketName string, objectName string) error {
	// Delete the file
	opts := minio.RemoveObjectOptions{GovernanceBypass: true}

	err := s.minioClient.RemoveObject(context.Background(), bucketName, objectName, opts)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GenerateTMPFilePath(fileName string) string {
	return fmt.Sprintf("%s/..%s/%s", config.ConfigsDirPath(), UPLOAD_TMP_DIR, fileName)
}
