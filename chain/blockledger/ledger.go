package blockledger

import (
	"log"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/dyfromnil/pdag/chain/blkstorage"
	"github.com/dyfromnil/pdag/globleconfig"
	cb "github.com/dyfromnil/pdag/proto-go/common"
)

// ReadWriter encapsulates the read/write functions of the ledger
type ReadWriter interface {
	CreateNextBlock(messages []*cb.Envelope, preRefNum int) *cb.Block
	Append(block *cb.Block) error
	VerifyCurrentBlock(*cb.Block) bool
}

// Ledger is a struct used to interact with a node's ledger
type Ledger struct {
	blockStore      blkstorage.BlockStore
	round           int64                    //轮次
	roundBlockNums  map[int64]int            //每轮的区块总数
	roundPreRefNums map[int64]int            //每轮剩余的总后继引用数
	roundDigestPost map[int64]map[string]int //第几轮哪个区块的目前的后继区块数
	preBlocksDigest map[string][]string      //当前区块的前驱区块digest,所有的string均为区块的digest
	digestToHash    map[string][]byte        //digest对应的hash值

	lock sync.Mutex
}

// NewLedger creates a new Ledger for interaction with the ledger
func NewLedger(blockStore blkstorage.BlockStore) *Ledger {
	fl := &Ledger{
		blockStore:      blockStore,
		round:           1,
		roundPreRefNums: make(map[int64]int),
		roundBlockNums:  make(map[int64]int),
		roundDigestPost: make(map[int64]map[string]int),
		preBlocksDigest: make(map[string][]string),
		digestToHash:    make(map[string][]byte),
	}
	gensisBlock := blkstorage.GensisBlock()
	fl.lock.Lock()
	defer fl.lock.Unlock()
	fl.blockStore.AddBlock(gensisBlock)
	digest := blkstorage.BlockHeaderDigest(gensisBlock.Header)
	log.Println("gensis digest: ", digest)
	fl.roundBlockNums[fl.round-1] = 1
	fl.roundPreRefNums[fl.round-1] = globleconfig.PostReference
	fl.setRoundDigestPost(0, digest, globleconfig.PostReference)
	fl.digestToHash[digest] = blkstorage.BlockHeaderHash(gensisBlock.Header)
	return fl
}

// Append a new block to the ledger
func (fl *Ledger) Append(block *cb.Block) error {
	fl.lock.Lock()
	defer fl.lock.Unlock()
	if err := fl.blockStore.AddBlock(block); err != nil {
		return err
	}

	bhd := blkstorage.BlockHeaderDigest(block.Header)
	for _, digest := range fl.preBlocksDigest[bhd] {
		fl.garbageCollect(digest)
	}

	return nil
}

// CreateNextBlock provides a utility way to construct the next block
func (fl *Ledger) CreateNextBlock(messages []*cb.Envelope, preRefNum int) *cb.Block {
	var err error
	blockData := &cb.BlockData{Data: make([][]byte, len(messages))}
	for i, msg := range messages {
		blockData.Data[i], err = proto.Marshal(msg)
		if err != nil {
			panic(err)
		}
	}

	fl.lock.Lock()
	defer fl.lock.Unlock()

	preDigest, isFull := fl.choosePreBlocksDigestAndIsFull(preRefNum)

	previousHash := [][]byte{}
	for _, digest := range preDigest {
		previousHash = append(previousHash, fl.digestToHash[digest])
	}

	header := &cb.BlockHeader{
		PreviousHash: previousHash,
		DataHash:     blkstorage.BlockDataHash(blockData),
		Round:        fl.round,
		Timestamp:    time.Now().UnixNano(),
	}

	block := blkstorage.NewBlock(header, blockData)

	digest := blkstorage.BlockHeaderDigest(header)
	fl.preBlocksDigest[digest] = preDigest
	fl.roundBlockNums[fl.round]++
	fl.roundPreRefNums[fl.round] += globleconfig.PostReference
	fl.setRoundDigestPost(fl.round, digest, globleconfig.PostReference)
	fl.digestToHash[digest] = blkstorage.BlockHeaderHash(header)

	if isFull {
		log.Println("round ", fl.round, "'s num of Blocks: ", fl.roundBlockNums[fl.round])
		fl.round++
	}

	return block
}

func (fl *Ledger) choosePreBlocksDigestAndIsFull(preRefNum int) ([]string, bool) {

	dgts := []string{}
	//从前一层找
	for digest, num := range fl.roundDigestPost[fl.round-1] {
		if num > 0 {
			// log.Println("<<<<<前一层:", fl.round-1, ">>>>> ", fl.roundDigestPost, "roundBlockNums: ", fl.roundBlockNums[fl.round-1], "this roundBlockNums: ", fl.roundBlockNums[fl.round], "前一层剩余引用数: ", fl.roundPreRefNums[fl.round-1])
			dgts = append(dgts, digest)
			fl.roundDigestPost[fl.round-1][digest]--
			fl.roundPreRefNums[fl.round-1]--
			preRefNum--
		}
		// log.Println("<<<<<前一层:", fl.round-1, ">>>>> ", fl.roundDigestPost, "roundBlockNums: ", fl.roundBlockNums[fl.round-1], "前一层剩余引用数: ", fl.roundPreRefNums[fl.round-1])

		if preRefNum == 0 {
			break
		}
	}
	isFull := false

	if fl.roundPreRefNums[fl.round-1] == 0 {
		isFull = true
	}

	//若前一层的所有后继均满，且当前区块的前驱引用数未达到最大，则从更前层寻找（更前层可能有共识失败后返还的后继索引计数）
	// if preRefNum > 0 {
	// 	for rd := range fl.roundDigestPost {
	// 		if rd != fl.round-1 {
	// 			for digest, num := range fl.roundDigestPost[rd] {
	// 				if num > 0 {
	// 					dgts = append(dgts, digest)
	// 					fl.roundDigestPost[rd][digest]--
	// 					preRefNum--
	// 				}
	// 				if preRefNum == 0 {
	// 					break
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	return dgts, isFull
}

//VerifyCurrentBlock for
func (fl *Ledger) VerifyCurrentBlock(block *cb.Block) bool {
	return true
}

//为多重映射开辟赋值
func (fl *Ledger) setRoundDigestPost(val int64, val2 string, b int) {
	if _, ok := fl.roundDigestPost[val]; !ok {
		fl.roundDigestPost[val] = make(map[string]int)
	}
	fl.roundDigestPost[val][val2] = b
}

func (fl *Ledger) garbageCollect(digest string) {
	for rd := range fl.roundDigestPost {
		if rd == fl.round {
			continue
		}
		if nums, ok := fl.roundDigestPost[rd][digest]; ok {
			if nums == 0 {
				delete(fl.digestToHash, digest)
				delete(fl.preBlocksDigest, digest)
				delete(fl.roundDigestPost[rd], digest)
				if len(fl.roundDigestPost[rd]) == 0 {
					delete(fl.roundDigestPost, rd)
				}
			}
		}
	}
}
