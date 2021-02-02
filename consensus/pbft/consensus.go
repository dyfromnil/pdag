package pbft

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net"
	"sync"

	"github.com/dyfromnil/pdag/chain/blkstorage"
	"github.com/dyfromnil/pdag/client"
	"github.com/dyfromnil/pdag/consensus"
	"github.com/dyfromnil/pdag/globleconfig"
	cb "github.com/dyfromnil/pdag/proto-go/common"

	"google.golang.org/grpc"
)

// Server for
type Server struct {
	ch    *chain
	envCh chan cb.Envelope
}

type chain struct {
	support   consensus.ConsenterSupport
	sendChan  chan *message
	exitChan  chan struct{}
	blockChan chan *cb.Block

	lastEnvWaitingNum int //上一时刻等待被打包的transaction数量
	lastDiff          int //上次的diff
	lastPreRef        int //上个区块的前驱引用数

	lock sync.Mutex

	blockPool           map[string]*cb.Block       //临时消息池，消息摘要对应消息本体
	prePareConfirmCount map[string]map[string]bool //存放收到的prepare数量(至少需要收到并确认2f个)，根据摘要来对应
	commitConfirmCount  map[string]map[string]bool //存放收到的commit数量（至少需要收到并确认2f+1个），根据摘要来对应
	isCommitBordcast    map[string]bool            //该笔消息是否已进行Commit广播
	isReply             map[string]bool            //该笔消息是否已对客户端进行Reply
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
func (pb *Server) HandleChain(support consensus.ConsenterSupport) consensus.Chain {
	return pb.ch
}

// Start for
func (pb *Server) Start() {
	server := grpc.NewServer()
	cb.RegisterPbftServer(server, pb)

	lis, err := net.Listen("tcp", pb.ch.support.GetIdendity().GetSelfAddr())
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	go server.Serve(lis)
	pb.ch.Start()
	if pb.ch.support.GetIdendity().GetNodeID() == globleconfig.LeaderNodeID {
		go pb.ch.prePrepare()
	}
}

func newChain(support consensus.ConsenterSupport) *chain {
	return &chain{
		support:   support,
		sendChan:  make(chan *message, 20000),
		exitChan:  make(chan struct{}),
		blockChan: make(chan *cb.Block, 200),

		lastEnvWaitingNum: 0,
		lastDiff:          0,
		lastPreRef:        1,

		blockPool:           make(map[string]*cb.Block),
		prePareConfirmCount: make(map[string]map[string]bool),
		commitConfirmCount:  make(map[string]map[string]bool),
		isCommitBordcast:    make(map[string]bool),
		isReply:             make(map[string]bool),
	}
}

func (ch *chain) Start() {
	go ch.createBlock(0.7)
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

func (ch *chain) createBlock(rate float32) {
	for {
		msg := <-ch.sendChan
		batches, _ := ch.support.BlockCutter().Ordered(msg.normalMsg)

		for _, batch := range batches {
			diff := float32(len(ch.sendChan)-ch.lastEnvWaitingNum)*rate + (1-rate)*float32(ch.lastDiff)

			var refNum int
			scale := abs(diff / float32(ch.lastDiff))
			if diff > 0 {
				refNum = min(int(ch.lastPreRef*int(1+scale)), 10)
			} else {
				refNum = max(int(ch.lastPreRef*int(1-scale)), 1)
			}

			ch.lastPreRef = refNum

			block := ch.support.CreateNextBlock(batch, refNum)
			ch.blockChan <- block

			ch.lastEnvWaitingNum = len(ch.sendChan)
			ch.lastDiff = int(diff)
		}
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(a float32) float32 {
	if a < 0 {
		a = -a
	}
	return a
}

func (ch *chain) prePrepare() {
	for {
		//获取消息摘要
		block := <-ch.blockChan
		digest := blkstorage.BlockHeaderDigest(block.Header)
		log.Printf("已将block存入临时消息池")
		//存入临时消息池
		ch.lock.Lock()
		ch.blockPool[digest] = block
		ch.lock.Unlock()
		//主节点对消息摘要进行签名
		digestByte, _ := hex.DecodeString(digest)
		signInfo := ch.support.GetIdendity().RsaSignWithSha256(digestByte, ch.support.GetIdendity().GetSelfPivKey())
		//拼接成PrePrepare，准备发往follower节点
		pp := &cb.PrePrepareMsg{
			Block:  block,
			Digest: digest,
			Sign:   signInfo,
		}

		log.Printf("正在向其他节点进行进行PrePrepare广播 ...")
		//进行PrePrepare广播
		ch.broadcastPrePrepare(pp)
		log.Printf("PrePrepare广播完成")
	}
}

func (ch *chain) broadcastPrePrepare(pp *cb.PrePrepareMsg) {
	for _, i := range ch.support.GetIdendity().GetClusterAddrs() {
		if i == ch.support.GetIdendity().GetSelfAddr() {
			continue
		}
		go func(i string) {
			conn, err := grpc.Dial(i, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("grpc.Dial err: %v", err)
			}
			defer conn.Close()

			client := cb.NewPbftClient(conn)
			_, err = client.HandlePrePrepare(context.Background(), pp)
			// resp, err := client.HandlePrePrepare(context.Background(), pp)
			if err != nil {
				log.Fatalf("client.Search err: %v", err)
			}
			// log.Printf("resp: %t", resp.GetResCode())
		}(i)
	}
}

func (ch *chain) broadcastPrepare(p *cb.PrepareMsg) {
	for _, i := range ch.support.GetIdendity().GetClusterAddrs() {
		if i == ch.support.GetIdendity().GetSelfAddr() {
			continue
		}
		go func(i string) {
			conn, err := grpc.Dial(i, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("grpc.Dial err: %v", err)
			}
			defer conn.Close()

			client := cb.NewPbftClient(conn)
			_, err = client.HandlePrepare(context.Background(), p)
			if err != nil {
				log.Fatalf("client.Search err: %v", err)
			}
			// log.Printf("resp: %t", resp.GetResCode())
		}(i)
	}
}

func (ch *chain) broadcastCommit(c *cb.CommitMsg) {
	for _, i := range ch.support.GetIdendity().GetClusterAddrs() {
		if i == ch.support.GetIdendity().GetSelfAddr() {
			continue
		}
		go func(i string) {
			conn, err := grpc.Dial(i, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("grpc.Dial err: %v", err)
			}
			defer conn.Close()

			client := cb.NewPbftClient(conn)
			_, err = client.HandleCommit(context.Background(), c)
			if err != nil {
				log.Fatalf("client.Search err: %v", err)
			}
			// log.Printf("resp: %t", resp.GetResCode())
		}(i)
	}
}

//为多重映射开辟赋值
func (ch *chain) setPrePareConfirmMap(val, val2 string, b bool) {
	if _, ok := ch.prePareConfirmCount[val]; !ok {
		ch.prePareConfirmCount[val] = make(map[string]bool)
	}
	ch.prePareConfirmCount[val][val2] = b
}

//为多重映射开辟赋值
func (ch *chain) setCommitConfirmMap(val, val2 string, b bool) {
	if _, ok := ch.commitConfirmCount[val]; !ok {
		ch.commitConfirmCount[val] = make(map[string]bool)
	}
	ch.commitConfirmCount[val][val2] = b
}

//HandlePrePrepare for
func (pb *Server) HandlePrePrepare(ctx context.Context, pp *cb.PrePrepareMsg) (*cb.Response, error) {
	pb.ch.lock.Lock()
	defer pb.ch.lock.Unlock()
	log.Printf("本节点已接收到主节点发来的PrePrepare ...")

	if !pb.ch.support.VerifyCurrentBlock(pp.Block) {
		log.Fatalln("该block与本地账本存在冲突！！！")
	}

	//获取主节点的公钥，用于数字签名验证
	primaryNodePubKey := pb.ch.support.GetIdendity().GetPubKey(globleconfig.LeaderNodeID)
	digestByte, _ := hex.DecodeString(pp.Digest)
	if digest := blkstorage.BlockHeaderDigest(pp.Block.Header); digest != pp.Digest {
		log.Fatalln("信息摘要对不上，拒绝进行prepare广播")
		return &cb.Response{ResCode: false}, errors.New("HandlePrePrepare digest error")
	} else if !pb.ch.support.GetIdendity().RsaVerySignWithSha256(digestByte, pp.Sign, primaryNodePubKey) {
		log.Printf("主节点签名验证失败！,拒绝进行prepare广播")
		return &cb.Response{ResCode: false}, errors.New("HandlePrePrepare VerySign error")
	} else {
		//将信息存入临时消息池
		log.Printf("已将消息存入临时节点池")
		// pb.ch.lock.Lock()
		pb.ch.blockPool[pp.Digest] = pp.Block
		log.Println("Pool add digest:", digest)

		// pb.ch.lock.Unlock()
		//节点使用私钥对其签名
		sign := pb.ch.support.GetIdendity().RsaSignWithSha256(digestByte, pb.ch.support.GetIdendity().GetSelfPivKey())
		//拼接成Prepare
		pre := &cb.PrepareMsg{
			Digest: pp.Digest,
			NodeID: pb.ch.support.GetIdendity().GetNodeID(),
			Sign:   sign,
		}

		//进行准备阶段的广播
		log.Printf("正在进行Prepare广播 ...")
		pb.ch.broadcastPrepare(pre)
		log.Printf("Prepare广播完成")

		// 接收到Preprepare后发现早已收到足够多后一阶段prepare消息，^_^!节点间延迟太高啦！ @_@几个阶段杂揉了，有时间再改!_!
		if pb.canPrepared(digest) {
			log.Println("本节点已收到至少2f个节点(包括本地节点)发来的Prepare信息 ...")
			//节点使用私钥对其签名
			sign := pb.ch.support.GetIdendity().RsaSignWithSha256(digestByte, pb.ch.support.GetIdendity().GetSelfPivKey())
			c := &cb.CommitMsg{
				Digest: digest,
				NodeID: pb.ch.support.GetIdendity().GetNodeID(),
				Sign:   sign,
			}

			//进行提交信息的广播
			log.Println("正在进行commit广播", digest)
			pb.ch.broadcastCommit(c)
			pb.ch.isCommitBordcast[digest] = true
			log.Printf("commit广播完成")

			if pb.canCommit(digest) {
				log.Printf("本节点已收到至少2f + 1 个节点(包括本地节点)发来的Commit信息 ...")
				//将消息信息，提交到本地消息池中！
				pb.ch.support.Append(pb.ch.blockPool[c.Digest])
				log.Println("Delete digest from pool: ", c.Digest)
				log.Println(pb.ch.support.GetIdendity().GetNodeID() + "节点已将当前block存入本地账本中")
				log.Printf("正在reply客户端 ...")
				reply := client.Msg{
					Digest:     c.Digest,
					NumOfTranc: len(pb.ch.blockPool[c.Digest].Data.Data),
					NodeID:     pb.ch.support.GetIdendity().GetNodeID(),
				}
				replyBytes, err := json.Marshal(reply)
				if err != nil {
					log.Fatal("Marshal Error!")
				}
				tcpDial(replyBytes, globleconfig.ClientAddr)
				pb.ch.isReply[c.Digest] = true
				log.Printf("reply完毕")
			}
		}
	}
	return &cb.Response{ResCode: true}, nil
}

// HandlePrepare for
func (pb *Server) HandlePrepare(ctx context.Context, p *cb.PrepareMsg) (*cb.Response, error) {
	pb.ch.lock.Lock()
	defer pb.ch.lock.Unlock()
	log.Printf("本节点已接收到%s节点发来的Prepare ...", p.NodeID)
	//获取消息源节点的公钥，用于数字签名验证
	MessageNodePubKey := pb.ch.support.GetIdendity().GetPubKey(p.NodeID)
	digestByte, _ := hex.DecodeString(p.Digest)
	// pb.ch.lock.Lock()
	_, ok := pb.ch.blockPool[p.Digest]
	// pb.ch.lock.Unlock()
	if !ok {
		log.Println("当前临时消息池无此摘要，等待Preprepare消息&暂存当前Prepare消息至prepare池")
		pb.ch.setPrePareConfirmMap(p.Digest, p.NodeID, true)
		return &cb.Response{ResCode: false}, nil
	} else if !pb.ch.support.GetIdendity().RsaVerySignWithSha256(digestByte, p.Sign, MessageNodePubKey) {
		log.Printf("节点签名验证失败！,拒绝执行commit广播")
		return &cb.Response{ResCode: false}, errors.New("HandlePrepare VerySign error")
	} else {
		// pb.ch.lock.Lock()
		pb.ch.setPrePareConfirmMap(p.Digest, p.NodeID, true)
		//如果节点至少收到了2f个prepare的消息（包括自己）,并且没有进行过commit广播，则进行commit广播
		// pb.ch.lock.Lock()
		//获取消息源节点的公钥，用于数字签名验证
		if pb.canPrepared(p.Digest) {
			log.Println("本节点已收到至少2f个节点(包括本地节点)发来的Prepare信息 ...", p.Digest)
			//节点使用私钥对其签名
			sign := pb.ch.support.GetIdendity().RsaSignWithSha256(digestByte, pb.ch.support.GetIdendity().GetSelfPivKey())
			c := &cb.CommitMsg{
				Digest: p.Digest,
				NodeID: pb.ch.support.GetIdendity().GetNodeID(),
				Sign:   sign,
			}

			//进行提交信息的广播
			log.Println("正在进行commit广播", p.Digest)
			pb.ch.broadcastCommit(c)
			pb.ch.isCommitBordcast[p.Digest] = true
			log.Printf("commit广播完成")

			// 接收到prepare后发现早已收到足够多后一阶段commit消息，^_^!节点间延迟太高啦！ @_@几个阶段杂揉了，有时间再改!_!
			if pb.canCommit(p.Digest) {
				log.Printf("本节点已收到至少2f + 1 个节点(包括本地节点)发来的Commit信息 ...")
				//将消息信息，提交到本地消息池中！
				pb.ch.support.Append(pb.ch.blockPool[p.Digest])
				log.Println("Delete digest from pool: ", p.Digest)
				log.Println(pb.ch.support.GetIdendity().GetNodeID() + "节点已将当前block存入本地账本中")
				log.Printf("正在reply客户端 ...")
				reply := client.Msg{
					Digest:     c.Digest,
					NumOfTranc: len(pb.ch.blockPool[c.Digest].Data.Data),
					NodeID:     pb.ch.support.GetIdendity().GetNodeID(),
				}
				replyBytes, err := json.Marshal(reply)
				if err != nil {
					log.Fatal("Marshal Error!")
				}
				tcpDial(replyBytes, globleconfig.ClientAddr)
				pb.ch.isReply[p.Digest] = true
				log.Printf("reply完毕")
			}
		}
		// pb.ch.lock.Unlock()
	}
	return &cb.Response{ResCode: true}, nil
}

func (pb *Server) canPrepared(digest string) bool {
	count := len(pb.ch.prePareConfirmCount[digest])
	// pb.ch.lock.Unlock()
	//因为主节点不会发送Prepare，所以不包含自己
	specifiedCount := 0
	nodeCount := len(pb.ch.support.GetIdendity().GetClusterAddrs())
	if pb.ch.support.GetIdendity().GetNodeID() == globleconfig.LeaderNodeID {
		specifiedCount = nodeCount / 3 * 2
	} else {
		specifiedCount = (nodeCount / 3 * 2) - 1
	}
	return count >= specifiedCount && !pb.ch.isCommitBordcast[digest]
}

//HandleCommit for
func (pb *Server) HandleCommit(ctx context.Context, c *cb.CommitMsg) (*cb.Response, error) {
	pb.ch.lock.Lock()
	defer pb.ch.lock.Unlock()
	log.Printf("本节点已接收到%s节点发来的Commit ... ", c.NodeID)
	//获取消息源节点的公钥，用于数字签名验证
	MessageNodePubKey := pb.ch.support.GetIdendity().GetPubKey(c.NodeID)
	digestByte, _ := hex.DecodeString(c.Digest)
	// pb.ch.lock.Lock()
	_, ok := pb.ch.prePareConfirmCount[c.Digest]
	// pb.ch.lock.Unlock()
	if !ok {
		pb.ch.setCommitConfirmMap(c.Digest, c.NodeID, true)
		log.Println("当前prepare池无此摘要，等待prepare消息&暂存当前commit消息至CommitConfirm池")
		return &cb.Response{ResCode: false}, nil
	} else if !pb.ch.support.GetIdendity().RsaVerySignWithSha256(digestByte, c.Sign, MessageNodePubKey) {
		log.Printf("节点签名验证失败！,拒绝将信息持久化到本地消息池")
		return &cb.Response{ResCode: false}, errors.New("HandleCommit VerySign error")
	} else {
		// pb.ch.lock.Lock()
		pb.ch.setCommitConfirmMap(c.Digest, c.NodeID, true)

		if pb.canCommit(c.Digest) {
			log.Printf("本节点已收到至少2f + 1 个节点(包括本地节点)发来的Commit信息 ...")
			//将消息信息，提交到本地消息池中！
			pb.ch.support.Append(pb.ch.blockPool[c.Digest])
			log.Println("Delete digest from pool: ", c.Digest)
			log.Println(pb.ch.support.GetIdendity().GetNodeID() + "节点已将当前block存入本地账本中")
			log.Printf("正在reply客户端 ...")
			reply := client.Msg{
				Digest:     c.Digest,
				NumOfTranc: len(pb.ch.blockPool[c.Digest].Data.Data),
				NodeID:     pb.ch.support.GetIdendity().GetNodeID(),
			}
			replyBytes, err := json.Marshal(reply)
			if err != nil {
				log.Fatal("Marshal Error!")
			}
			tcpDial(replyBytes, globleconfig.ClientAddr)
			pb.ch.isReply[c.Digest] = true
			log.Printf("reply完毕")
		}
		// pb.ch.lock.Unlock()
	}
	return &cb.Response{ResCode: true}, nil
}

func (pb *Server) canCommit(digest string) bool {
	count := len(pb.ch.commitConfirmCount[digest])
	//如果节点至少收到了2f+1个commit消息（包括自己）,并且节点没有回复过,并且已进行过commit广播，则提交信息至本地消息池，并reply成功标志至客户端！
	nodeCount := len(pb.ch.support.GetIdendity().GetClusterAddrs())
	return count >= nodeCount/3*2 && !pb.ch.isReply[digest] && pb.ch.isCommitBordcast[digest]
}

//使用tcp发送消息
func tcpDial(context []byte, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("connect error", err)
		return
	}

	_, err = conn.Write(context)
	if err != nil {
		log.Fatal(err)
	}
	conn.Close()
}
