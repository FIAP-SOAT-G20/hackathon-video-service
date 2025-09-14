package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/FIAP-SOAT-G20/hackathon-video-service/internal/infrastructure/handler/request"
)

func main() {
	// Test that the UpdateVideoBodyRequest can handle Hash field
	jsonData := `{"status": "PROCESSING", "hash": "abc123hash456"}`
	var req request.UpdateVideoBodyRequest
	err := json.Unmarshal([]byte(jsonData), &req)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	fmt.Printf("Success! Status: %s, Hash: %s\n", req.Status, req.Hash)

	// Test UpdateVideoPartilBodyRequest too
	var partialReq request.UpdateVideoPartilBodyRequest
	err = json.Unmarshal([]byte(jsonData), &partialReq)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	fmt.Printf("Partial Success! Status: %s, Hash: %s\n", partialReq.Status, partialReq.Hash)
}
