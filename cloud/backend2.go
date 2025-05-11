package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from backend on port %s", r.URL.Path)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Starting server on port 9002...")
	http.ListenAndServe(":9002", nil)
}
