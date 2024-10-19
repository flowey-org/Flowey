package main

import (
	"fmt"
	"log"
	"os"

	"flowey/db"
	"flowey/server"
)

func main() {
	if len(os.Args) < 2 {
		goto usage
	}

	switch os.Args[1] {
	case "db":
		db.Main();
		return
	case "server":
		if err := server.Main(); err != nil {
			log.Fatal(err)
		}
		return
	}

usage:
	fmt.Fprintln(os.Stderr, `Usage of flowey:
  db       handle the database
  server   run the server`)
}
