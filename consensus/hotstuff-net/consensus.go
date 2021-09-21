package hotstuffnet

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/dyfromnil/pdag/chain/blkstorage"
	"github.com/dyfromnil/pdag/client"
	"github.com/dyfromnil/pdag/consensus"
	"github.com/dyfromnil/pdag/globleconfig"
	cb "github.com/dyfromnil/pdag/proto-go/common"

	"google.golang.org/grpc"
)

// Server for
type Server struct {
	ch *chain
	// envCh chan cb.Envelope
}

type chain struct {
	support  consensus.ConsenterSupport
	sendChan chan *message
	exitChan chan struct{}

	lock  sync.Mutex
	cntCh chan bool

	createTime           map[string]int64 //区块创建时间
	numsOfCommitedBlocks int              //已commit的区块数量
	totalDelay           int64            //已commit的区块的总确认延迟
	fileOfDelay          *os.File         //记录区块确认延迟

	blockPool          map[string]*cb.Block       //临时消息池，消息摘要对应消息本体
	commitConfirmCount map[string]map[string]bool //存放收到的commit数量（至少需要收到并确认2f+1个），根据摘要来对应
	isReply            map[string]bool            //该笔消息是否已replay client
}

type message struct {
	normalMsg *cb.Envelope
}

//NewServer for
func NewServer(support consensus.ConsenterSupport) consensus.Consenter {
	return &Server{
		ch: newChain(support),
	}
}

//HandleChain for
func (hs *Server) HandleChain(support consensus.ConsenterSupport) consensus.Chain {
	return hs.ch
}

// Start for
func (hs *Server) Start() {
	server := grpc.NewServer()
	cb.RegisterHotStuffServer(server, hs)

	lis, err := net.Listen("tcp", hs.ch.support.GetIdendity().GetSelfAddr())
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	go server.Serve(lis)
	hs.ch.Start()
	if hs.ch.support.GetIdendity().GetNodeID() == globleconfig.LeaderNodeID {
		go hs.ch.propose()
	}
}

func newChain(support consensus.ConsenterSupport) *chain {
	if err := os.MkdirAll("./log", 0755); err != nil {
		log.Fatalln("error while creating dir:'./log'")
	}
	fileOfDelay, err := os.OpenFile("./log/delay.log", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal("error opening delay log file writer for file delay.log")
	}
	return &chain{
		support:  support,
		sendChan: make(chan *message, 200000),
		exitChan: make(chan struct{}),
		cntCh:    make(chan bool, globleconfig.NumOfConsensusGoroutine), //num of hotstuff goroutines

		createTime:           make(map[string]int64),
		numsOfCommitedBlocks: 0,
		totalDelay:           0,
		fileOfDelay:          fileOfDelay,

		blockPool:          make(map[string]*cb.Block),
		commitConfirmCount: make(map[string]map[string]bool),
		isReply:            make(map[string]bool),
	}
}

func (ch *chain) Start() {
	// go ch.createBlock(0.7)
	go ch.statDelay()
}

func (ch *chain) statDelay() {
	ticker := time.NewTicker(time.Duration(time.Second * 5))
	for {
		<-ticker.C
		ch.fileOfDelay.WriteString(fmt.Sprintln("numsOfCommitedBlocks: ", ch.numsOfCommitedBlocks, "totalDelay: ", float64(ch.totalDelay)/1e6))
	}
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
		//log.Fatal("Exiting")
		return nil
	}
}

// Errored only closes on exit
func (ch *chain) Errored() <-chan struct{} {
	return ch.exitChan
}

