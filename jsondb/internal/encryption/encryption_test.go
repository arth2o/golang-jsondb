package encryption

import (
	"bytes"
	"testing"
)

func TestNewEncryptor(t *testing.T) {
    tests := []struct {
        name    string
        key     string
        wantErr bool
    }{
        {
            name:    "Valid Key",
            key:     "0123456789abcdef0123456789abcdef",
            wantErr: false,
        },
        {
            name:    "Invalid Key Length",
            key:     "too_short",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            enc, err := NewEncryptor(tt.key)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewEncryptor() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && enc == nil {
                t.Error("NewEncryptor() returned nil encryptor with no error")
            }
        })
    }
}

func TestEncryptDecrypt(t *testing.T) {
    key := "0123456789abcdef0123456789abcdef"
    enc, err := NewEncryptor(key)
    if err != nil {
        t.Fatalf("Failed to create encryptor: %v", err)
    }

    original := []byte("test data")
    encrypted, err := enc.Encrypt(original)
    if err != nil {
        t.Fatalf("Encryption failed: %v", err)
    }

    decrypted, err := enc.Decrypt(encrypted)
    if err != nil {
        t.Fatalf("Decryption failed: %v", err)
    }

    if !bytes.Equal(original, decrypted) {
        t.Errorf("Decrypted data doesn't match original: got %q, want %q", decrypted, original)
    }
}

func TestDecryptInvalidData(t *testing.T) {
    key := "0123456789abcdef0123456789abcdef"
    enc, err := NewEncryptor(key)
    if err != nil {
        t.Fatalf("Failed to create encryptor: %v", err)
    }

    tests := []struct {
        name    string
        data    []byte
        wantErr bool
    }{
        {
            name:    "Too Short",
            data:    []byte("short"),
            wantErr: true,
        },
        {
            name:    "Empty Data",
            data:    []byte{},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := enc.Decrypt(tt.data)
            if (err != nil) != tt.wantErr {
                t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}