package main

import (
	"database/sql"
	"fmt"
	"net/http"

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

const (
	limitQuery = "10"
)

type DB struct {
	db *sql.DB
}

type QueryFilter struct {
	Limit      string
	KodeBiller string
	Biller     string
	Subbiller  string
}

func Connect() *DB {
	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s", host, port, dbname, user, password, sslmode)
	conn, err := sql.Open("postgres", dsn)
	panicOnErr(err)

	err = conn.Ping()
	panicOnErr(err)

	fmt.Println("DB Connected Successfully")
	return &DB{db: conn}
}

func (d DB) Close() {
	defer d.db.Close()
}

func buildFilter(r *http.Request) QueryFilter {
	var filter QueryFilter
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

	return filter
}

func sqlFilterInquiry(r *http.Request, filter QueryFilter) string {
	sql := "select ts, ts_resp, rc, biller, subbiller, id_byr, nama from inquiry"

	sql += " where biller = '" + filter.Biller + "'"
	if filter.Subbiller != "" {
		sql += " and subbiller = '" + filter.Subbiller + "'"
	}
	sql += " order by ts desc limit " + filter.Limit
	fmt.Println(sql)

	return sql
}

func (d DB) GetLatestInquiryRecords(r *http.Request) ([]Inquiry, error) {
	var inquirys []Inquiry

	filter := buildFilter(r)
	sql := sqlFilterInquiry(r, filter)

	rows, err := d.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// loop through
	for rows.Next() {
		var inq Inquiry
		err := rows.Scan(&inq.Ts, &inq.TsResp, &inq.RC, &inq.Biller, &inq.Subbiller, &inq.IdByr, &inq.Nama)
		if err != nil {
			return nil, err
		}
		inquirys = append(inquirys, inq)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	d.Close()
	return inquirys, nil
}
