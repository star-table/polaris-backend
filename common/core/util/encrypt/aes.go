package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

func AesDecrypt(key string, encrypt string) (string, error) {
	kbs := SHA256(key)
	decode, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", err
	}
	if len(decode) < aes.BlockSize {
		return "", errors.New("密文太短啦")
	}
	iv := decode[:aes.BlockSize]
	block, err := aes.NewCipher(kbs)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	plantText := make([]byte, len(decode))
	blockMode.CryptBlocks(plantText, decode)
	plantText = PKCS7UnPadding(plantText)
	plantText = plantText[aes.BlockSize:]
	return string(plantText), nil
}

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

func SHA256(source string) []byte {
	mac := sha256.New()
	mac.Write([]byte(source))
	return mac.Sum(nil)
}

func AesEncrypt(plainText string, key string) (string, error) {
	bKey := sha256.Sum256([]byte(key))
	bPlaintext := pKCS7Padding([]byte(plainText))
	if len(bPlaintext)%aes.BlockSize != 0 {
		return "", errors.New("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(bKey[:])
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(bPlaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], bPlaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func pKCS7Padding(cipherText []byte) []byte {
	padding := aes.BlockSize - len(cipherText)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func SHA1(source string) string {
	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(source))
	return hex.EncodeToString(sha1Hash.Sum([]byte(nil)))
}
