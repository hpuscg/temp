package main

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
)

const (
	host     = "192.168.100.235"
	port     = 5432
	user     = "postgres"
	password = "deepglint"
	dbname   = "libraT"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s " +
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("client to postgresql err: ", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("ping postgresql err :", err)
	} else {
		fmt.Println("successfully connected!")
	}
}
