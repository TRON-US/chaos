package main

import (
	registry "github.com/TRON-US/chaos/zookeeper"
	"github.com/TRON-US/chaos/zookeeper/balancer"
	"github.com/TRON-US/chaos/zookeeper/example/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	registry.RegisterResolver("zk", []string{"127.0.0.1:2181"}, "test", "v1.0")
	//c, err := grpc.Dial("zk:///", grpc.WithInsecure(), grpc.WithBalancerName(balancer.RoundRobin))
	balancer.InitRoundRobin()
	c, err := grpc.Dial("zk:///", grpc.WithInsecure(), grpc.WithBalancerName(balancer.RoundRobin))
	if err != nil {
		log.Printf("grpc dial: %s", err)
		return
	}
	defer c.Close()
	client := proto.NewTestClient(c)

	for i := 0; i < 5000; i++ {
		resp, err := client.Say(context.Background(), &proto.SayReq{Content: "round robin"})
		if err != nil {
			log.Println("aa:", err)
			time.Sleep(time.Second)
			continue
		}
		time.Sleep(time.Second)
		log.Printf(resp.Content)
	}
}
