package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host = "url-microservice_postgres_1" // set this host for running service on Docker
	//host = "localhost" // set this host for running service on localhost
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
