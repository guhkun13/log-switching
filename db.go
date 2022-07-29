package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host       = "localhost"
	port       = 5432
	dbname     = "switching"
	user       = "postgres"
	password   = "postgres"
	sslmode    = "disable"
	limitQuery = "10"
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

func (d DB) Close() {
	defer d.db.Close()
}

func sqlFilterInquiry(r *http.Request, filter QueryFilter) string {
	sql := "select ts, ts_resp, AGE(ts_resp, ts) as elapsed_time, rc, biller, subbiller, id_byr, nama, bank, channel, terminal from inquiry"

	if filter.Biller != "" {
		sql += " where biller = '" + filter.Biller + "'"

		if filter.Subbiller != "" {
			sql += " and subbiller = '" + filter.Subbiller + "'"
		}
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
		err := rows.Scan(&inq.Ts, &inq.TsResp, &inq.ElapsedTime, &inq.RC, &inq.Biller, &inq.Subbiller, &inq.IdByr, &inq.Nama, &inq.Bank, &inq.Channel, &inq.Terminal)
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
