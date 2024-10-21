package db

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"flowey/utils"
)

var db *sql.DB

func Main(args []string) {
	defaultPath, err := utils.GetDefaultPath()
	if err != nil {
		log.Fatal(err)
	}

	flagSet := flag.NewFlagSet("flowey db", flag.ExitOnError)
	path := flagSet.String("path", defaultPath, "path to the database file")

	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, `Usage of flowey db:
  add       add a user to the database
  prepare   prepare a database`)
		fmt.Fprintln(os.Stderr)
		flagSet.PrintDefaults()
	}

	flagSet.Parse(args)
	nextArgs := utils.PopSlice(flagSet.Args())

	if len(args) < 1 {
		flagSet.Usage()
		return
	}

	switch flagSet.Arg(0) {
	case "add":
		if err := AddCmd(nextArgs, *path); err != nil {
			log.Fatal(err)
		}
	case "prepare":
		if err := PrepareCmd(nextArgs, *path); err != nil {
			log.Fatal(err)
		}
	default:
		flagSet.Usage()
	}
}
