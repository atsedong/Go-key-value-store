package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/alextse/go-key-value-store/internal/store"
)

type Handler struct {
	store *store.Store
}

type StoreRequest struct {
	Data string `json:"data"`
	ID   string `json:"id"`
}

type StoreResponse struct {
	ID            string `json:"id"`
	EncryptionKey string `json:"encryption_key"`
}

type RetrieveRequest struct {
	ID            string `json:"id"`
	EncryptionKey string `json:"encryption_key"`
}

type RetrieveResponse struct {
	Data string `json:"data"`
}

type UpdateRequest struct {
	ID            string `json:"id"`
	Data          string `json:"data"`
	EncryptionKey string `json:"encryption_key"`
}

type DeleteRequest struct {
	ID            string `json:"id"`
	EncryptionKey string `json:"encryption_key"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// New creates a new Handler with the given store
func New(s *store.Store) *Handler {
	return &Handler{store: s}
}

// Store handles POST /store
func (h *Handler) Store(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "method not allowed"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to read body"})
		return
	}
	defer r.Body.Close()

	var req StoreRequest
	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid JSON"})
		return
	}

	if req.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "id is required"})
		return
	}

	if req.Data == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "data is required"})
		return
	}

	key, err := h.store.StoreData(req.ID, []byte(req.Data))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(StoreResponse{
		ID:            req.ID,
		EncryptionKey: key,
	})
}

// Retrieve handles GET /retrieve
func (h *Handler) Retrieve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "method not allowed"})
		return
	}

	id := r.URL.Query().Get("id")
	key := r.URL.Query().Get("encryption_key")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "id parameter is required"})
		return
	}

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "encryption_key parameter is required"})
		return
	}

	data, err := h.store.RetrieveData(id, key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RetrieveResponse{Data: string(data)})
}

// Update handles PUT /update
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "method not allowed"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "failed to read body"})
		return
	}
	defer r.Body.Close()

	var req UpdateRequest
	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "invalid JSON"})
		return
	}

	if req.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "id is required"})
		return
	}

	if req.EncryptionKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "encryption_key is required"})
		return
	}

	if req.Data == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "data is required"})
		return
	}

	if err := h.store.UpdateData(req.ID, []byte(req.Data), req.EncryptionKey); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{})
}

// Delete handles DELETE /delete
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "method not allowed"})
		return
	}

	id := r.URL.Query().Get("id")
	key := r.URL.Query().Get("encryption_key")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "id parameter is required"})
		return
	}

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "encryption_key parameter is required"})
		return
	}

	if err := h.store.DeleteData(id, key); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{})
}

// Health handles GET /health
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
