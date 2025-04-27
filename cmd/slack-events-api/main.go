package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	log.Print("Starting Slack Events API server")

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
