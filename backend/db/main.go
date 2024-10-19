package db

import (
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func Main() {
	if len(os.Args) < 3 {
		goto usage
	}

	switch os.Args[2] {
	case "init":
		if err := runInit(); err != nil {
			log.Fatal(err)
		}
		return
	}

usage:
	fmt.Fprintln(os.Stderr, `Usage of flowey db:
  init   initialize a new database`)
}
