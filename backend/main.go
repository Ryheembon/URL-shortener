package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// URLMapping stores the mapping between short and long URLs
type URLMapping struct {
	LongURL string `json:"url"`
}

// Response represents the API response
type Response struct {
	ShortURL string `json:"shortUrl"`
}

// In-memory storage for URL mappings
var urlMap = make(map[string]string)

// CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from your React app
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	// Initialize random number generator
	rand.Seed(time.Now().UnixNano())

	//Handle URL shortening with CORS middleware
	http.HandleFunc("/shorten", corsMiddleware(handleShorten))

	//Handle URL redirection
	http.HandleFunc("/", handleRedirect)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid request method", http.StatusBadRequest)
		return
	}
	var mapping URLMapping
	if err := json.NewDecoder(r.Body).Decode(&mapping); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	//Generate a short URL
	ShortCode := generateShortCode()
	urlMap[ShortCode] = mapping.LongURL

	response := Response{
		ShortURL: "http://localhost:8080/" + ShortCode,
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/")
	if longURL, exists := urlMap[shortCode]; exists {
		http.Redirect(w, r, longURL, http.StatusMovedPermanently)
		return
	}
	http.Error(w, "URL not found", http.StatusNotFound)
}

func generateShortCode() string {
	// Generate a random 6-character string using letters and numbers
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6

	// Create a byte slice to store the result
	result := make([]byte, length)

	// Generate random bytes
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
