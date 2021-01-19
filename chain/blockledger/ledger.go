package blockledger

import (
	cb "github.com/dyfromnil/pdag/proto-go/common"
)

//Reader for
type Reader interface {
	TipsBlock() []*cb.Block
	//CreateNextBlock returns the next block and tips blocks of ledger
	CreateNextBlock(messages []*cb.Envelope) (*cb.Block, []*cb.Block)
}

// Writer allows the caller to modify the ledger
type Writer interface {
	// Append a new block to the ledger and change the tips state of ledger
	Append(block *cb.Block, tips []*cb.Block) error
}

// ReadWriter encapsulates the read/write functions of the ledger
type ReadWriter interface {
	Reader
	Writer
}
