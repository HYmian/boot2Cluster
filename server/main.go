package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/HYmian/boot2Cluster/conf"
	"github.com/HYmian/boot2Cluster/connector"
)

var (
	clusterNum *string = flag.String("n", "0", "设置集群数量")
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags)

	l, err := net.ListenTCP(
		"tcp4",
		&net.TCPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: 34616,
		})
	if err != nil {
		log.Fatalf("start listen error: %s", err.Error())
	}

	c := make(chan *connector.Conn, 10)
	go waitConn(c)

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			log.Printf("accept conn error: %s", err.Error())
			continue
		}

		c <- connector.NewConn(conn)
	}
}

func waitConn(c <-chan *connector.Conn) {
	num := 2

	for {
		select {
		case conn := <-c:
			agent := conn.RemoteAddr().String()
			ip := strings.Split(agent, ":")[0]

			if num == 0 {
				//new agent
			}

			num--

			if num == 0 {
				//init cluster
				if err := conf.Exec("start-dfs.sh"); err != nil {
					log.Printf("exec start-dfs.sh error: %s", err.Error())
					break
				}
			}

			data, err := conn.ReadPacket()
			if err != nil {
				log.Printf("write to agent ok error: %s", err.Error())
			}

			if data[0] == connector.COM_REGISTER {
				cmd := fmt.Sprintf(`echo %s >> $HADOOP_HOME/etc/hadoop/slaves`, data[1:])
				if err := conf.Exec(cmd); err != nil {
					log.Printf("exec register error: %s", err.Error())
					break
				}

				cmd = fmt.Sprintf(`echo "%s %s" >> /etc/hosts`, ip, data[1:])
				if err := conf.Exec(cmd); err != nil {
					log.Printf("exec register error: %s", err.Error())
					break
				}

				conn.Close()
			}
		}
	}
}
