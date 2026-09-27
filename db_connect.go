package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB(short_url string, original_url string) {
	connStr := "host=127.0.0.1 port=5432 user=postgres password=password123 dbname=urlshortener sslmode=disable"

	var err error

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database!")

}

func insertURL(shortURL, originalURL string) error {
	sqlStatement := `INSERT INTO url_shortener (short_url, original_url, created_at) VALUES ($1, $2, NOW())`
	_, err := db.Exec(sqlStatement, shortURL, originalURL)
	return err
}

func getOriginalURL(shortCode string) (string, error) {
	sqlStatement := `SELECT original_url FROM url_shortener WHERE short_url = $1`
	var originalURL string
	err := db.QueryRow(sqlStatement, shortCode).Scan(&originalURL)
	return originalURL, err
}
