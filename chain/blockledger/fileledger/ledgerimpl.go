package fileledger

import (
	"github.com/dyfromnil/pdag/chain/blkstorage"
	cb "github.com/dyfromnil/pdag/proto-go/common"
	"github.com/golang/protobuf/proto"
)

// FileLedger is a struct used to interact with a node's ledger
type FileLedger struct {
	blockStore BlockStore
	// signal     chan struct{}
}

// BlockStore defines the interface to interact with deliver when using a
// file ledger
type BlockStore interface {
	AddBlock(block *cb.Block, tips []*cb.Block) error
	TipsBlock() []*cb.Block
}

// NewFileLedger creates a new FileLedger for interaction with the ledger
func NewFileLedger(blockStore BlockStore) *FileLedger {
	// return &FileLedger{blockStore: blockStore, signal: make(chan struct{})}
	return &FileLedger{blockStore: blockStore}
}

//TipsBlock for
func (fl *FileLedger) TipsBlock() []*cb.Block {
	return fl.blockStore.TipsBlock()
}

// Append a new block to the ledger
func (fl *FileLedger) Append(block *cb.Block, tips []*cb.Block) error {
	err := fl.blockStore.AddBlock(block, tips)
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
func (fl *FileLedger) CreateNextBlock(messages []*cb.Envelope) (*cb.Block, []*cb.Block) {
	var previousBlockHash [][]byte
	var err error

	blockHeaders := []*cb.BlockHeader{}
	tipsList := fl.TipsBlock()
	for _, bk := range tipsList {
		blockHeaders = append(blockHeaders, bk.Header)
	}

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

	block := blkstorage.NewBlock(previousBlockHash)
	block.Header.DataHash = blkstorage.BlockDataHash(data)
	block.Data = data

	return block, tipsList
}
