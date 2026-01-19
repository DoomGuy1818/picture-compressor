package storage

import "errors"

var (
	ErrPictureNotFound     = errors.New("picture not found")
	ErrPictureDoesNotExist = errors.New("picture does not exist")
)
