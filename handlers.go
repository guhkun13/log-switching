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

func logHandler(rw http.ResponseWriter, r *http.Request) {

	filter := buildFilter(r)
	fmt.Println(filter)

	// urus html
	funcMap := template.FuncMap{
		"getType": GetType,
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseGlob("./tmpl/*")
	panicOnErr(err)
	tmpl.ExecuteTemplate(rw, "header", nil)

	var ctx HtmlCtx
	ctx.Filter = filter
	ctx.Status = true

	errVal := filter.ValidateInput()
	if errVal != nil {
		ctx.Status = false
		ctx.Message = errVal.Error()
		tmpl.ExecuteTemplate(rw, "content", ctx)
		tmpl.ExecuteTemplate(rw, "footer", nil)
		return
	}

	var data []Inquiry
	if filter.Action == ACT_INQUIRY {
		data = apiGetInquiry(r)
	} else if filter.Action == ACT_PAYMENT {
		data = apiGetPayment(r)
	}
	ctx.Data = data
	tmpl.ExecuteTemplate(rw, "content", ctx)
	tmpl.ExecuteTemplate(rw, "footer", nil)
}

func GetType(subbiller string) string {
	if len(subbiller) > 0 {
		return "P2H"
	}
	return "H2H"
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
