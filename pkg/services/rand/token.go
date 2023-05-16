package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RememberTokenBytes = 32

func RememberToken() (string, error) {
	return StringEncode(RememberTokenBytes)
}

func StringEncode(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Bytes will generate n random bytes for remember tokens
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err

	}
	return b, nil
}

// Validation that Token has been createated at least with 32bytes lenght
func NBytes(base64String string) (int, error) {
	b, err := base64.URLEncoding.DecodeString(base64String)
	if err != nil {
		return -1, nil
	}
	return len(b), nil
}
