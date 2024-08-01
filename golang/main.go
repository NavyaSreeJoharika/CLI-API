package main

import (
	"fmt"
	"net/http"
)

func world(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello")
}

func main() {
	http.HandleFunc("/world", world)

	fmt.Println("Server is starting on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}

/*
			package main

import (
	"net/http"
)

func main() {

http.HandleFunc("/Hello-World", func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World Printed"))
})

http.ListenAndServe(":8080", nil)
}

*/
