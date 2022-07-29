package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"reflect"
)

func buildFilter(r *http.Request) QueryFilter {
	var filter QueryFilter
	action := r.URL.Query().Get("action")
	limit := r.URL.Query().Get("limit")
	kodeBiller := r.URL.Query().Get("biller")

	if limit == "" {
		limit = limitQuery
	}

	if len(kodeBiller) == 8 {
		filter.Biller = kodeBiller[:4]
		filter.Subbiller = kodeBiller[4:]
	}
	filter.Limit = limit
	filter.Action = action

	return filter
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
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
