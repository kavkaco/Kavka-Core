package file_manager

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/kavkaco/Kavka-Core/internal/model"
)

var correctFormats = []string{
	"png", "jpg", "jpeg", "gif",
}

var ErrInvalidFileForma = errors.New("invalid file format")

type FileManager interface {
	SetFile(fileName model.UserID, path string) error
	Write(chunk []byte) error
	Close() error
}

func NewFileManager() FileManager {
	return &serverFileMng{}
}

type serverFileMng struct {
	filePath   string
	outputFile *os.File
}

func (f *serverFileMng) SetFile(fileName, path string) error {
	if f.filePath != "" {
		return nil
	}

	isCorrectFormat := IsCorrectFormat(fileName)
	if !isCorrectFormat {
		return ErrInvalidFileForma
	}

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	f.filePath = filepath.Join(path, fileName)

	file, createErr := os.Create(f.filePath)
	if createErr != nil {
		return err
	}

	f.outputFile = file

	return nil
}

func (f *serverFileMng) Write(chunk []byte) error {
	if f.outputFile == nil {
		return nil
	}
	_, err := f.outputFile.Write(chunk)
	return err
}

func (f *serverFileMng) Close() error {
	return f.outputFile.Close()
}

func IsCorrectFormat(s string) bool {
	regex := regexp.MustCompile(`\.[^.]+$`)
	match := regex.FindString(s)
	if match != "" {
		extension := strings.ToLower(match[1:])
		for _, f := range correctFormats {
			if extension == f {
				return true
			}
		}
	}
	return false
}
