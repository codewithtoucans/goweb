package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const SessionTokenBytes = 32

func bytes(length int) ([]byte, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	return randomBytes, nil
}

func String(length int) (string, error) {
	randomBytes, err := bytes(length)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}

func SessionToken() (string, error) {
	return String(SessionTokenBytes)
}
