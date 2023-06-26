package types

import (
	"crypto/sha256"
	"log"

	pb "github.com/golang/protobuf/proto"
	"github.com/koshkaj/bloq/crypto"
	"github.com/koshkaj/bloq/proto"
)

func SignTransaction(pk *crypto.PrivateKey, tx *proto.Transaction) *crypto.Signature {
	return pk.Sign(HashTransaction(tx))
}

func HashTransaction(tx *proto.Transaction) []byte {
	b, err := pb.Marshal(tx)
	if err != nil {
		log.Fatal(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]

}

func VerifyTransaction(tx *proto.Transaction) bool {
	for _, inp := range tx.Inputs {
		if len(inp.Signature) == 0 {
			log.Fatal("the transaction has no signature")
		}
		var (
			sig    = crypto.SignatureFromBytes(inp.Signature)
			pubKey = crypto.PublicKeyFromBytes(inp.PublicKey)
		)
		inp.Signature = nil
		if !sig.Verify(pubKey, HashTransaction(tx)) {
			return false
		}
	}
	return true
}
