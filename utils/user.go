package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenToken(size int) string {
	token := make([]byte, size)
	_, _ = rand.Read(token)
	return base64.StdEncoding.EncodeToString(token)
}
