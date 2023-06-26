package main

import (
	"context"
	"log"
	"time"

	"github.com/koshkaj/bloq/crypto"
	"github.com/koshkaj/bloq/node"
	"github.com/koshkaj/bloq/proto"
	"github.com/koshkaj/bloq/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	makeNode(":3000", []string{}, true)
	time.Sleep(1 * time.Second)
	makeNode(":3001", []string{":3000"}, false)
	time.Sleep(3 * time.Second)
	makeNode(":4001", []string{":3001"}, false)
	for {
		time.Sleep(time.Millisecond * 100)
		makeTransaction()
	}
}

func makeNode(listenAddr string, bootstrapNodes []string, isValidator bool) *node.Node {
	cfg := node.ServerConfig{
		Version:    "bloq-1.0",
		ListenAddr: listenAddr,
	}
	if isValidator {
		cfg.PrivateKey = crypto.GeneratePrivateKey()
	}
	n := node.New(cfg)
	go n.Start(listenAddr, bootstrapNodes)
	return n
}

func makeTransaction() {
	client, err := grpc.Dial(":3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	c := proto.NewNodeClient(client)
	privKey := crypto.GeneratePrivateKey()

	tx := &proto.Transaction{
		Inputs: []*proto.TxInput{
			{
				PrevTxHash:   util.RandomHash(),
				PrevOutIndex: 0,
				PublicKey:    privKey.Public().Bytes(),
			},
		},
		Outputs: []*proto.TxOutput{
			{
				Amount:  99,
				Address: privKey.Public().Address().Bytes(),
			},
		},
		Version: 1,
	}

	_, err = c.HandleTransaction(context.TODO(), tx)

	if err != nil {
		log.Fatal(err)
	}
}
