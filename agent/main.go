package main

import (
	"flag"
	"log"
	"net"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/HYmian/boot2Cluster/conf"
	"github.com/HYmian/boot2Cluster/connector"
)

var (
	server *string = flag.String("s", "", "server address and port")
	config *string = flag.String("conf", "./conf.yml", "配置文件路径")
)

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags)

	cfg, err := conf.LoadConfig(*config)
	if err != nil {
		log.Printf("load config error: %s", err.Error())
	}

	boot := conf.NewBoot(cfg, 1)

	co, err := grpc.Dial(*server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer co.Close()
	cc := connector.NewRegisterClient(co)

	_, err = cc.Notify(context.Background(),
		&connector.Inform{Node: map[string]string{"host": "airCraft"}},
	)
	if err != nil {
		log.Printf("notify error: %s", err.Error())
	}
	return

	conn, err := net.Dial("tcp4", *server)
	if err != nil {
		log.Fatalf("dail to server %s error: %s", server, err.Error())
	}

	c := connector.NewConn(conn)

	agent := conn.RemoteAddr().String()
	ip := strings.Split(agent, ":")[0]

	data, err := c.ReadPacket()
	if err != nil {
		log.Printf("write to agent ok error: %s", err.Error())
	}

	if data[0] == connector.COM_REGISTER {
		boot.AddNode(string(data[1:]), ip, 1)
	}

	if err = c.WriteRegister(); err != nil {
		log.Printf("register to server error: %s", err.Error())
	}
	c.Close()
}
