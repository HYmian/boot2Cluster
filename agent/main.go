package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/HYmian/boot2Cluster/conf"
	"github.com/HYmian/boot2Cluster/connector"
)

var (
	serverAddress *string = flag.String("s", "", "server address and port")
	config        *string = flag.String("conf", "./conf.yml", "配置文件路径")
	port          *int    = flag.Int("p", 8602, "监听端口")
	host          *string = flag.String("host", "", "注册用的hostname")
	mode          *string = flag.String("mode", "client", "运行模式")

	boot *conf.Boot
)

type server struct{}

func (s *server) Notify(ctx context.Context, in *connector.Notification) (*connector.Empty, error) {
	boot.Nodes = make([]conf.Node, 0, len(in.Inform))
	for _, inform := range in.Inform {
		boot.Nodes = append(boot.Nodes, inform.Node)
		log.Println(inform.Node)
	}

	boot.Entry()

	return &connector.Empty{}, nil
}

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg, err := conf.LoadConfig(*config)
	if err != nil {
		log.Printf("load config error: %s", err.Error())
	}

	boot = conf.NewBoot(cfg, 1)

	address := strings.Split(*serverAddress, ";")

	for _, a := range address {
		co, err := grpc.Dial(a, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		c := connector.NewRegisterClient(co)

		host, err := os.Hostname()
		if err != nil {
			log.Printf("get hostname error %s", err.Error())
		}
		_, err = c.Registe(context.Background(),
			&connector.Inform{
				Node: map[string]string{
					"host": host,
					"port": fmt.Sprintf("%d", *port),
					"mode": *mode,
				},
			},
		)
		if err != nil {
			log.Printf("register error: %s", err.Error())
		}
		co.Close()
	}

	if *mode == "client" {
		return
	}

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

	s := grpc.NewServer()
	connector.RegisterNotifierServer(s, &server{})
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
