package pbft

import (
	"context"
	"log"
	"sync"

	"github.com/dyfromnil/pdag/chain/blkstorage"
	"github.com/dyfromnil/pdag/consensus"
	// "github.com/dyfromnil/pdag/msp"
	cb "github.com/dyfromnil/pdag/proto-go/common"

	"google.golang.org/grpc"
)

type consenter struct{}

type chain struct {
	support   consensus.ConsenterSupport
	sendChan  chan *message
	exitChan  chan struct{}
	blockChan chan *cb.Block
	tipsList  []*cb.Block

	lock sync.Mutex
	//临时消息池，消息摘要对应消息本体
	blockPool map[string]*cb.Block
	//存放收到的prepare数量(至少需要收到并确认2f个)，根据摘要来对应
	prePareConfirmCount map[string]map[string]bool
	//存放收到的commit数量（至少需要收到并确认2f+1个），根据摘要来对应
	commitConfirmCount map[string]map[string]bool
	//该笔消息是否已进行Commit广播
	isCommitBordcast map[string]bool
	//该笔消息是否已对客户端进行Reply
	isReply map[string]bool
}

type message struct {
	normalMsg *cb.Envelope
}

//New for
func New() consensus.Consenter {
	return &consenter{}
}

func (pbft *consenter) HandleChain(support consensus.ConsenterSupport) consensus.Chain {
	return newChain(support)
}

func newChain(support consensus.ConsenterSupport) *chain {
	return &chain{
		support:   support,
		sendChan:  make(chan *message),
		exitChan:  make(chan struct{}),
		blockChan: make(chan *cb.Block),
		tipsList:  []*cb.Block{},

		blockPool:           make(map[string]*cb.Block),
		prePareConfirmCount: make(map[string]map[string]bool),
		commitConfirmCount:  make(map[string]map[string]bool),
		isCommitBordcast:    make(map[string]bool),
		isReply:             make(map[string]bool),
	}
}

func (ch *chain) Start() {
	go ch.createBlock()
}

func (ch *chain) Halt() {
	select {
	case <-ch.exitChan:
		// Allow multiple halts without panic
	default:
		close(ch.exitChan)
	}
}

func (ch *chain) WaitReady() error {
	return nil
}

// Order accepts normal messages for ordering
func (ch *chain) Order(env *cb.Envelope) error {
	select {
	case ch.sendChan <- &message{
		normalMsg: env,
	}:
		return nil
	case <-ch.exitChan:
		log.Fatal("Exiting")
		return nil
	}
}

// Errored only closes on exit
func (ch *chain) Errored() <-chan struct{} {
	return ch.exitChan
}

func (ch *chain) createBlock() {
	for {
		msg := <-ch.sendChan
		batches, _ := ch.support.BlockCutter().Ordered(msg.normalMsg)

		for _, batch := range batches {
			block, tipsList := ch.support.CreateNextBlock(batch)
			ch.tipsList = tipsList
			log.Printf("num of tips:", len(tipsList))
			ch.blockChan <- block
			// ch.support.Append(block, tipsList)
		}
	}
}

func (ch *chain) prePrepare() {
	for {
		//获取消息摘要
		block := <-ch.blockChan
		digest := blkstorage.BlockHeaderDigest(block.Header)
		log.Printf("已将block存入临时消息池")
		//存入临时消息池
		ch.blockPool[digest] = block
		//主节点对消息摘要进行签名
		digestByte, _ := hex.DecodeString(digest)
		signInfo := p.RsaSignWithSha256(digestByte, p.node.rsaPrivKey)
		//拼接成PrePrepare，准备发往follower节点
		pp := PrePrepare{*r, digest, p.sequenceID, signInfo}
		b, err := json.Marshal(pp)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("正在向其他节点进行进行PrePrepare广播 ...")
		//进行PrePrepare广播
		p.broadcast(cPrePrepare, b)
		fmt.Println("PrePrepare广播完成")
	}
}
