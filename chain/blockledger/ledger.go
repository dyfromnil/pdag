package blockledger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/dyfromnil/pdag/chain/blkstorage"
	"github.com/dyfromnil/pdag/globleconfig"
	cb "github.com/dyfromnil/pdag/proto-go/common"
)

// ReadWriter encapsulates the read/write functions of the ledger
type ReadWriter interface {
	CreateNextBlock(messages []*cb.Envelope, NumOfTransactionsInPool int) *cb.Block
	Append(block *cb.Block) error
	VerifyCurrentBlock(*cb.Block) bool
}

// Ledger is a struct used to interact with a node's ledger
type Ledger struct {
	blockStore           blkstorage.BlockStore
	round                int64                    //轮次
	roundTransactionNums int64                    //当前层打包交易总数
	roundBlockNums       map[int64]int            //每轮的区块总数
	roundPreRefNums      map[int64]int            //每轮剩余的总后继引用数
	roundDigestPost      map[int64]map[string]int //第几轮哪个区块的目前的后继区块数
	preBlocksDigest      map[string][]string      //当前区块的前驱区块digest,所有的string均为区块的digest
	digestToHash         map[string][]byte        //digest对应的hash值

	lock sync.Mutex

	roundCreateTransactionNums int64 //当前round创建时交易池中的交易数量

	lastEnvWaitingNum int     //上一时刻等待被打包的transaction数量
	lastDiff          float32 //上次的diff
	lastPreRef        int     //上个区块的前驱引用数

	fileOfRefNum        *os.File //记录索引数和交易池中的交易数量变化
	fileOfRoundBlockNum *os.File //记录当前层区块数量和层结束时与创建时交易总量的变化
}

// NewLedger creates a new Ledger for interaction with the ledger
func NewLedger(blockStore blkstorage.BlockStore) *Ledger {
	if err := os.MkdirAll("./log", 0755); err != nil {
		log.Fatalln("error while creating dir:'./log'")
	}
	fileOfRefNum, err := os.OpenFile("./log/refNum.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("error opening refNum log file writer for file refNum.csv")
	}
	fileOfRoundBlockNum, err := os.OpenFile("./log/roundBlockNum.csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("error opening roundBlockNum log file writer for file roundBlockNum.csv")
	}

	fl := &Ledger{
		blockStore:           blockStore,
		round:                1,
		roundTransactionNums: 0,
		roundPreRefNums:      make(map[int64]int),
		roundBlockNums:       make(map[int64]int),
		roundDigestPost:      make(map[int64]map[string]int),
		preBlocksDigest:      make(map[string][]string),
		digestToHash:         make(map[string][]byte),

		roundCreateTransactionNums: 0,

		lastEnvWaitingNum: 0,
		lastDiff:          1,
		lastPreRef:        1,

		fileOfRefNum:        fileOfRefNum,
		fileOfRoundBlockNum: fileOfRoundBlockNum,
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
func (fl *Ledger) CreateNextBlock(messages []*cb.Envelope, NumOfTransactionsInPool int) *cb.Block {
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

	preRefNum := fl.computePreRefNum(NumOfTransactionsInPool)
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

	fl.updateLedgerInfo(digest, preDigest, len(messages), header, isFull, NumOfTransactionsInPool)

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

func (fl *Ledger) updateLedgerInfo(digest string, preDigest []string, messagesLen int, header *cb.BlockHeader, isFull bool, NumOfTransactionsInPool int) {
	fl.preBlocksDigest[digest] = preDigest
	fl.roundTransactionNums += int64(messagesLen)
	fl.roundBlockNums[fl.round]++
	fl.roundPreRefNums[fl.round] += globleconfig.PostReference
	fl.setRoundDigestPost(fl.round, digest, globleconfig.PostReference)
	fl.digestToHash[digest] = blkstorage.BlockHeaderHash(header)

	if isFull {
		roundTransactionDiff := int64(NumOfTransactionsInPool) + fl.roundTransactionNums - fl.roundCreateTransactionNums
		log.Println("当前层打包交易数量: ", fl.roundTransactionNums)
		log.Println("当前round创建时交易池中的交易数量", fl.roundCreateTransactionNums)
		log.Println("第", fl.round, "层结束时与创建时交易总量的变化: ", roundTransactionDiff, "当前层平均交易速率（tps/区块）", roundTransactionDiff/int64(fl.roundBlockNums[fl.round]), "当前层的区块数量: ", fl.roundBlockNums[fl.round])
		fl.fileOfRoundBlockNum.WriteString(fmt.Sprintln(fl.round, "\t", roundTransactionDiff, "\t", roundTransactionDiff/int64(fl.roundBlockNums[fl.round]), "\t", fl.roundBlockNums[fl.round]))

		// log.Println("round ", fl.round, "'s num of Blocks: ", fl.roundBlockNums[fl.round])
		fl.roundTransactionNums = 0
		fl.round++
		fl.roundCreateTransactionNums = int64(NumOfTransactionsInPool)
	}
}

func (fl *Ledger) computePreRefNum(NumOfTransactionsInPool int) int {
	diff := float32(NumOfTransactionsInPool-fl.lastEnvWaitingNum)*globleconfig.Rate + (1-globleconfig.Rate)*float32(fl.lastDiff) //与上一区块打包时交易池交易量的变化值
	// diff := float32(NumOfTransactionsInPool - fl.lastEnvWaitingNum)

	var preRefNum int
	var scale float32
	if diff*fl.lastDiff > 0 {
		scale = abs(diff / fl.lastDiff)
	} else {
		scale = abs((diff - fl.lastDiff) / fl.lastDiff)
	}

	if diff >= 0 {
		preRefNum = int(float32(fl.lastPreRef) / scale)
		preRefNum = max(preRefNum, 1)
		preRefNum = min(preRefNum, globleconfig.PostReference)
		log.Println("diff>0")
	} else {
		preRefNum = int(float32(fl.lastPreRef) * scale)
		preRefNum = max(preRefNum, globleconfig.PostReference+1)
		preRefNum = min(preRefNum, globleconfig.PreReference)
		log.Println("diff<0")
	}

	fl.lastPreRef = preRefNum

	// log.Println("refNum (k) : ", refNum, " NumOfTransactionsInPool: ", NumOfTransactionsInPool, "lastNumOfTransactionsInPool: ", fl.lastEnvWaitingNum, " diff: ", diff, " lastDiff: ", fl.lastDiff, " scale: ", scale)
	// log.Println("refNum : ", refNum, " 当前交易池中待打包交易量: ", NumOfTransactionsInPool)
	fl.fileOfRefNum.WriteString(fmt.Sprintln(preRefNum, "\t", NumOfTransactionsInPool))

	// if fl.roundBlockNums[fl.round-1] >= globleconfig.NumOfConsensusGoroutine {
	if fl.roundBlockNums[fl.round-1] >= 10 {
		if preRefNum < globleconfig.PostReference {
			preRefNum = globleconfig.PostReference
		}
	}

	fl.lastEnvWaitingNum = NumOfTransactionsInPool
	fl.lastDiff = diff

	log.Println("preRefNum: ", preRefNum)
	return preRefNum
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
