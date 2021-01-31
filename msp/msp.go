package msp

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/dyfromnil/pdag/globleconfig"
)

// IdentityProvider for
type IdentityProvider interface {
	RsaSignWithSha256(data []byte, keyBytes []byte) []byte
	RsaVerySignWithSha256(data, signData, keyBytes []byte) bool
	GetPubKey(nodeID string) []byte
	GetPivKey(nodeID string) []byte
	GetSelfPivKey() []byte
	GetSelfPubKey() []byte
	GetClusterAddrs() map[string]string
	GetNodeID() string
	GetAddr(nodeID string) string
	GetSelfAddr() string
}

type keyPair struct {
	privKey []byte
	pubKey  []byte
}

// Identity for
type Identity struct {
	nodeID  string
	addr    string
	keyPair keyPair
	// key:nodeID, value:addr
	clusterAddrs map[string]string
	// key:nodeID, value:keyPair
	clusterKeyPairs map[string]keyPair
}

// NewIdt for
func NewIdt(nodeid string) *Identity {
	cKPs := initClusterKeyPairs()

	return &Identity{
		nodeID:          nodeid,
		addr:            globleconfig.NodeTable[nodeid],
		keyPair:         cKPs[nodeid],
		clusterAddrs:    globleconfig.NodeTable,
		clusterKeyPairs: cKPs,
	}
}

func initClusterKeyPairs() map[string]keyPair {
	cKPs := make(map[string]keyPair)
	for key := range globleconfig.NodeTable {
		cKPs[key] = keyPair{
			getPivKeyFromFile(key),
			getPubKeyFromFile(key),
		}
	}
	return cKPs
}

//RsaSignWithSha256 for
func (iden *Identity) RsaSignWithSha256(data []byte, keyBytes []byte) []byte {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("private key error"))
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		panic(err)
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		panic(err)
	}

	return signature
}

//RsaVerySignWithSha256 for
func (iden *Identity) RsaVerySignWithSha256(data, signData, keyBytes []byte) bool {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], signData)
	if err != nil {
		panic(err)
	}
	return true
}

// GetPubKey for
func (iden *Identity) GetPubKey(nodeID string) []byte {
	return iden.clusterKeyPairs[nodeID].pubKey
}

// GetPivKey for
func (iden *Identity) GetPivKey(nodeID string) []byte {
	return iden.clusterKeyPairs[nodeID].privKey
}

func getPubKeyFromFile(nodeID string) []byte {
	key, err := ioutil.ReadFile("Keys/" + nodeID + "/" + nodeID + "_RSA_PUB")
	if err != nil {
		log.Panic(err)
	}
	return key
}

func getPivKeyFromFile(nodeID string) []byte {
	key, err := ioutil.ReadFile("Keys/" + nodeID + "/" + nodeID + "_RSA_PIV")
	if err != nil {
		log.Panic(err)
	}
	return key
}

//GetSelfPivKey for
func (iden *Identity) GetSelfPivKey() []byte {
	return iden.keyPair.privKey
}

//GetSelfPubKey for
func (iden *Identity) GetSelfPubKey() []byte {
	return iden.keyPair.pubKey
}

// GetClusterAddrs for
func (iden *Identity) GetClusterAddrs() map[string]string {
	return iden.clusterAddrs
}

//GetNodeID for
func (iden *Identity) GetNodeID() string {
	return iden.nodeID
}

//GetSelfAddr for
func (iden *Identity) GetSelfAddr() string {
	return iden.addr
}

//GetAddr for
func (iden *Identity) GetAddr(nodeID string) string {
	return iden.clusterAddrs[nodeID]
}
