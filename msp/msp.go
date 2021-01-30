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
	"os"
	"strconv"

	"github.com/dyfromnil/pdag/globleconfig"
)

// IdentityProvider for
type IdentityProvider interface {
	GenRsaKeys()
	GetKeyPair() (prvkey, pubkey []byte)
	RsaSignWithSha256(data []byte, keyBytes []byte) []byte
	RsaVerySignWithSha256(data, signData, keyBytes []byte) bool
	GetPubKey(nodeID string) []byte
	GetPivKey(nodeID string) []byte
	GetSelfPivKey() []byte
	GetSelfPubKey() []byte
	GetClusterAddrs() map[string]string
	GetNodeID() string
	GetAddr() string
}

// Identity for
type Identity struct {
	nodeID     string
	addr       string
	rsaPrivKey []byte
	rsaPubKey  []byte
	cluster    map[string]string
}

// NewIdt for
func NewIdt(nodeid string) Identity {
	return Identity{
		nodeID:     nodeid,
		addr:       globleconfig.NodeTable[nodeid],
		rsaPrivKey: getPivKey(nodeid),
		rsaPubKey:  getPubKey(nodeid),
		cluster:    globleconfig.NodeTable,
	}
}

// GenRsaKeys for
func (iden *Identity) GenRsaKeys() {
	if !isExist("./Keys") {
		fmt.Println("检测到还未生成公私钥目录，正在生成公私钥 ...")
		err := os.Mkdir("Keys", 0666)
		if err != nil {
			log.Panic()
		}

		n := len(globleconfig.NodeTable)

		for i := 0; i <= n; i++ {
			if !isExist("./Keys/N" + strconv.Itoa(i)) {
				err := os.Mkdir("./Keys/N"+strconv.Itoa(i), 0666)
				if err != nil {
					log.Panic()
				}
			}
			priv, pub := iden.GetKeyPair()
			privFileName := "Keys/N" + strconv.Itoa(i) + "/N" + strconv.Itoa(i) + "_RSA_PIV"
			file, err := os.OpenFile(privFileName, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				log.Panic(err)
			}
			defer file.Close()
			file.Write(priv)

			pubFileName := "Keys/N" + strconv.Itoa(i) + "/N" + strconv.Itoa(i) + "_RSA_PUB"
			file2, err := os.OpenFile(pubFileName, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				log.Panic(err)
			}
			defer file2.Close()
			file2.Write(pub)
		}
		fmt.Println("已为节点们生成RSA公私钥")
	}
}

//GetKeyPair for
func (iden *Identity) GetKeyPair() (prvkey, pubkey []byte) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	prvkey = pem.EncodeToMemory(block)
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubkey = pem.EncodeToMemory(block)
	return
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
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
	return getPubKey(nodeID)
}

func getPubKey(nodeID string) []byte {
	key, err := ioutil.ReadFile("Keys/" + nodeID + "/" + nodeID + "_RSA_PUB")
	if err != nil {
		log.Panic(err)
	}
	return key
}

// GetPivKey for
func (iden *Identity) GetPivKey(nodeID string) []byte {
	return getPivKey(nodeID)
}

func getPivKey(nodeID string) []byte {
	key, err := ioutil.ReadFile("Keys/" + nodeID + "/" + nodeID + "_RSA_PIV")
	if err != nil {
		log.Panic(err)
	}
	return key
}

//GetSelfPivKey for
func (iden *Identity) GetSelfPivKey() []byte {
	return iden.rsaPrivKey
}

//GetSelfPubKey for
func (iden *Identity) GetSelfPubKey() []byte {
	return iden.rsaPubKey
}

// GetClusterAddrs for
func (iden *Identity) GetClusterAddrs() map[string]string {
	return iden.cluster
}

//GetNodeID for
func (iden *Identity) GetNodeID() string {
	return iden.nodeID
}

//GetAddr for
func (iden *Identity) GetAddr() string {
	return iden.addr
}