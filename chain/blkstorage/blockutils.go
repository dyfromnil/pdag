package blkstorage

import (
	"bytes"
	"crypto/sha256"

	"github.com/golang/protobuf/proto"

	cb "github.com/dyfromnil/pdag/proto-go/common"
)

// NewBlock constructs a block with no data and no metadata.
func NewBlock(header *cb.BlockHeader, data *cb.BlockData) *cb.Block {
	return &cb.Block{
		Header: header,
		Data:   data,
	}
}

//BlockHeaderBytes for
func BlockHeaderBytes(b *cb.BlockHeader) []byte {
	result, err := proto.Marshal(b)
	if err != nil {
		panic(err)
	}
	return result
}

//BlocksHeaderHash for
func BlocksHeaderHash(bs []*cb.BlockHeader) [][]byte {
	blocksHeaderHash := [][]byte{}
	for _, bh := range bs {
		sum := sha256.Sum256(BlockHeaderBytes(bh))
		blocksHeaderHash = append(blocksHeaderHash, sum[:])
	}
	if len(blocksHeaderHash) == 0 {
		gensisBytes := sha256.Sum256([]byte("gensis block"))
		blocksHeaderHash = append(blocksHeaderHash, gensisBytes[:])
	}
	return blocksHeaderHash
}

//BlockDataHash for
func BlockDataHash(b *cb.BlockData) []byte {
	sum := sha256.Sum256(bytes.Join(b.Data, nil))
	return sum[:]
}
