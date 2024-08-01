package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var (
	urls    = make(map[string]string)
	mu      sync.Mutex
	id      int
	baseURL = "http://localhost:8080/"
)

func main() {
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)
	fmt.Println("URL shortener service running on http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	id++
	shortID := fmt.Sprintf("%d", id)
	shortURL := baseURL + shortID
	urls[shortID] = longURL

	fmt.Fprintf(w, "Shortened URL: %s\n", shortURL)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortID := strings.TrimPrefix(r.URL.Path, "/")
	if shortID == "" {
		http.Error(w, "Short URL ID is required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	longURL, exists := urls[shortID]
	mu.Unlock()

	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}
