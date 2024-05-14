package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
)

func md(k []byte) string {
	rs := md5.Sum(k)
	return hex.EncodeToString(rs[:])
}
func Cipher(keyPhrase []byte, value *[]byte) ([]byte, error) {
	gcm := gcmInstance(keyPhrase)
	nonce := make([]byte, gcm.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return []byte{}, nil
	}
	return gcm.Seal(nonce, nonce, *value, nil), nil
}
func gcmInstance(keyPhrase []byte) cipher.AEAD {
	hashedPhrase := md(keyPhrase)
	aesBlock, err := aes.NewCipher([]byte(hashedPhrase))
	if err != nil {
		log.Fatal(err)
	}
	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		log.Fatalln(err)
	}
	return gcm
}
func UnCipher(keyPhrase, ciphered []byte) (*[]byte, error) {
	gcm := gcmInstance(keyPhrase)
	nonceSize := gcm.NonceSize()
	nonce, cipheredText := ciphered[:nonceSize], ciphered[nonceSize:]
	data, err := gcm.Open(nil, nonce, cipheredText, nil)
	return &data, err
}
