package models

import "errors"

var (
	ErrNotFound        = errors.New("email or password was wrong")
	ErrGalleryNotFound = errors.New("Gallery not found")
)
