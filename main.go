package main

import (
	"context"
	"log"

	"github.com/koshkaj/bloq/node"
	"github.com/koshkaj/bloq/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	makeNode(":3000", []string{})
	makeNode(":4000", []string{":3000"})
	// go func() {
	// 	for {
	// 		time.Sleep(1 * time.Second)
	// 		makeTransaction(listenAddr)
	// 	}
	// }()
	// log.Fatal(node.Start(listenAddr))
	select {}
}

func makeNode(listenAddr string, bootstrapNodes []string) *node.Node {
	n := node.New()
	go n.Start(listenAddr)
	if len(bootstrapNodes) > 0 {
		if err := n.BootstrapNetwork(bootstrapNodes); err != nil {
			log.Fatal(err)
		}
	}
	return n
}

func makeTransaction(listenAddr string) {
	client, err := grpc.Dial(listenAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	c := proto.NewNodeClient(client)
	v := &proto.Version{
		Version:    "bloq-1",
		Height:     1,
		ListenAddr: ":4000",
	}

	_, err = c.Handshake(context.TODO(), v)
	if err != nil {
		log.Fatal(err)
	}
}
