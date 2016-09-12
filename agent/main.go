package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/HYmian/boot2Cluster/connector"
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags)

	server := os.Getenv("SERVER")

	conn, err := net.Dial("tcp4", server)
	if err != nil {
		log.Fatalf("dail to server %s:%d error: %s", server, 34616, err.Error())
	}

	c := connector.NewConn(conn)
	c.WriteRegister()
	c.Close()
}
