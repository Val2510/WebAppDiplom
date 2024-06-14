package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchBooks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		books := []Book{}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(books)
		if err != nil {
			t.Fatalf("Failed to encode response JSON: %v", err)
		}
	}))
	defer server.Close()

	query := "test"
	req, err := http.NewRequest("GET", server.URL+"?query="+query, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var responseBooks []Book
	err = json.NewDecoder(resp.Body).Decode(&responseBooks)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(responseBooks) != 0 {
		t.Fatalf("Expected empty list of books in response, got %d books", len(responseBooks))
	}
}
