package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	"github.com/HYmian/boot2Cluster/conf"
	"github.com/HYmian/boot2Cluster/connector"
)

var (
	clusterNum *uint   = flag.Uint("n", 0, "设置集群数量")
	port       *int    = flag.Int("p", 34616, "监听端口")
	config     *string = flag.String("conf", "./conf.yml", "配置文件路径")

	c = make(chan conf.Node, 10)
)

type server struct{}

func (s *server) Registe(ctx context.Context, in *connector.Inform) (*connector.Empty, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		log.Println("failed to get peer from ctx")
	}

	agent := pr.Addr.String()
	ip := strings.Split(agent, ":")[0]

	in.Node["IP"] = ip
	c <- in.Node

	return &connector.Empty{}, nil
}

func NewNotification(nodes conf.Nodes) *connector.Notification {
	informs := make([]*connector.Inform, 0, len(nodes))

	for _, node := range nodes {
		informs = append(informs, NewInform(node))
	}

	return &connector.Notification{Inform: informs}
}

func NewInform(node conf.Node) *connector.Inform {
	return &connector.Inform{Node: node}
}

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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

	go waitConn(c, boot)

	s := grpc.NewServer()
	connector.RegisterRegisterServer(s, &server{})
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func waitConn(c <-chan conf.Node, boot *conf.Boot) {
	for {
		select {
		case node := <-c:
			boot.AddNode(node)

			if boot.LiveCommand != "" {
				if err := boot.ExecLiveCommand(); err != nil {
					log.Printf("exec live command error: %s", err.Error())
					break
				}
			}

			for _, n := range boot.Nodes {
				co, err := grpc.Dial(fmt.Sprintf("%s:%s", n["IP"], n["port"]), grpc.WithInsecure())
				if err != nil {
					log.Fatalf("did not connect to agent %s: %s", n["Host"], err.Error())
				}
				c := connector.NewNotifierClient(co)

				if _, err := c.Notify(context.Background(), NewNotification(boot.Nodes)); err != nil {
					log.Printf("notify to agent %s error: %s", n["Host"], err.Error())
				}
				co.Close()
			}

			if *clusterNum == uint(len(boot.Nodes)) && boot.BootCommand != "" {
				//init cluster
				if err := boot.ExecBootCommand(); err != nil {
					log.Printf("exec boot command error: %s", err.Error())
				}
				break
			}
		}
	}
}
