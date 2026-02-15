package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetDefaultSweaterScore(t *testing.T) {
	// Test default value when env var is not set
	os.Unsetenv("DEFAULT_SWEATER_SCORE")
	score := getDefaultSweaterScore()
	if score != 10 {
		t.Errorf("Expected default score 10, got %d", score)
	}

	// Test valid env var value
	os.Setenv("DEFAULT_SWEATER_SCORE", "7")
	score = getDefaultSweaterScore()
	if score != 7 {
		t.Errorf("Expected score 7, got %d", score)
	}

	// Test invalid env var value (out of range - too high)
	os.Setenv("DEFAULT_SWEATER_SCORE", "15")
	score = getDefaultSweaterScore()
	if score != 10 {
		t.Errorf("Expected default score 10 for out of range value, got %d", score)
	}

	// Test invalid env var value (out of range - too low)
	os.Setenv("DEFAULT_SWEATER_SCORE", "0")
	score = getDefaultSweaterScore()
	if score != 10 {
		t.Errorf("Expected default score 10 for out of range value, got %d", score)
	}

	// Test invalid env var value (not a number)
	os.Setenv("DEFAULT_SWEATER_SCORE", "invalid")
	score = getDefaultSweaterScore()
	if score != 10 {
		t.Errorf("Expected default score 10 for invalid value, got %d", score)
	}

	// Cleanup
	os.Unsetenv("DEFAULT_SWEATER_SCORE")
}

func TestGetWorkshopHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/workshop", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getWorkshopHandler)
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check content type
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Handler returned wrong content type: got %v want %v", contentType, "application/json")
	}

	// Check response body
	var response Workshop
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	// Check SweaterScore is within valid range
	if response.SweaterScore < 1 || response.SweaterScore > 10 {
		t.Errorf("SweaterScore should be between 1 and 10, got %d", response.SweaterScore)
	}
}

func TestPostWorkshopHandler_ValidData(t *testing.T) {
	newWorkshop := Workshop{
		Name:         "New Workshop",
		Date:         "2/15/2026",
		Presentator:  "Test Presenter",
		SweaterScore: 8,
		Participants: []string{"Alice", "Bob"},
	}

	body, _ := json.Marshal(newWorkshop)
	req, err := http.NewRequest("POST", "/workshop", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postWorkshopHandler)
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check response body
	var response Workshop
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	if response.Name != newWorkshop.Name {
		t.Errorf("Expected name %s, got %s", newWorkshop.Name, response.Name)
	}

	if response.SweaterScore != newWorkshop.SweaterScore {
		t.Errorf("Expected SweaterScore %d, got %d", newWorkshop.SweaterScore, response.SweaterScore)
	}
}

func TestPostWorkshopHandler_InvalidJSON(t *testing.T) {
	req, err := http.NewRequest("POST", "/workshop", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postWorkshopHandler)
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Check error message
	if rr.Body.String() != "Invalid JSON data" {
		t.Errorf("Expected error message 'Invalid JSON data', got '%s'", rr.Body.String())
	}
}

func TestPostWorkshopHandler_SweaterScoreTooLow(t *testing.T) {
	newWorkshop := Workshop{
		Name:         "Test Workshop",
		Date:         "2/15/2026",
		Presentator:  "Test Presenter",
		SweaterScore: 0,
		Participants: []string{"Alice"},
	}

	body, _ := json.Marshal(newWorkshop)
	req, err := http.NewRequest("POST", "/workshop", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postWorkshopHandler)
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Check error message
	if rr.Body.String() != "SweaterScore must be between 1 and 10" {
		t.Errorf("Expected validation error message, got '%s'", rr.Body.String())
	}
}

func TestPostWorkshopHandler_SweaterScoreTooHigh(t *testing.T) {
	newWorkshop := Workshop{
		Name:         "Test Workshop",
		Date:         "2/15/2026",
		Presentator:  "Test Presenter",
		SweaterScore: 11,
		Participants: []string{"Alice"},
	}

	body, _ := json.Marshal(newWorkshop)
	req, err := http.NewRequest("POST", "/workshop", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postWorkshopHandler)
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Check error message
	if rr.Body.String() != "SweaterScore must be between 1 and 10" {
		t.Errorf("Expected validation error message, got '%s'", rr.Body.String())
	}
}

func TestWorkshopHandler_GET(t *testing.T) {
	req, err := http.NewRequest("GET", "/workshop", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(WorkshopHandler)
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestWorkshopHandler_POST(t *testing.T) {
	newWorkshop := Workshop{
		Name:         "POST Test Workshop",
		Date:         "2/15/2026",
		Presentator:  "Test Presenter",
		SweaterScore: 5,
		Participants: []string{"Alice"},
	}

	body, _ := json.Marshal(newWorkshop)
	req, err := http.NewRequest("POST", "/workshop", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(WorkshopHandler)
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestWorkshopHandler_MethodNotAllowed(t *testing.T) {
	methods := []string{"PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		req, err := http.NewRequest(method, "/workshop", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(WorkshopHandler)
		handler.ServeHTTP(rr, req)

		// Check status code
		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("Handler returned wrong status code for %s: got %v want %v", method, status, http.StatusMethodNotAllowed)
		}

		// Check error message
		if rr.Body.String() != "Method not allowed" {
			t.Errorf("Expected error message 'Method not allowed', got '%s'", rr.Body.String())
		}
	}
}
