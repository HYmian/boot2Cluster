package main

import (
	"flag"
	"log"
	"net"
	"strings"

	"github.com/HYmian/boot2Cluster/conf"
	"github.com/HYmian/boot2Cluster/connector"
)

var (
	clusterNum *uint   = flag.Uint("n", 0, "设置集群数量")
	port       *int    = flag.Int("p", 34616, "监听端口")
	config     *string = flag.String("conf", "./conf.yml", "配置文件路径")
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags)

	cfg, err := conf.LoadConfig(*config)
	if err != nil {
		log.Printf("load config error: %s", err.Error())
	}

	boot := conf.NewBoot(cfg, *clusterNum)

	l, err := net.ListenTCP(
		"tcp4",
		&net.TCPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: *port,
		})
	if err != nil {
		log.Fatalf("start listen error: %s", err.Error())
	}
	log.Printf("start listen in 0.0.0.0:%d", *port)

	c := make(chan *connector.Conn, 10)
	go waitConn(c, boot)

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			log.Printf("accept conn error: %s", err.Error())
			continue
		}
		log.Printf("accept a connection from %s", conn.RemoteAddr().String())

		c <- connector.NewConn(conn)
	}
}

func waitConn(c <-chan *connector.Conn, boot *conf.Boot) {
	i := uint(0)
	for {
		select {
		case conn := <-c:
			agent := conn.RemoteAddr().String()
			ip := strings.Split(agent, ":")[0]

			i++

			if err := conn.WriteRegister(); err != nil {
				log.Printf("register to agent error: %s", err.Error())
				break
			}

			data, err := conn.ReadPacket()
			if err != nil {
				log.Printf("write to agent ok error: %s", err.Error())
			}

			if data[0] == connector.COM_REGISTER {
				boot.AddNode(string(data[1:]), ip, i)

				conn.Close()
			}

			if boot.LiveCommand != "" {
				if err = boot.ExecLiveCommand(); err != nil {
					log.Printf("exec live command error: %s", err.Error())
					break
				}
			}

			if *clusterNum == i && boot.BootCommand != "" {
				//init cluster
				if err := boot.ExecBootCommand(); err != nil {
					log.Printf("exec boot command error: %s", err.Error())
				}
				break
			}
		}
	}
}
