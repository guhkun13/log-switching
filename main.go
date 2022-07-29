package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"log"
	"net/http"

	"database/sql"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	dbname   = "switching"
	user     = "postgres"
	password = "postgres"
	sslmode  = "disable"
)

type Inquiry struct {
	Ts        time.Time  `json:"ts"`
	TsResp    NullTime   `json:"ts_resp"`
	RC        NullString `json:"rc"`
	Biller    string     `json:"biller"`
	IdByr     string     `json:"id_byr"`
	Nama      NullString `json:"nama"`
	Subbiller NullString `json:"subbiller"`
	// bank      string
	// channel   string
	// terminal  string
}

type NullString sql.NullString

func (ns *NullString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ns = NullString{s.String, false}
	} else {
		*ns = NullString{s.String, true}
	}

	return nil
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

type NullTime sql.NullTime

func (ns *NullTime) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ns = NullTime{s.Time, false}
	} else {
		*ns = NullTime{s.Time, true}
	}

	return nil
}

func (ns *NullTime) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.Time)
}

/** Handler Here **/

/** Main **/
func main() {
	fmt.Println("Init server")

	// init chi & middleware
	chi := chi.NewRouter()
	chi.Use(middleware.Logger)
	chi.Use(middleware.Recoverer)
	chi.Use(render.SetContentType(render.ContentTypeJSON))

	// start Database Connection
	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s", host, port, dbname, user, password, sslmode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("DB Connected Successfully")

	// start handler
	chi.Get("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Halo"))
	})

	chi.Get("/inquiry", func(rw http.ResponseWriter, r *http.Request) {
		// query
		sql := "select ts, ts_resp, rc, biller, subbiller, id_byr, nama from inquiry order by ts desc limit 10"
		rows, err := db.Query(sql)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// prepare
		var inquirys []Inquiry

		fmt.Println(rows)
		// loop through
		for rows.Next() {
			var inq Inquiry
			if err := rows.Scan(&inq.Ts, &inq.TsResp, &inq.RC, &inq.Biller, &inq.Subbiller, &inq.IdByr, &inq.Nama); err != nil {
				log.Fatal(err)
			}
			inquirys = append(inquirys, inq)
		}
		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}
		response, _ := json.Marshal(inquirys)
		rw.Write(response)
	})

	// start server
	port := ":8000"
	log.Fatal(http.ListenAndServe(port, chi))
}
