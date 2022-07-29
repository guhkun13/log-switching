package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	chi := chi.NewRouter()

	chi.Use(middleware.Logger)

	chi.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Halo dunia"))
	})

	log.Fatal(http.ListenAndServe(":8000", chi))
}
