package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"golang.org/x/exp/slices"
	"net/http"
	"reflect"
	"strconv"
)

const (
	ACT_INQUIRY = "inquiry"
	ACT_PAYMENT = "payment"
)

func GetValidActions() []string {
	result := make([]string, 2)
	result = append(result, ACT_INQUIRY)
	result = append(result, ACT_PAYMENT)

	return result
}

func buildFilter(r *http.Request) QueryFilter {
	var filter QueryFilter
	action := r.URL.Query().Get("action")
	limit := r.URL.Query().Get("limit")
	kodeBiller := r.URL.Query().Get("kodeBiller")

	if limit == "" {
		limit = limitQuery
	}
	if kodeBiller != "" {
		if len(kodeBiller) == 8 {
			filter.Biller = kodeBiller[:4]
			filter.Subbiller = kodeBiller[4:]
		} else {
			filter.Biller = kodeBiller
		}
	}
	filter.Limit = limit
	filter.Action = action
	filter.KodeBiller = kodeBiller

	return filter
}

func validateFilter(f QueryFilter) error {
	// validate kodeBiller
	lenKodeBiller := len(f.KodeBiller)
	if f.KodeBiller != "" && (lenKodeBiller != 4 || lenKodeBiller == 8) {
		return errors.New("invalid biller length. Must be 4 or 8")
	}
	if _, err := strconv.Atoi(f.KodeBiller); err != nil {
		return errors.New("invalid biller format. Must be int")
	}

	// validate action
	validActions := GetValidActions()
	if exist := slices.Contains(validActions, f.Action); !exist {
		return errors.New("invalid action. Must be INQUIRY or PAYMENT")
	}

	return nil
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

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
