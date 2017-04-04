package database

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var dbNameFlag string
var dbIpFlag string
var dbUserFlag string
var dbPassFlag string

type DB struct {
	Db *sql.DB
}

func init() {
	flag.StringVar(&dbIpFlag, "db-ip", "127.0.0.1", "Database IP address")
	flag.StringVar(&dbNameFlag, "db-name", "data", "Database Name")
	flag.StringVar(&dbUserFlag, "db-user", "postgres", "Database User")
	flag.StringVar(&dbPassFlag, "db-pass", "postgres", "Database User Password")
}

func Connect(db *DB) error {
	var err error
	os.Setenv("PGSSLMODE", "disable")
	(*db).Db, err = sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s host=%s password=%s", dbUserFlag, dbNameFlag, dbIpFlag, dbPassFlag))
	return err
}

func GetResult(db *DB, query string) (*sql.Rows, error) {
	var err error

	rows, err := db.Db.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return rows, nil
}

func InsertArray(db *DB, table string, values []string) error {
	var err error
	var stmt *sql.Stmt
	var query string

	query = fmt.Sprintf("insert into %s values ($1", table)
	for i := 1; i < len(values); i++ {
		query = fmt.Sprintf("%s,$%d", query, i+1)
	}
	query = fmt.Sprintf("%s)", query)

	stmt, err = db.Db.Prepare(query)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()

	values2 := make([]interface{}, len(values))
	for index, value := range values {
		if value == "" {
			values2[index] = nil
		} else {
			values2[index] = value
		}
	}

	_, err = stmt.Exec(values2...)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
