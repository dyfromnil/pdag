package blkstorage

import (
	cb "github.com/dyfromnil/pdag/proto-go/common"
)

// BlockStore - filesystem based implementation for `BlockStore`
type BlockStore struct {
	conf    *Conf
	fileMgr *blockfileMgr
}

// newBlockStore constructs a `BlockStore`
func newBlockStore(id string, conf *Conf) (*BlockStore, error) {
	fileMgr, err := newBlockfileMgr(id, conf)
	if err != nil {
		return nil, err
	}

	return &BlockStore{conf, fileMgr}, nil
}

// AddBlock adds a new block
func (store *BlockStore) AddBlock(block *cb.Block, tips []*cb.Block) error {
	result := store.fileMgr.addBlock(block, tips)
	return result
}

//TipsBlock for
func (store *BlockStore) TipsBlock() []*cb.Block {
	return store.fileMgr.TipsBlock()
}
