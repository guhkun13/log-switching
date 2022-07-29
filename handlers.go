package main

import (
	"encoding/json"
	"net/http"
)

func rootHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Halo"))
}

func inquiryHandler(rw http.ResponseWriter, r *http.Request) {
	db := Connect()
	inquirys, err := db.GetLatestInquiryRecords(r)
	panicOnErr(err)

	response, _ := json.Marshal(inquirys)
	rw.Write(response)
}
