package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

type Workshop struct {
	Name         string   `json:"name"`
	Date         string   `json:"date"`
	Presentator  string   `json:"presentator"`
	SweaterScore int      `json:"sweaterScore"`
	Participants []string `json:"participants"`
}

var workshop = Workshop{
	Name:         "ALM Workshop",
	Date:         "1/12/2025",
	Presentator:  "AE Consultants",
	SweaterScore: getDefaultSweaterScore(),
	Participants: []string{"John Doe", "Mary Little Lamb", "Chuck Norris", "Ting Lee"},
}

func getDefaultSweaterScore() int {
	if envScore := os.Getenv("DEFAULT_SWEATER_SCORE"); envScore != "" {
		if score, err := strconv.Atoi(envScore); err == nil && score >= 1 && score <= 10 {
			return score
		}
	}
	return 10
}

func getWorkshopHandler(w http.ResponseWriter, r *http.Request) {
	// Generate a random SweaterScore between 1 and 10
	workshop.SweaterScore = rand.Intn(10) + 1

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode the struct to JSON and write it to the response
	json.NewEncoder(w).Encode(workshop)
}

func postWorkshopHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming JSON data into a new Workshop struct
	var newWorkshop Workshop
	err := json.NewDecoder(r.Body).Decode(&newWorkshop)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON data"))
		return
	}

	// Validate SweaterScore is between 1 and 10
	if newWorkshop.SweaterScore < 1 || newWorkshop.SweaterScore > 10 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("SweaterScore must be between 1 and 10"))
		return
	}

	// Update the workshop details
	workshop = newWorkshop

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(workshop)
}

func WorkshopHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getWorkshopHandler(w, r)
	case "POST":
		postWorkshopHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}
