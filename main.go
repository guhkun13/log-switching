package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"log"
	"net/http"
)

/** Handler Here **/

func rootHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Halo"))
}

func inquiryHandler(rw http.ResponseWriter, r *http.Request) {
	db := Connect()
	inquirys, err := db.GetLatestInquiryRecords()
	panicOnErr(err)

	response, _ := json.Marshal(inquirys)
	rw.Write(response)
}

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
