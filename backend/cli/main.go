package main

import (
	"database/sql"
	"log"
	"flag"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	path := flag.String("db", "flowey.db", "path to the database file")
	flag.Parse()

	db, err := sql.Open("sqlite3", *path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	query := `
CREATE TABLE users(
  id INTEGER NOT NULL PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL
);`
	if _, err := db.Exec(query); err != nil {
		log.Println(err)
		return
	}
}
