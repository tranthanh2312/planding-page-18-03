package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Querier struct {
	DB *sql.DB
}

const (
	PASSWORD = "123456"
)

func ConnectDB() (db *sql.DB) {
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/defaultdb", PASSWORD))
	if err != nil {
		panic(err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Cannot ping to database server:", err.Error())
	}
	return db
}

func NewQuerier() *Querier {
	return &Querier{
		DB: ConnectDB(),
	}
}
