package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func rootHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Halo"))
}

type HtmlCtx struct {
	Data    interface{}
	Filter  QueryFilter
	Error   error
	Status  bool
	Message string
}

const (
	ACT_INQUIRY = "inquiry"
	ACT_PAYMENT = "payment"
)

var tmpl *template.Template

func logHandler(rw http.ResponseWriter, r *http.Request) {
	var err error
	tmpl, err = template.ParseGlob("./tmpl/*")
	if err != nil {
		panicOnErr(err)
	}

	// tmpl, err := template.New("index.html").Funcs(
	// 	template.FuncMap{
	// 		"getBillerType": func(biller string, subbiller string) string {
	// 			if subbiller != "" {
	// 				return "P2H"
	// 			}
	// 			return "H2H"
	// 		},
	// 	},
	// ).ParseFiles("tmpl/index.html")
	// panicOnErr(err)

	var ctx HtmlCtx

	filter := buildFilter(r)
	fmt.Println(filter)
	ctx.Filter = filter
	ctx.Status = true

	var data []Inquiry
	if filter.Action == ACT_INQUIRY {
		data = apiGetInquiry(r)
	} else if filter.Action == ACT_PAYMENT {
		data = apiGetPayment(r)
	} else {
		fmt.Println("ERROR invalid action")
		data = nil
		ctx.Status = false
		ctx.Message = "Action not valid"
	}
	ctx.Data = data

	tmpl.ExecuteTemplate(rw, "base.html", ctx)
}

func inquiryHandler(rw http.ResponseWriter, r *http.Request) {
	data := apiGetInquiry(r)
	response, _ := json.Marshal(data)
	rw.Write(response)
}

func apiGetInquiry(r *http.Request) []Inquiry {
	db := Connect()
	data, err := db.GetLatestInquiryRecords(r)
	panicOnErr(err)
	db.Close()
	return data
}

func apiGetPayment(r *http.Request) []Inquiry {
	db := Connect()
	data, err := db.GetLatestInquiryRecords(r)
	panicOnErr(err)
	db.Close()
	return data
}
