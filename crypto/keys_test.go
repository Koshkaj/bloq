package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	privKey := GeneratePrivateKey()

	assert.Equal(t, len(privKey.Bytes()), privKeyLen)

	pubKey := privKey.Public()

	assert.Equal(t, len(pubKey.Bytes()), pubKeyLen)
}

func TestPrivateKeySign(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	msg := []byte("foo bar baz")

	sig := privKey.Sign(msg)
	assert.True(t, sig.Verify(pubKey, msg))

	assert.False(t, sig.Verify(pubKey, []byte("wsp")))

	invalidPrivKey := GeneratePrivateKey()
	invalidPubKey := invalidPrivKey.Public()
	assert.False(t, sig.Verify(invalidPubKey, msg))
}

func TestPublicKeyToAddress(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	address := pubKey.Address()

	assert.Equal(t, addressLen, len(address.Bytes()))
}

func TestNewPrivateKeyFromString(t *testing.T) {
	var (
		seed       = "2ec4d1620be036b2be86892effb4d6b3dd3f50262391b174e4a8628bb038360b"
		privKey    = NewPrivateKeyFromString(seed)
		addressStr = "79480dec6b8b0a299f3af77a4657a08493875c0b"
	)
	assert.Equal(t, privKeyLen, len(privKey.Bytes()))

	address := privKey.Public().Address()
	assert.Equal(t, address.String(), addressStr)
}
