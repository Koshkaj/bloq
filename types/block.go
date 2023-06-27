package types

import (
	"bytes"
	"crypto/sha256"
	"log"

	"github.com/cbergoon/merkletree"
	pb "github.com/golang/protobuf/proto"
	"github.com/koshkaj/bloq/crypto"
	"github.com/koshkaj/bloq/proto"
)

type TxHash struct {
	hash []byte
}

func NewTxHash(hash []byte) TxHash {
	return TxHash{
		hash: hash,
	}
}

func (h TxHash) CalculateHash() ([]byte, error) {
	return h.hash, nil
}

func (h TxHash) Equals(other merkletree.Content) (bool, error) {
	equals := bytes.Equal(h.hash, other.(TxHash).hash)
	return equals, nil
}

func VerifyBlock(b *proto.Block) bool {
	if len(b.Transactions) > 0 {
		if !VerifyRootHash(b) {
			return false
		}
	}
	if len(b.PublicKey) != crypto.PubKeyLen {
		return false
	}
	if len(b.Signature) != crypto.SignatureLen {
		return false
	}
	sig := crypto.SignatureFromBytes(b.Signature)
	pubKey := crypto.PublicKeyFromBytes(b.PublicKey)
	hash := HashBlock(b)
	return sig.Verify(pubKey, hash)
}

func SignBlock(pk *crypto.PrivateKey, b *proto.Block) *crypto.Signature {
	if len(b.Transactions) > 0 {
		tree, err := GetMerkleTree(b)
		if err != nil {
			log.Fatal(err)
		}
		b.Header.RootHash = tree.MerkleRoot()
	}
	hash := HashBlock(b)
	sig := pk.Sign(hash)
	b.PublicKey = pk.Public().Bytes()
	b.Signature = sig.Bytes()
	return sig
}

// returns SHA256 of the header
func HashBlock(block *proto.Block) []byte {
	return HashHeader(block.Header)
}

func VerifyRootHash(b *proto.Block) bool {
	tree, err := GetMerkleTree(b)
	if err != nil {
		return false
	}
	valid, err := tree.VerifyTree()
	if err != nil {
		return false
	}
	if !valid {
		return false
	}
	return bytes.Equal(b.Header.RootHash, tree.MerkleRoot())
}

func GetMerkleTree(b *proto.Block) (*merkletree.MerkleTree, error) {
	list := make([]merkletree.Content, len(b.Transactions))
	for i := 0; i < len(b.Transactions); i++ {
		list[i] = NewTxHash(HashTransaction(b.Transactions[i]))
	}
	t, err := merkletree.NewTree(list)
	if err != nil {
		return nil, err
	}
	b.Header.RootHash = t.MerkleRoot()
	return t, nil
}

func HashHeader(header *proto.Header) []byte {
	b, err := pb.Marshal(header)
	if err != nil {
		log.Fatal(err)
	}
	hash := sha256.Sum256(b)
	return hash[:]
}
