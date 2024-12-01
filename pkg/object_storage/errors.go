package object_storage

import "errors"

var (
	ErrMaxFilesizeExceeded = errors.New("max file size exceeded")
	ErrInvalidFileFormat   = errors.New("invalid file format")
)
