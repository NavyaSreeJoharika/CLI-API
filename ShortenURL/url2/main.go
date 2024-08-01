package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var urlMap = make(map[string]string)

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
	urlMap[shortURL] = request.URL

	fmt.Printf("URL Map afetr adding: %+v\n", urlMap)

	response := map[string]string{"short_url": "http://localhost:8080/" + shortURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// originalURLHandler handles returning the original URL.
func originalURLHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[len("/original/"):]

	originalURL, exists := urlMap[shortURL]

	fmt.Printf("URL Map during retrieval: %+v\n ", urlMap)
	fmt.Printf("Short URL: %s, Ecists: %v\n", shortURL, exists)

	if !exists {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"original_url": originalURL})
}

func main() {
	http.HandleFunc("/shorten", shortenURLHandler)
	http.HandleFunc("/original/", originalURLHandler)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Server is starting on 8080 port")
}
