package main

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"log"
	"net/http"
)

/** Main **/
func main() {
	fmt.Println("Init server")

	// init chi & middleware
	chi := chi.NewRouter()
	chi.Use(middleware.Logger)
	chi.Use(middleware.Recoverer)
	chi.Use(render.SetContentType(render.ContentTypeJSON))

	// start handler
	chi.Get("/", rootHandler)
	chi.Get("/inquiry", inquiryHandler)

	// start server
	port := ":8000"
	log.Fatal(http.ListenAndServe(port, chi))
}