func (ch *chain) propose() {
	for {
		msg := <-ch.sendChan
		log.Println("???????????????????????")
		batches, _ := ch.support.BlockCutter().Ordered(msg.normalMsg)

		for _, batch := range batches {
			ch.cntCh <- true

			NumOfTransactionsInPool := len(ch.sendChan) //当前区块打包时，交易池交易量
			block := ch.support.CreateNextBlock(batch, NumOfTransactionsInPool)
			//获取消息摘要
			digest := blkstorage.BlockHeaderDigest(block.Header)
			//log.Printf("已将block存入临时消息池")
			//存入临时消息池
			ch.lock.Lock()
			ch.blockPool[digest] = block
			ch.createTime[digest] = time.Now().UnixNano()
			ch.lock.Unlock()
			//主节点对消息摘要进行签名
			digestByte, _ := hex.DecodeString(digest)
			signInfo := ch.support.GetIdendity().RsaSignWithSha256(digestByte, ch.support.GetIdendity().GetSelfPivKey())
			//拼接成PrePrepare，准备发往follower节点
			proposal := &cb.ProposeMsg{
				Block:  block,
				Digest: digest,
				Sign:   signInfo,
			}

			log.Printf("正在向其他节点进行进行Porposal广播 ...")
			ch.broadcastProposal(proposal)
		}
	}
}

func (ch *chain) broadcastProposal(proposal *cb.ProposeMsg) {
	for _, addr := range ch.support.GetIdendity().GetClusterAddrs() {
		go func(addr string) {
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("grpc.Dial err: %v", err)
			}
			defer conn.Close()

			client := cb.NewHotStuffClient(conn)
			_, err = client.OnReceiveProposal(context.Background(), proposal)
			// resp, err := client.HandlePrePrepare(context.Background(), pp)
			if err != nil {
				log.Fatalf("client.Search err: %v", err)
			}
			// //log.Printf("resp: %t", resp.GetResCode())
		}(addr)
	}
}

//HandlePrePrepare for
func (hs *Server) OnReceiveProposal(ctx context.Context, proposeMsg *cb.ProposeMsg) (*cb.Response, error) {
	log.Println("---------------------------------------------------------????????????????--------------------------")
	hs.ch.lock.Lock()
	defer hs.ch.lock.Unlock()
	//log.Printf("本节点已接收到主节点发来的PrePrepare ...")

	if !hs.ch.support.VerifyCurrentBlock(proposeMsg.Block) {
		log.Fatalln("该block与本地账本存在冲突！！！")
	}

	//获取主节点的公钥，用于数字签名验证
	primaryNodePubKey := hs.ch.support.GetIdendity().GetPubKey(globleconfig.LeaderNodeID)
	digestByte, _ := hex.DecodeString(proposeMsg.Digest)
	if digest := blkstorage.BlockHeaderDigest(proposeMsg.Block.Header); digest != proposeMsg.Digest {
		//log.Fatalln("信息摘要对不上，拒绝进行prepare广播")
		return &cb.Response{ResCode: false}, errors.New("HandlePrePrepare digest error")
	} else if !hs.ch.support.GetIdendity().RsaVerySignWithSha256(digestByte, proposeMsg.Sign, primaryNodePubKey) {
		//log.Printf("主节点签名验证失败！,拒绝进行prepare广播")
		return &cb.Response{ResCode: false}, errors.New("HandlePrePrepare VerySign error")
	} else {
		//将信息存入临时消息池
		//log.Printf("已将消息存入临时节点池")
		// pb.ch.lock.Lock()
		hs.ch.blockPool[proposeMsg.Digest] = proposeMsg.Block
		//log.Println("Pool add digest:", digest)

		// pb.ch.lock.Unlock()
		//节点使用私钥对其签名
		sign := hs.ch.support.GetIdendity().RsaSignWithSha256(digestByte, hs.ch.support.GetIdendity().GetSelfPivKey())
		//拼接成Prepare
		voteMsg := &cb.VoteMsg{
			Digest: proposeMsg.Digest,
			NodeID: hs.ch.support.GetIdendity().GetNodeID(),
			Sign:   sign,
		}

		go hs.ch.Vote(voteMsg)
		go hs.ch.update(proposeMsg.Block, digest)
	}
	return &cb.Response{ResCode: true}, nil
}

