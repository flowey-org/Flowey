package db

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"flowey/utils"
)

func getPath() string {
	dataDir, err := utils.GetDataDir()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(dataDir, "flowey.db")
}

func Main(args []string) {
	flagSet := flag.NewFlagSet("flowey db", flag.ExitOnError)
	path := flagSet.String("path", getPath(), "path to the database file")

	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, `Usage of flowey db:
  add    add a user to the database
  init   initialize a new database`)
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
		if err := Add(nextArgs, *path); err != nil {
			log.Fatal(err)
		}
		return
	case "init":
		if err := Init(nextArgs, *path); err != nil {
			log.Fatal(err)
		}
		return
	}

	flagSet.Usage()
}
