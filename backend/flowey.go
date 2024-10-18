package main

import (
	"flag"
	"fmt"
	"log"
)

func run(address string) error {
	server := NewServer(address)
	return server.ListenAndServe()
}

func main() {
	ip := flag.String("ip", "0.0.0.0", "ip to bind to")
	port := flag.Int("port", 80, "port to bind to")
	flag.Parse()

	address := fmt.Sprintf("%s:%d", *ip, *port)
	if err := run(address); err != nil {
		log.Fatal(err)
	}
}
