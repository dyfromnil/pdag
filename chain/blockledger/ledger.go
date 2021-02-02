package blockledger

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/dyfromnil/pdag/chain/blkstorage"
	cb "github.com/dyfromnil/pdag/proto-go/common"
)

// ReadWriter encapsulates the read/write functions of the ledger
type ReadWriter interface {
	CreateNextBlock(messages []*cb.Envelope) *cb.Block
	Append(block *cb.Block) error
	VerifyCurrentBlock(*cb.Block) bool
}

// FileLedger is a struct used to interact with a node's ledger
type FileLedger struct {
	blockStore      blkstorage.BlockStore
	round           int                    //轮次
	roundDigestPost map[int]map[string]int //第几轮哪个区块的目前的后继区块数
	preBlocksDigest map[string][]string    //当前区块的前驱区块,所有的string均为区块的digest
	digestToHash    map[string][]byte      //digest对应的hash值

	// signal     chan struct{}
}

// NewFileLedger creates a new FileLedger for interaction with the ledger
func NewFileLedger(blockStore blkstorage.BlockStore) *FileLedger {
	fl := &FileLedger{
		blockStore:      blockStore,
		round:           1,
		roundDigestPost: make(map[int]map[string]int),
		preBlocksDigest: make(map[string][]string),
	}
	fl.blockStore.AddBlock(gensisBlock())
	return fl
}

func gensisBlock() *cb.Block {
	time := time.Now()
	gensis := fmt.Sprintln("This is gensis block, created ", time)
	gensisHash := sha256.Sum256([]byte(gensis))
	selfHash := [][]byte{}
	selfHash = append(selfHash, gensisHash[:])

	header := &cb.BlockHeader{
		PreviousHash: selfHash,
		DataHash:     gensisHash[:],
		Round:        0,
		Timestamp:    time.UnixNano(),
	}
	gensisBlock := &cb.Block{
		Header: header,
		Data:   &cb.BlockData{Data: selfHash},
	}
	return gensisBlock
}

// Append a new block to the ledger
func (fl *FileLedger) Append(block *cb.Block, tips []*cb.Block) error {
	err := fl.blockStore.AddBlock(block)
	// if err == nil {
	// 	close(fl.signal)
	// 	fl.signal = make(chan struct{})
	// }
	return err
}

// CreateNextBlock provides a utility way to construct the next block from
// contents and metadata for a given ledger
// XXX This will need to be modified to accept marshaled envelopes
//     to accommodate non-deterministic marshaling
func (fl *FileLedger) CreateNextBlock(messages []*cb.Envelope) *cb.Block {
	var previousBlockHash [][]byte
	var err error

	blockHeaders := []*cb.BlockHeader{}

	previousBlockHash = blkstorage.BlocksHeaderHash(blockHeaders)

	data := &cb.BlockData{
		Data: make([][]byte, len(messages)),
	}

	for i, msg := range messages {
		data.Data[i], err = proto.Marshal(msg)
		if err != nil {
			panic(err)
		}
	}

	header := &cb.BlockHeader{
		PreviousHash: previousBlockHash,
		DataHash:     blkstorage.BlockDataHash(data),
	}
	block := blkstorage.NewBlock(header, data)

	return block
}

//VerifyCurrentBlock for
func (fl *FileLedger) VerifyCurrentBlock(*cb.Block) bool {
	return true
}

//为多重映射开辟赋值
func (fl *FileLedger) setRoundDigestPost(val int, val2 string, b int) {
	if _, ok := fl.roundDigestPost[val]; !ok {
		fl.roundDigestPost[val] = make(map[string]int)
	}
	fl.roundDigestPost[val][val2] = b
}
