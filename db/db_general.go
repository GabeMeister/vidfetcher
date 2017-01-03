package db

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	// For postgres db
	_ "github.com/lib/pq"
)

// CreateDBInstance creates an instance of a database connection to
// be used throughout the duration of the program
func CreateDBInstance() *sql.DB {
	dbinfo := dbConnectionStr()
	youtubeDB, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}

	return youtubeDB
}

func dbConnectionStr() string {
	return fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable",
		"gabemeister", "qwerQWER1234!", "104.236.163.200", "youtubecollections")
}

// SelectColumnFromTable fetches all the channel ids in the channels table
// maxRows of 0 does not put limit on number of rows returned
func SelectColumnFromTable(youtubeDB *sql.DB, column string, table string, maxRows uint) []string {
	var sql string
	if maxRows == 0 {
		sql = fmt.Sprintf("select %s from %s;", column, table)
	} else {
		sql = fmt.Sprintf("select %s from %s limit %d;", column, table, maxRows)
	}

	rows, err := youtubeDB.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var colValue string
		err := rows.Scan(&colValue)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, strings.TrimSpace(colValue))
	}

	return results
}

// GeneratePlaceHolders generates a comma separated list of postgres place holders
// i.e. $1, $2, $3
func GeneratePlaceHolders(count int) string {
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	return strings.Join(placeholders, ",")
}
