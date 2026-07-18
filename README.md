# Go Key-Value Store Microservice

A secure, encrypted key-value store microservice written in Go. Each entry is encrypted with its own unique AES-256-GCM encryption key, providing strong encryption guarantees.

## Features

- **Per-Entry Encryption**: Each stored value is encrypted with a unique 256-bit AES key
- **Secure Encryption**: Uses AES-256-GCM (Galois/Counter Mode) for authenticated encryption
- **Thread-Safe**: Internal store uses RWMutex for safe concurrent access
- **Simple HTTP API**: RESTful API for all operations
- **Client Library**: Go client package for easy integration into other microservices
- **Easy to Run**: No external dependencies, just standard Go libraries

## Architecture

### Components

- **cmd/server/**: HTTP server implementation
- **internal/store/**: In-memory encrypted key-value store with thread safety
- **internal/crypto/**: Encryption/decryption utilities (AES-256-GCM)
- **internal/handler/**: HTTP request handlers
- **pkg/kvstore/**: Client library for calling the store from other services
- **example/**: Example client usage

## API Endpoints

### Store
**POST /store**

Stores data and returns a unique encryption key.

Request:
```json
{
  "id": "user-123",
  "data": "sensitive information"
}
```

Response (201 Created):
```json
{
  "id": "user-123",
  "encryption_key": "base64-encoded-256-bit-key"
}
```

### Retrieve
**GET /retrieve?id=user-123&encryption_key=base64-encoded-key**

Retrieves and decrypts data.

Response (200 OK):
```json
{
  "data": "sensitive information"
}
```

Response (404 Not Found):
```json
{
  "error": "entry not found for id: user-123"
}
```

### Update
**PUT /update**

Updates existing data (must provide correct encryption key).

Request:
```json
{
  "id": "user-123",
  "data": "new data",
  "encryption_key": "base64-encoded-key"
}
```

Response (200 OK):
```json
{}
```

### Delete
**DELETE /delete?id=user-123&encryption_key=base64-encoded-key**

Deletes data from the store (requires correct encryption key).

Response (200 OK):
```json
{}
```

### Health
**GET /health**

Health check endpoint.

Response (200 OK):
```json
{
  "status": "ok"
}
```

## Building and Running

### Prerequisites
- Go 1.21 or later

### Build and Run Server

```bash
# Build the server
go build -o server ./cmd/server

# Run the server (listens on :8080 by default)
./server
```

Or run directly:
```bash
go run ./cmd/server/main.go
```

### Run Example Client

In another terminal:

```bash
# Run the example (make sure server is running on localhost:8080)
go run ./example/main.go
```

Expected output:
```
=== Key-Value Store Client Example ===

1. Storing data...
   Stored data with ID: user-123
   Encryption key: base64-encoded-key

2. Retrieving data...
   Retrieved data: Hello, Secret World!
   Matches original: true

3. Updating data...
   Data updated successfully
   New data: Updated Secret World!

4. Testing with wrong encryption key...
   Expected error: retrieve failed: failed to decrypt data: ...

5. Deleting data...
   Data deleted successfully
   Expected error when retrieving deleted data: retrieve failed: entry not found for id: user-123

6. Storing multiple entries...
   Stored entry-1 with key: ...
   Stored entry-2 with key: ...
   Stored entry-3 with key: ...

7. Retrieving all entries...
   entry-1: First entry
   entry-2: Second entry
   entry-3: Third entry

=== All tests completed successfully ===
```

### Running with Docker

#### Prerequisites
- Docker installed

#### Build and Run with Docker

```bash
# Build the Docker image
docker build -t go-kvstore .

# Run the container
docker run -p 8080:8080 go-kvstore
```

#### Docker Usage Examples

Test the running service:
```bash
# Health check
curl http://localhost:8080/health

# Store data
curl -X POST http://localhost:8080/store \
  -H "Content-Type: application/json" \
  -d '{"id":"docker-test","data":"Hello from Docker"}'

# Retrieve data
curl "http://localhost:8080/retrieve?id=docker-test&encryption_key=YOUR_KEY_HERE"
```
