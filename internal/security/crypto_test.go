package security

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	secretData := []byte("sk-proj-super-secret-api-key-12345")

	encrypted, err := Encrypt(secretData, key)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	decrypted, err := Decrypt(encrypted.Ciphertext, encrypted.Nonce, key)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if !bytes.Equal(secretData, decrypted) {
		t.Errorf("expected %s, got %s", string(secretData), string(decrypted))
	}
}
