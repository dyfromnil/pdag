package blkstorage

import (
	"bytes"
	"crypto/sha256"
	"encoding/asn1"
	cb "github.com/dyfromnil/pdag/proto-go/common"
)

// NewBlock constructs a block with no data and no metadata.
func NewBlock(previousHash [][]byte) *cb.Block {
	block := &cb.Block{}
	block.Header = &cb.BlockHeader{}
	block.Header.PreviousHash = previousHash
	block.Header.DataHash = []byte{}
	block.Data = &cb.BlockData{}
	return block
}

type asn1Header struct {
	PreviousHash [][]byte
	DataHash     []byte
}

//BlockHeaderBytes for
func BlockHeaderBytes(b *cb.BlockHeader) []byte {
	asn1Header := asn1Header{
		PreviousHash: b.PreviousHash,
		DataHash:     b.DataHash,
	}
	result, err := asn1.Marshal(asn1Header)
	if err != nil {
		// Errors should only arise for types which cannot be encoded, since the
		// BlockHeader type is known a-priori to contain only encodable types, an
		// error here is fatal and should not be propagated
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
