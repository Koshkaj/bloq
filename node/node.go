package node

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/koshkaj/bloq/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
)

type Node struct {
	version    string
	listenAddr string
	logger     *zap.SugaredLogger

	peerLock sync.RWMutex
	peers    map[proto.NodeClient]*proto.Version

	proto.UnimplementedNodeServer
}

func New() *Node {
	loggerconfig := zap.NewDevelopmentConfig()
	loggerconfig.EncoderConfig.TimeKey = ""
	logger, _ := loggerconfig.Build()
	return &Node{
		version: "bloq-1",
		peers:   make(map[proto.NodeClient]*proto.Version),
		logger:  logger.Sugar(),
	}
}

func (n *Node) BootstrapNetwork(addrs []string) error {
	for _, addr := range addrs {
		c, err := makeNodeClient(addr)
		if err != nil {
			return err
		}
		v, err := c.Handshake(context.Background(), n.getVersion())
		if err != nil {
			fmt.Println("handshake error : ", err)
			continue
		}

		n.addPeer(c, v)
	}
	return nil
}

func (n *Node) getVersion() *proto.Version {
	return &proto.Version{
		Version:    "bloq-1",
		Height:     0,
		ListenAddr: n.listenAddr,
		PeerList:   n.getPeerList(),
	}
}

func (n *Node) getPeerList() []string {
	n.peerLock.RLock()
	defer n.peerLock.RUnlock()

	peers := []string{}
	for _, version := range n.peers {
		peers = append(peers, version.ListenAddr)
	}
	return peers

}

func (n *Node) Start(listenAddr string) error {
	n.listenAddr = listenAddr
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	proto.RegisterNodeServer(grpcServer, n)

	// go n.bootstrapNetwork()

	n.logger.Info("node running: ", n.listenAddr)
	return grpcServer.Serve(ln)
}

func (n *Node) addPeer(c proto.NodeClient, v *proto.Version) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	n.logger.Debugf("[%s] new peer connected (%s) - height (%d)", n.listenAddr, v.ListenAddr, v.Height)
	n.peers[c] = v
}

func (n *Node) deletePeer(c proto.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()
	delete(n.peers, c)
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	peer, _ := peer.FromContext(ctx)
	fmt.Println("received tx from :", peer)
	return &proto.Ack{}, nil
}

func (n *Node) Handshake(ctx context.Context, v *proto.Version) (*proto.Version, error) {
	c, err := makeNodeClient(v.ListenAddr)
	if err != nil {
		return nil, err
	}
	n.addPeer(c, v)
	return n.getVersion(), nil
}

func makeNodeClient(listenAddr string) (proto.NodeClient, error) {
	c, err := grpc.Dial(listenAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return proto.NewNodeClient(c), nil
}
