package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var urlMap = make(map[string]string) //an empty map is created

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano()) //Without this line, the output would be the same every time the program runs.
}

// generateShortLink generates a random string for the short link.
func generateShortLink(n int) string { //func function_name(Parameter-list)(Return_type){
	b := make([]byte, n) //SLICE  an empty slice is created with byte data type ,n is length of the slice
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	//everythime a new index number is generated and stored in b slice .
	return string(b)
	//When you use string(b), it converts each byte in the slice to its corresponding character in the resulting string.
} //uptill this inside everything is called function body, here is close of a function

// shortenURLHandler handles the URL shortening.
func shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL := generateShortLink(6) //shorturl="anNIge"
	urlMap[shortURL] = request.URL

	fmt.Printf("URL Map afetr adding: %+v\n", urlMap) //URL Map afetr adding: map[anNIge:http://example99.com]
	// %+v: This prints the value along with the field names (for structs)
	response := map[string]string{"short_url": "http://localhost:8080/" + shortURL}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// originalURLHandler handles returning the original URL.
func originalURLHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[len("/original/"):]

	originalURL, exists := urlMap[shortURL]

	fmt.Printf("URL Map during retrieval: %+v\n ", urlMap)      //URL Map during retrieval: map[anNIge:http://example99.com]
	fmt.Printf("Short URL: %s, Ecists: %v\n", shortURL, exists) //Short URL: anNIge, Ecists: true

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
