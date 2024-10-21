package db

import "log"

func Close() {
	if db != nil {
		db.Close()
		db = nil
		log.Printf("closed the connection to the database")
	}
}
