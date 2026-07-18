package store

import (
	"testing"
)

func TestStoreRetrieve(t *testing.T) {
	s := New()
	id := "test-id"
	data := []byte("test data")

	// Store data
	key, err := s.StoreData(id, data)
	if err != nil {
		t.Fatalf("StoreData failed: %v", err)
	}

	if key == "" {
		t.Error("Encryption key should not be empty")
	}

	// Retrieve data
	retrieved, err := s.RetrieveData(id, key)
	if err != nil {
		t.Fatalf("RetrieveData failed: %v", err)
	}

	if string(retrieved) != string(data) {
		t.Errorf("Retrieved data doesn't match. Got %s, expected %s", string(retrieved), string(data))
	}
}

func TestRetrieveWithWrongKey(t *testing.T) {
	s := New()
	id := "test-id"
	data := []byte("test data")

	_, _ = s.StoreData(id, data)

	_, err := s.RetrieveData(id, "wrong-key")
	if err == nil {
		t.Error("RetrieveData with wrong key should fail")
	}
}

func TestRetrieveNonExistent(t *testing.T) {
	s := New()

	_, err := s.RetrieveData("nonexistent", "some-key")
	if err == nil {
		t.Error("RetrieveData for nonexistent entry should fail")
	}
}

func TestUpdate(t *testing.T) {
	s := New()
	id := "test-id"
	originalData := []byte("original")
	newData := []byte("updated")

	key, _ := s.StoreData(id, originalData)

	// Update
	err := s.UpdateData(id, newData, key)
	if err != nil {
		t.Fatalf("UpdateData failed: %v", err)
	}

	// Verify update
	retrieved, _ := s.RetrieveData(id, key)
	if string(retrieved) != string(newData) {
		t.Errorf("Updated data doesn't match. Got %s, expected %s", string(retrieved), string(newData))
	}
}

func TestUpdateWithWrongKey(t *testing.T) {
	s := New()
	id := "test-id"
	data := []byte("data")

	s.StoreData(id, data)

	err := s.UpdateData(id, []byte("new"), "wrong-key")
	if err == nil {
		t.Error("UpdateData with wrong key should fail")
	}
}

func TestUpdateNonExistent(t *testing.T) {
	s := New()

	err := s.UpdateData("nonexistent", []byte("data"), "some-key")
	if err == nil {
		t.Error("UpdateData for nonexistent entry should fail")
	}
}

func TestDelete(t *testing.T) {
	s := New()
	id := "test-id"
	data := []byte("test data")

	key, _ := s.StoreData(id, data)

	// Delete
	err := s.DeleteData(id, key)
	if err != nil {
		t.Fatalf("DeleteData failed: %v", err)
	}

	// Verify deletion
	_, err = s.RetrieveData(id, key)
	if err == nil {
		t.Error("RetrieveData after delete should fail")
	}
}

func TestDeleteWithWrongKey(t *testing.T) {
	s := New()
	id := "test-id"

	s.StoreData(id, []byte("data"))

	err := s.DeleteData(id, "wrong-key")
	if err == nil {
		t.Error("DeleteData with wrong key should fail")
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := New()

	err := s.DeleteData("nonexistent", "some-key")
	if err == nil {
		t.Error("DeleteData for nonexistent entry should fail")
	}
}

func TestEmptyID(t *testing.T) {
	s := New()

	_, err := s.StoreData("", []byte("data"))
	if err == nil {
		t.Error("StoreData with empty ID should fail")
	}

	_, err = s.RetrieveData("", "key")
	if err == nil {
		t.Error("RetrieveData with empty ID should fail")
	}

	err = s.UpdateData("", []byte("data"), "key")
	if err == nil {
		t.Error("UpdateData with empty ID should fail")
	}

	err = s.DeleteData("", "key")
	if err == nil {
		t.Error("DeleteData with empty ID should fail")
	}
}

func TestMultipleEntries(t *testing.T) {
	s := New()

	entries := map[string][]byte{
		"id-1": []byte("data-1"),
		"id-2": []byte("data-2"),
		"id-3": []byte("data-3"),
	}

	keys := make(map[string]string)

	// Store multiple entries
	for id, data := range entries {
		key, err := s.StoreData(id, data)
		if err != nil {
			t.Fatalf("StoreData failed for %s: %v", id, err)
		}
		keys[id] = key
	}

	// Retrieve and verify each entry
	for id, expectedData := range entries {
		retrieved, err := s.RetrieveData(id, keys[id])
		if err != nil {
			t.Fatalf("RetrieveData failed for %s: %v", id, err)
		}

		if string(retrieved) != string(expectedData) {
			t.Errorf("Data mismatch for %s. Got %s, expected %s", id, string(retrieved), string(expectedData))
		}
	}
}

func TestConcurrentOperations(t *testing.T) {
	s := New()
	done := make(chan bool)

	// Concurrent writes
	go func() {
		s.StoreData("id-1", []byte("data-1"))
		done <- true
	}()

	go func() {
		s.StoreData("id-2", []byte("data-2"))
		done <- true
	}()

	<-done
	<-done

	// Concurrent reads
	key1, _ := s.StoreData("id-read", []byte("test"))
	results := make(chan []byte)

	go func() {
		data, _ := s.RetrieveData("id-read", key1)
		results <- data
	}()

	go func() {
		data, _ := s.RetrieveData("id-read", key1)
		results <- data
	}()

	<-results
	<-results
}
