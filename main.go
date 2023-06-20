package main

import (
	"context"
	"log"
	"time"

	"github.com/koshkaj/bloq/node"
	"github.com/koshkaj/bloq/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	makeNode(":3000", []string{})
	time.Sleep(1 * time.Second)
	makeNode(":3001", []string{":3000"})
	time.Sleep(3 * time.Second)
	makeNode(":4001", []string{":3001"})
	// makeTransaction(listenAddr)
	// log.Fatal(node.Start(listenAddr))
	select {}
}

func makeNode(listenAddr string, bootstrapNodes []string) *node.Node {
	n := node.New()
	go n.Start(listenAddr, bootstrapNodes)
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
