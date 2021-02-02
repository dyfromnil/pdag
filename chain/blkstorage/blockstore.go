package blkstorage

import (
	cb "github.com/dyfromnil/pdag/proto-go/common"
)

// BlockStore defines the interface to interact with deliver when using a
// file ledger
type BlockStore interface {
	AddBlock(block *cb.Block) error
}

// BlkStore - filesystem based implementation for `BlockStore`
type BlkStore struct {
	conf    *Conf
	fileMgr *blockfileMgr
}

// NewBlockStore constructs a `BlockStore`
func NewBlockStore(conf *Conf) (*BlkStore, error) {
	fileMgr, err := newBlockfileMgr(conf)
	if err != nil {
		return nil, err
	}

	return &BlkStore{conf, fileMgr}, nil
}

// AddBlock adds a new block
func (store *BlkStore) AddBlock(block *cb.Block) error {
	result := store.fileMgr.addBlock(block)
	return result
}
