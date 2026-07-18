package store

import (
	"fmt"
	"sync"

	"github.com/alextse/go-key-value-store/internal/crypto"
)

// Entry represents an encrypted data entry
type Entry struct {
	EncryptedData string // base64-encoded encrypted data
}

// Store is an in-memory key-value store with per-entry encryption
type Store struct {
	mu    sync.RWMutex
	data  map[string]Entry
	keys  map[string][]byte // Store the actual keys for validation
}

// New creates a new Store instance
func New() *Store {
	return &Store{
		data: make(map[string]Entry),
		keys: make(map[string][]byte),
	}
}

// StoreData encrypts and stores data with a unique encryption key
func (s *Store) StoreData(id string, data []byte) (string, error) {
	if id == "" {
		return "", fmt.Errorf("id cannot be empty")
	}

	// Generate a new encryption key for this entry
	key, err := crypto.GenerateKey()
	if err != nil {
		return "", fmt.Errorf("failed to generate encryption key: %w", err)
	}

	// Encrypt the data
	encrypted, err := crypto.Encrypt(data, key)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[id] = Entry{EncryptedData: encrypted}
	s.keys[id] = key

	return crypto.EncodeKey(key), nil
}

// RetrieveData decrypts and retrieves data using the provided encryption key
func (s *Store) RetrieveData(id string, keyStr string) ([]byte, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}
	if keyStr == "" {
		return nil, fmt.Errorf("encryption key cannot be empty")
	}

	key, err := crypto.DecodeKey(keyStr)
	if err != nil {
		return nil, fmt.Errorf("invalid encryption key: %w", err)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.data[id]
	if !exists {
		return nil, fmt.Errorf("entry not found for id: %s", id)
	}

	// Decrypt the data
	plaintext, err := crypto.Decrypt(entry.EncryptedData, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return plaintext, nil
}

// UpdateData updates existing data with the provided encryption key
func (s *Store) UpdateData(id string, newData []byte, keyStr string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}
	if keyStr == "" {
		return fmt.Errorf("encryption key cannot be empty")
	}

	key, err := crypto.DecodeKey(keyStr)
	if err != nil {
		return fmt.Errorf("invalid encryption key: %w", err)
	}

	// Encrypt the new data with the same key
	encrypted, err := crypto.Encrypt(newData, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Verify the entry exists and key matches
	if storedKey, exists := s.keys[id]; !exists {
		return fmt.Errorf("entry not found for id: %s", id)
	} else if !bytesEqual(storedKey, key) {
		return fmt.Errorf("invalid encryption key for id: %s", id)
	}

	s.data[id] = Entry{EncryptedData: encrypted}
	return nil
}

// DeleteData removes data if the encryption key matches
func (s *Store) DeleteData(id string, keyStr string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}
	if keyStr == "" {
		return fmt.Errorf("encryption key cannot be empty")
	}

	key, err := crypto.DecodeKey(keyStr)
	if err != nil {
		return fmt.Errorf("invalid encryption key: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Verify the entry exists and key matches
	if storedKey, exists := s.keys[id]; !exists {
		return fmt.Errorf("entry not found for id: %s", id)
	} else if !bytesEqual(storedKey, key) {
		return fmt.Errorf("invalid encryption key for id: %s", id)
	}

	delete(s.data, id)
	delete(s.keys, id)
	return nil
}

// bytesEqual compares two byte slices
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
