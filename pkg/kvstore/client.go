package kvstore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client is a client for the key-value store microservice
type Client struct {
	baseURL string
	client  *http.Client
}

// StoreResponse represents the response from a Store operation
type StoreResponse struct {
	ID            string `json:"id"`
	EncryptionKey string `json:"encryption_key"`
}

// RetrieveResponse represents the response from a Retrieve operation
type RetrieveResponse struct {
	Data string `json:"data"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewClient creates a new key-value store client
// baseURL should be something like "http://localhost:8080"
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

// Store stores data in the key-value store and returns the encryption key
func (c *Client) Store(id string, data []byte) (encryptionKey string, err error) {
	if id == "" {
		return "", fmt.Errorf("id cannot be empty")
	}

	payload := map[string]string{
		"id":   id,
		"data": string(data),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/store", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		var errResp ErrorResponse
		json.Unmarshal(body, &errResp)
		return "", fmt.Errorf("store failed: %s (status %d)", errResp.Error, resp.StatusCode)
	}

	var storeResp StoreResponse
	if err := json.Unmarshal(body, &storeResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return storeResp.EncryptionKey, nil
}

// Retrieve retrieves data from the key-value store
func (c *Client) Retrieve(id string, encryptionKey string) ([]byte, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}
	if encryptionKey == "" {
		return nil, fmt.Errorf("encryption key cannot be empty")
	}

	query := url.Values{}
	query.Set("id", id)
	query.Set("encryption_key", encryptionKey)

	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/retrieve?"+query.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.Unmarshal(body, &errResp)
		return nil, fmt.Errorf("retrieve failed: %s (status %d)", errResp.Error, resp.StatusCode)
	}

	var retrieveResp RetrieveResponse
	if err := json.Unmarshal(body, &retrieveResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return []byte(retrieveResp.Data), nil
}

// Update updates data in the key-value store
func (c *Client) Update(id string, data []byte, encryptionKey string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}
	if encryptionKey == "" {
		return fmt.Errorf("encryption key cannot be empty")
	}

	payload := map[string]string{
		"id":               id,
		"data":             string(data),
		"encryption_key":   encryptionKey,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, c.baseURL+"/update", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.Unmarshal(body, &errResp)
		return fmt.Errorf("update failed: %s (status %d)", errResp.Error, resp.StatusCode)
	}

	return nil
}

// Delete deletes data from the key-value store
func (c *Client) Delete(id string, encryptionKey string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}
	if encryptionKey == "" {
		return fmt.Errorf("encryption key cannot be empty")
	}

	query := url.Values{}
	query.Set("id", id)
	query.Set("encryption_key", encryptionKey)

	req, err := http.NewRequest(http.MethodDelete, c.baseURL+"/delete?"+query.Encode(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.Unmarshal(body, &errResp)
		return fmt.Errorf("delete failed: %s (status %d)", errResp.Error, resp.StatusCode)
	}

	return nil
}
