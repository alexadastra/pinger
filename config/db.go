package config

import (
	"database/sql"
	"fmt"
)

const (
	host     = "url-microservice_postgres_1"
	port     = 5432
	user     = "root"
	password = "password"
	dbname   = "microservice-db"
)

var DB *sql.DB

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")
}
