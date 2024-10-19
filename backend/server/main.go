package server

import (
	"flag"
	"os"
)

func Main() error {
	flagSet := flag.NewFlagSet("flowey server", flag.ExitOnError)
	ip := flagSet.String("ip", "0.0.0.0", "ip to bind to")
	port := flagSet.Int("port", 80, "port to bind to")
	flagSet.Parse(os.Args[2:])

	if flagSet.NArg() > 0 {
		flagSet.Usage()
		return nil
	}

	server := NewServer(*ip, *port)
	return server.ListenAndServe()
}
