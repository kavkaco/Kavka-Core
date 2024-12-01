package object_storage

import (
	"bytes"
	"context"
	"os"

	"github.com/kavkaco/Kavka-Core/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const TEMP_DIR = "tmp"

type Object struct {
	Name string
	Size int
}

type Service struct {
	minioClient *minio.Client
	TempDirPath string
}

func New(cfg *config.Config) (*Service, error) {
	minioClient, err := minio.New(cfg.MinIO.Url, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIO.AccessKey, cfg.MinIO.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	tempDirPath := wd + "/" + TEMP_DIR + "/"

	return &Service{minioClient, tempDirPath}, nil
}

func (s *Service) Upload(ctx context.Context, bucketName string, objectName string, fileData []byte, fileSize int64, contentType string) (*Object, error) {
	br := bytes.NewReader(fileData)

	// Upload the file
	_, err := s.minioClient.PutObject(ctx, bucketName,
		objectName, br, fileSize, minio.PutObjectOptions{
			ContentType: contentType,
		})
	if err != nil {
		return nil, err
	}

	return &Object{
		Name: objectName,
		Size: len(fileData),
	}, err
}

func (s *Service) Delete(ctx context.Context, bucketName string, objectName string) error {
	opts := minio.RemoveObjectOptions{GovernanceBypass: true}

	err := s.minioClient.RemoveObject(ctx, bucketName, objectName, opts)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Get(ctx context.Context, bucketName string, objectName string) (string, string, error) {
	filePath := s.TempDirPath + objectName

	err := s.minioClient.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		return "", "", err
	}

	stat, err := s.minioClient.StatObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return "", "", err
	}

	return filePath, stat.ContentType, nil
}
