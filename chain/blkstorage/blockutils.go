package blkstorage

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"

	cb "github.com/dyfromnil/pdag/proto-go/common"
)

//GensisBlock for
func GensisBlock() *cb.Block {
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

//BlockHeaderDigest for
func BlockHeaderDigest(b *cb.BlockHeader) string {
	sum := sha256.Sum256(BlockHeaderBytes(b))
	return hex.EncodeToString(sum[:])
}

//BlocksHeaderHash for
func BlocksHeaderHash(bs []*cb.BlockHeader) [][]byte {
	blocksHeaderHash := [][]byte{}
	for _, bh := range bs {
		blocksHeaderHash = append(blocksHeaderHash, BlockHeaderBytes(bh))
	}
	return blocksHeaderHash
}

//BlockHeaderHash for
func BlockHeaderHash(header *cb.BlockHeader) []byte {
	sum := sha256.Sum256(BlockHeaderBytes(header))
	return sum[:]
}

//BlockDataHash for
func BlockDataHash(b *cb.BlockData) []byte {
	sum := sha256.Sum256(bytes.Join(b.Data, nil))
	return sum[:]
}
