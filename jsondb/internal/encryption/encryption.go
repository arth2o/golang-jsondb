package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
)

var (
	ErrInvalidKeyLength = errors.New("encryption key must be 32 bytes long")
	ErrEncryption      = errors.New("encryption failed")
	ErrDecryption      = errors.New("decryption failed")
)

type Encryptor struct {
	block cipher.Block
}

func NewEncryptor(key string) (*Encryptor, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("encryption key must be exactly 32 bytes, got %d", len(key))
	}
	
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}
	
	return &Encryptor{block: block}, nil
}

func (e *Encryptor) Encrypt(data []byte) ([]byte, error) {
	log.Printf("Encrypting data of length: %d", len(data))
	
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %v", err)
	}

	stream := cipher.NewCTR(e.block, nonce)
	encrypted := make([]byte, len(data))
	stream.XORKeyStream(encrypted, data)

	// Combine nonce and encrypted data
	result := make([]byte, len(nonce)+len(encrypted))
	copy(result, nonce)
	copy(result[len(nonce):], encrypted)

	log.Printf("Encrypted data length: %d", len(result))
	
	return result, nil
}

func (e *Encryptor) Decrypt(data []byte) ([]byte, error) {
	log.Printf("Decrypting data of length: %d", len(data))
	
	if len(data) < 16 {
		return nil, fmt.Errorf("encrypted data too short: %d bytes", len(data))
	}

	nonce := data[:16]
	ciphertext := data[16:]
	
	stream := cipher.NewCTR(e.block, nonce)
	decrypted := make([]byte, len(ciphertext))
	stream.XORKeyStream(decrypted, ciphertext)
	
	log.Printf("Decrypted data length: %d", len(decrypted))
	
	return decrypted, nil
}