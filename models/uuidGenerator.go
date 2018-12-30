package models

import (
	"crypto/rand"
	"encoding/base64"
)

//GenerateTransactionId - generate random string
func GenerateTransactionId(len int) string {
	buff := make([]byte, len)
	rand.Read(buff)
	str := base64.StdEncoding.EncodeToString(buff)
	return str[:len]
}
