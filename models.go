package main

import (
	"database/sql"
	"time"
)

type NullString sql.NullString
type NullTime sql.NullTime

type Inquiry struct {
	Ts          time.Time  `json:"ts"`
	TsResp      NullTime   `json:"ts_resp"`
	RC          NullString `json:"rc"`
	Biller      string     `json:"biller"`
	IdByr       string     `json:"id_byr"`
	Nama        NullString `json:"nama"`
	Subbiller   NullString `json:"subbiller"`
	Bank        string     `json:"bank"`
	Channel     string     `json:"channel"`
	Terminal    string     `json:"terminal"`
	ElapsedTime NullString `json:"elapsedTime"`
}

type QueryFilter struct {
	Action     string
	Limit      string
	KodeBiller string
	Biller     string
	Subbiller  string
}

type HtmlCtx struct {
	Data    interface{}
	Filter  QueryFilter
	Error   error
	Status  bool
	Message string
}
