package server

import (
	"flag"
	"log"
	"os"

	"flowey/db"
	"flowey/utils"
)

func Main() error {
	defaultPath, err := utils.GetDefaultPath()
	if err != nil {
		log.Fatal(err)
	}

	flagSet := flag.NewFlagSet("flowey server", flag.ExitOnError)
	ip := flagSet.String("ip", "0.0.0.0", "ip to bind to")
	path := flagSet.String("path", defaultPath, "path to the database file")
	port := flagSet.Int("port", 80, "port to bind to")
	flagSet.Parse(os.Args[2:])

	if flagSet.NArg() > 0 {
		flagSet.Usage()
		return nil
	}

	if err = db.Prepare(*path); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := NewServer(*ip, *port)
	err = server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
