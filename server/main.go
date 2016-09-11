package server

import (
	"flag"
	"log"
	"net"
)

var (
	clusterNum *string = flag.String("n", "0", "设置集群数量")
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags)

	l, err := net.ListenTCP("tcp4", &net.TCPAddr{
		IP:   []byte("0.0.0.0"),
		Port: 34616,
	})
	if err != nil {
		log.Fatalf("start listen error: %s", err.Error())
	}

	c := make(chan net.Conn, 10)
	go waitConn(c)

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			log.Printf("accept conn error: %s", err.Error())
			continue
		}

		c <- conn
	}
}

func waitConn(c <-chan net.Conn) {
	num := 10

	for {
		select {
		case conn := <-c:
			agent := conn.RemoteAddr().String()
			log.Println(agent) // TODO

			if num == 0 {
				//new agent
			}

			num--

			if num == 0 {
				//init cluster
			}

			_, err := conn.Write([]byte("OK"))
			if err != nil {
				log.Printf("write to agent ok error: %s", err.Error())
			}
		}
	}
}
