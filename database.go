package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewConnect() (DB, error) {
	connect := "postgres://root@127.0.0.1:26257/santri_db?sslmode=disable"
	db, err := sql.Open("postgres", connect)
	if err != nil {
		log.Fatal("Cannot Connection", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB unreachable:", err)
	}
	return DB{db}, err
}
