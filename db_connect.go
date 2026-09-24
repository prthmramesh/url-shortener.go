package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func connectDB() {
	connStr := "host=127.0.0.1 port=5432 user=postgres password=password123 dbname=urlshortener sslmode=disable"

	fmt.Println(connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database!")

	sqlStatement := `INSERT INTO url_shortener (short_url, original_url, created_at) VALUES ($1, $2, NOW())`
	_, err = db.Exec(sqlStatement, "short1", "https://www.example.com/long-url-1")
	if err != nil {
		log.Fatal(err)
	}

	sqlStatement2 := `SELECT short_url, original_url, created_at FROM url_shortener`
	rows, err := db.Query(sqlStatement2)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var shortURL, originalURL string
		var createdAt time.Time
		err = rows.Scan(&shortURL, &originalURL, &createdAt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Short URL: %s, Original URL: %s, Created At: %s\n", shortURL, originalURL, createdAt)
	}
}
