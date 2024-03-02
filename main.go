package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// echo server

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request for -> ", r.URL.Path)

		fmt.Fprintf(w, "Hello from %s", r.URL.Path)
	})

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", mux)
}
