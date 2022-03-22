package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"himatro-api/internal/config"
)

var secret_key = config.SecretKey()
var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher([]byte(secret_key))

	if err != nil {
		return "", err
	}

	plainTextAsByte := []byte(plainText)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainTextAsByte))
	cfb.XORKeyStream(cipherText, plainTextAsByte)

	return encode(cipherText), nil
}

func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
