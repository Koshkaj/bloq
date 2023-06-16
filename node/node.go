package node

import (
	"context"
	"fmt"

	"github.com/koshkaj/bloq/proto"
	"google.golang.org/grpc/peer"
)

type Node struct {
	version string
	// peers map[net.Addr]*grpc.ClientConn
	proto.UnimplementedNodeServer
}

func New() *Node {
	return &Node{
		version: "bloq-1",
	}
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	peer, _ := peer.FromContext(ctx)
	fmt.Println("received tx from :", peer)
	return &proto.Ack{}, nil
}

func (n *Node) Handshake(ctx context.Context, v *proto.Version) (*proto.Version, error) {
	ourVersion := &proto.Version{
		Version: n.version,
		Height:  100,
	}
	peer, _ := peer.FromContext(ctx)
	fmt.Printf("received version from %s: %+v\n", v, peer.Addr)
	return ourVersion, nil
}
