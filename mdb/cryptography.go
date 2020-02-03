package mdb

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"syscall"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ssh/terminal"
)

func Encrypt(passphrase string, plaintext string) string {
	salt := make([]byte, 8)
	rand.Read(salt)
	key := deriveKey(passphrase, salt)
	iv := make([]byte, 12)
	rand.Read(iv)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	data := aesgcm.Seal(nil, iv, []byte(plaintext), nil)
	return fmt.Sprintf(
		"%s-%s-%s",
		hex.EncodeToString(salt),
		hex.EncodeToString(iv),
		hex.EncodeToString(data),
	)
}

func HashPassword(password string) string {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		panic(err)
	}
	return string(hashedPasswordBytes[:])
}

func CheckPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
	return err == nil
}

func ReadPassword() string {
	bytes, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func deriveKey(passphrase string, salt []byte) []byte {
	return pbkdf2.Key([]byte(passphrase), salt, 1000, 32, sha256.New)
}
