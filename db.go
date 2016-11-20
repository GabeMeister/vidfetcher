package main

import (
	"database/sql"
	"fmt"
	"log"
)

// CreateDBInstance creates an instance of a database connection to
// be used throughout the duration of the program
func CreateDBInstance() *sql.DB {
	dbinfo := dbConnectionStr()
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func dbConnectionStr() string {
	return fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable",
		"gabemeister", "qwerQWER1234!", "104.236.163.200", "youtubecollections")
}
