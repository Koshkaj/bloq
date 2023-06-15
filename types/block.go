package types

import (
	"crypto/sha256"
	"log"

	pb "github.com/golang/protobuf/proto"
	"github.com/koshkaj/bloq/crypto"
	"github.com/koshkaj/bloq/proto"
)

func SignBlock(pk *crypto.PrivateKey, b *proto.Block) *crypto.Signature {
	return pk.Sign(HashBlock(b))
}

// returns SHA256 of the header
func HashBlock(block *proto.Block) []byte {
	b, err := pb.Marshal(block)
	if err != nil {
		log.Fatal(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]
}
