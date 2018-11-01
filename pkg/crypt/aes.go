package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/luckywinds/rshell/types"
	"io"
)

func AesEncrypt(plaintext string, cfg types.Cfg) (string, error) {
	if cfg.Passcrypttype != "aes" || cfg.Passcryptkey == "" || len(cfg.Passcryptkey) != 32 {
		return "", fmt.Errorf("Crypt failed. Passcrypttype must be %s and passcryptkey length must be %d.", "aes", 32)
	}

	block, err := aes.NewCipher([]byte(cfg.Passcryptkey))
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cipher.NewCFBEncrypter(block, iv).XORKeyStream(ciphertext[aes.BlockSize:],
		[]byte(plaintext))
	return hex.EncodeToString(ciphertext), nil

}

func AesDecrypt(d string, cfg types.Cfg) (string, error) {
	if d == "" {
		return "", nil
	}
	if cfg.Passcrypttype != "aes" || cfg.Passcryptkey == "" || len(cfg.Passcryptkey) != 32 {
		return "", fmt.Errorf("Crypt failed. Passcrypttype must be %s and passcryptkey length must be %d.", "aes", 32)
	}

	ciphertext, err := hex.DecodeString(d)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(cfg.Passcryptkey))
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("Decrypt failed, [%s].", "ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}
