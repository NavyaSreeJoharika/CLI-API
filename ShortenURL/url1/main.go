package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var urlMap = struct {
	sync.RWMutex
	mapping map[string]string
}{mapping: make(map[string]string)}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// generateShortLink generates a random string for the short link.
func generateShortLink(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// shortenURLHandler handles the URL shortening.
func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL := generateShortLink(6)
	urlMap.Lock()
	urlMap.mapping[shortURL] = request.URL
	urlMap.Unlock()

	response := map[string]string{"short_url": "http://localhost:8080/" + shortURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// redirectHandler handles the redirection to the original URL.
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[len("/redirect/"):]

	urlMap.RLock()
	originalURL, exists := urlMap.mapping[shortURL]
	urlMap.RUnlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"original_url": originalURL})
}

func main() {
	http.HandleFunc("/shorten", shortenURLHandler)
	http.HandleFunc("/redirect/", redirectHandler) // Modified the endpoint to include "/redirect/"
	http.ListenAndServe(":8080", nil)
	fmt.Println("URL shortener service running on http://localhost:8080/")

}
