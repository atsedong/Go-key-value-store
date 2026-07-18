package crypto

import (
	"testing"
)

func TestGenerateKey(t *testing.T) {
	key1, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	if len(key1) != KeySize {
		t.Errorf("Expected key size %d, got %d", KeySize, len(key1))
	}

	// Generate another key and ensure they're different
	key2, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	if bytesEqual(key1, key2) {
		t.Error("Generated keys should be different")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	plaintext := []byte("Hello, World!")

	// Encrypt
	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if ciphertext == "" {
		t.Error("Encrypted data should not be empty")
	}

	// Decrypt
	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytesEqual(decrypted, plaintext) {
		t.Errorf("Decrypted data doesn't match original. Got %s, expected %s", string(decrypted), string(plaintext))
	}
}

func TestEncryptWithWrongKey(t *testing.T) {
	key1, _ := GenerateKey()
	key2, _ := GenerateKey()

	plaintext := []byte("Secret data")
	ciphertext, _ := Encrypt(plaintext, key1)

	_, err := Decrypt(ciphertext, key2)
	if err == nil {
		t.Error("Decryption with wrong key should fail")
	}
}

func TestEncodeDecodeKey(t *testing.T) {
	key, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey failed: %v", err)
	}

	// Encode
	encoded := EncodeKey(key)
	if encoded == "" {
		t.Error("Encoded key should not be empty")
	}

	// Decode
	decoded, err := DecodeKey(encoded)
	if err != nil {
		t.Fatalf("DecodeKey failed: %v", err)
	}

	if !bytesEqual(decoded, key) {
		t.Error("Decoded key doesn't match original")
	}
}

func TestInvalidKeySize(t *testing.T) {
	plaintext := []byte("test")
	invalidKey := []byte("short")

	_, err := Encrypt(plaintext, invalidKey)
	if err == nil {
		t.Error("Encrypt with invalid key size should fail")
	}

	_, err = Decrypt("ZmFrZQ==", invalidKey)
	if err == nil {
		t.Error("Decrypt with invalid key size should fail")
	}
}

func TestEmptyData(t *testing.T) {
	key, _ := GenerateKey()
	plaintext := []byte("")

	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt of empty data failed: %v", err)
	}

	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if len(decrypted) != 0 {
		t.Errorf("Expected empty decrypted data, got %v", decrypted)
	}
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
