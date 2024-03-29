package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {

	connStr, ok := os.LookupEnv("DB_CONNECTION_STRING")

	if !ok {
		log.Println("DB_CONNECTION_STRING variable not set")
	}
	if connStr == "" {
		log.Fatal("DB_CONNECTION_STRING environment variable not set")
	}

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	log.Println("Connected to the database")

	return db, nil
}
