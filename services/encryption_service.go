package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"

	"golang.org/x/crypto/bcrypt"
)

// USERS ULTIMATELY REQUIRE AN APPLICATION NAME AND PASSWORD FOR
func Encrypt(data []byte, passphrase string) (string, error) {
	block, _ := aes.NewCipher([]byte(writeHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return string(ciphertext), nil
}

func Decrypt(data []byte, passphrase string) (string, error) {
	key := []byte(writeHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func HashWithBcrypt(key string) (string, error) {
	data := []byte(key)
	cost := 10
	if bytesHash, err := bcrypt.GenerateFromPassword(data, cost); err != nil {
		return "", err
	}
	return string(bytesHash), nil
}

func CompareWithBcrypt(hash, text string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(text))
	return err == nil
}

func writeHash(key string) string {
	enc := md5.New()
	enc.Write([]byte(key))
	return hex.EncodeToString(enc.Sum(nil))
}
