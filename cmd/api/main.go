package main

import (
	"log"
	"net/http"

	"github.com/Zerkina/url-shortener/internal/handlers"
)

func main() {
	handler := handlers.NewHandler()
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.MainPage)
	mux.HandleFunc("/{id}", handler.RedirectHandler)

	log.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
