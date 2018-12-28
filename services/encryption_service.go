package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

/*
 ULTIMATELY WE WANT TO BUILD A SYSTEM THAT IS REALLY TRUST WORTHY, SO WE'LL
 BE STORING THEIR APP_PASSWORDS AND THE DATABASE ADDRESS. THE LOGIC WILL BE TO ENCRYPT
 THE ADDRESS + THE NAME OF THE APP (SINCE THIS IS UNIQUE), WE'D USE THIS TO GENERATE THE API_KEY
 THEN PERSIST THE HASHED API_KEY FOR THE APPLICATION KEY. THEN TRY TO GET THE
*/

var secret = flag.String("hmacsecret", "Thisshouldnotbeusedever", "password for generating tokens")

type CustomSlice []uint8

func (u CustomSlice) MarshalJSON() ([]byte, error) {
	if u == nil {
		return []byte("null"), nil
	}
	var result string
	result = strings.Join(strings.Fields(fmt.Sprintf("%d", u)), ",")
	return []byte(result), nil
}

func Encrypt(data []byte, passphrase string) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(writeHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}
	// this Passphrase has to be an os.LookUP : Potential vulnerability
	ciphertext := gcm.Seal(nil, generateNonce(passphrase, gcm), data, nil)
	return ciphertext, nil
}

func Decrypt(data []byte, passphrase string) (string, error) {
	key := []byte(writeHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "1", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "2", err
	}
	plaintext, err := gcm.Open(nil, generateNonce(passphrase, gcm), data, nil)
	if err != nil {
		return "3", err
	}
	return string(plaintext), nil
}

func HashWithBcrypt(key string) (string, error) {
	data := []byte(key)
	cost := 10
	bytesHash, err := bcrypt.GenerateFromPassword(data, cost)
	if err != nil {
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

func EncodeWithJwt(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tk, err := token.SignedString([]byte(*secret))
	if err != nil {
		log.Println(err)
	}
	return tk
}

func ParseFromJwToken(tk string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tk, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Sorry thats an invalid token sent")
		}
		return []byte(*secret), nil
	})
	if err != nil {
		return map[string]interface{}{}, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return map[string]interface{}{}, errors.New("Sorry we couldn't parse the token sent")
}

func generateNonce(key string, gcm cipher.AEAD) []byte {
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(strings.NewReader(key), nonce); err != nil {
		panic(err.Error())
	}
	return nonce
}
