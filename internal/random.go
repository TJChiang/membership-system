package internal

import (
	"crypto/rand"
	"encoding/base64"
)

// https://gist.github.com/dopey/c69559607800d2f2f90b1b1ed4e550fb

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomStringURLSafe(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}
