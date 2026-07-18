package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alextse/go-key-value-store/internal/handler"
	"github.com/alextse/go-key-value-store/internal/store"
)

const (
	DefaultPort = ":8080"
)

func main() {
	kvStore := store.New()
	h := handler.New(kvStore)

	mux := http.NewServeMux()
	mux.HandleFunc("/store", h.Store)
	mux.HandleFunc("/retrieve", h.Retrieve)
	mux.HandleFunc("/update", h.Update)
	mux.HandleFunc("/delete", h.Delete)
	mux.HandleFunc("/health", h.Health)

	server := &http.Server{
		Addr:    DefaultPort,
		Handler: mux,
	}

	fmt.Printf("Starting key-value store server on %s\n", DefaultPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
