package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"himatro-api/internal/util"
)

func Decrypt(encrypted string) (string, error) {
	block, err := aes.NewCipher([]byte(secret_key))

	if err != nil {
		util.LogErr("ERROR", "Decrypt failed to create new chiper", err.Error())
		return "", err
	}

	cipherText, err := decode(encrypted)

	if err != nil {
		util.LogErr("ERROR", "Decrypt failed to decode", err.Error())
		return "", err
	}

	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}

func decode(encrypted string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)

	if err != nil {
		util.LogErr("ERROR", "Decode process failed", err.Error())
		return []byte(""), errors.New("authentication process failed")
	}

	return data, nil
}
