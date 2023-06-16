package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/koshkaj/bloq/node"
	"github.com/koshkaj/bloq/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	node := node.New()
	listenAddr := ":3000"
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	proto.RegisterNodeServer(grpcServer, node)
	fmt.Printf("node running on port %s\n", listenAddr)
	go func() {
		for {
			time.Sleep(1 * time.Second)
			makeTransaction(listenAddr)
		}
	}()
	grpcServer.Serve(ln)
}

func makeTransaction(listenAddr string) {
	client, err := grpc.Dial(listenAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	c := proto.NewNodeClient(client)
	v := &proto.Version{
		Version: "bloq-1",
		Height:  1,
	}

	_, err = c.Handshake(context.TODO(), v)
	if err != nil {
		log.Fatal(err)
	}
}
