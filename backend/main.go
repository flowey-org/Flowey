package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"flowey/db"
	"flowey/server"
	"flowey/utils"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	flagSet := flag.NewFlagSet("flowey", flag.ExitOnError)

	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, `Usage of flowey:
  db       interact with the database
  server   run the server`)
	}

	if len(os.Args) < 2 {
		flagSet.Usage()
		return
	}

	args := os.Args[1:]
	flagSet.Parse(args)
	nextArgs := utils.PopSlice(flagSet.Args())

	switch args[0] {
	case "db":
		db.Main(nextArgs)
		return
	case "server":
		if err := server.Main(); err != nil {
			log.Fatal(err)
		}
		return
	}

	flagSet.Usage()
}
