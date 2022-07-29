package main

import (
	"database/sql"
	"fmt"

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

type DB struct {
	db *sql.DB
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

func (d DB) GetLatestInquiryRecords() ([]Inquiry, error) {
	var inquirys []Inquiry

	sql := "select ts, ts_resp, rc, biller, subbiller, id_byr, nama from inquiry order by ts desc limit 10"
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

	return inquirys, nil
}