func (hs *Server) OnReceiveVote(ctx context.Context, voteMsg *cb.VoteMsg) (*cb.Response, error) {
	hs.ch.lock.Lock()
	defer hs.ch.lock.Unlock()
	//log.Printf("本节点已接收到主节点发来的PrePrepare ...")

	MessageNodePubKey := hs.ch.support.GetIdendity().GetPubKey(voteMsg.NodeID)
	digestByte, _ := hex.DecodeString(voteMsg.Digest)
	_, ok := hs.ch.blockPool[voteMsg.Digest]

	if !ok {
		return &cb.Response{ResCode: false}, nil
	} else if !hs.ch.support.GetIdendity().RsaVerySignWithSha256(digestByte, voteMsg.Sign, MessageNodePubKey) {
		log.Printf("节点签名验证失败！,拒绝执行commit广播")
		return &cb.Response{ResCode: false}, errors.New("HandlePrepare VerySign error")
	} else {
		hs.ch.setCommitConfirmMap(voteMsg.Digest, voteMsg.NodeID, true)

		if hs.meetThreshold(voteMsg.Digest) {
			//将消息信息，提交到本地消息池中！
			hs.ch.support.Append(hs.ch.blockPool[voteMsg.Digest])
			hs.ch.numsOfCommitedBlocks++
			hs.ch.totalDelay += time.Now().UnixNano() - hs.ch.createTime[voteMsg.Digest]
			if hs.ch.support.GetIdendity().GetNodeID() == globleconfig.LeaderNodeID {
				<-hs.ch.cntCh
			}

			reply := client.Msg{
				Digest:     voteMsg.Digest,
				NumOfTranc: len(hs.ch.blockPool[voteMsg.Digest].Data.Data),
				NodeID:     hs.ch.support.GetIdendity().GetNodeID(),
			}
			replyBytes, err := json.Marshal(reply)
			if err != nil {
				log.Fatal("Marshal Error!")
			}
			go tcpDial(replyBytes, globleconfig.ClientAddr)
			hs.ch.isReply[voteMsg.Digest] = true

			log.Println("-----replay------success----------------------------------")

		}
	}
	return &cb.Response{ResCode: true}, nil
}

func (ch *chain) update(block *cb.Block, digest string) {
	ch.updateQcHigh("qc")
	ch.onCommit(block, digest)
}

func (ch *chain) updateQcHigh(qc string) {
}

func (ch *chain) onCommit(block *cb.Block, digest string) {
	ch.support.Append(ch.blockPool[digest])
}

func (ch *chain) Vote(voteMsg *cb.VoteMsg) {
	leaderAddr := ch.support.GetIdendity().GetLeaderAddr()
	conn, err := grpc.Dial(leaderAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	client := cb.NewHotStuffClient(conn)
	_, err = client.OnReceiveVote(context.Background(), voteMsg)
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}
}

//为多重映射开辟赋值
func (ch *chain) setCommitConfirmMap(val, val2 string, b bool) {
	if _, ok := ch.commitConfirmCount[val]; !ok {
		ch.commitConfirmCount[val] = make(map[string]bool)
	}
	ch.commitConfirmCount[val][val2] = b
}

func (ht *Server) meetThreshold(digest string) bool {
	count := len(ht.ch.commitConfirmCount[digest])
	//如果节点至少收到了2f+1个commit消息（包括自己）,并且节点没有回复过,并且已进行过commit广播，则提交信息至本地消息池，并reply成功标志至客户端！
	nodeCount := len(ht.ch.support.GetIdendity().GetClusterAddrs())
	return count >= nodeCount/3*2 && !ht.ch.isReply[digest]
}

//使用tcp发送消息
func tcpDial(context []byte, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		//log.Println("connect error", err)
		return
	}

	_, err = conn.Write(context)
	if err != nil {
		log.Fatal(err)
	}
	conn.Close()
}
