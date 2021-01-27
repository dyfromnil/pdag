package node

import (
	"io"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/dyfromnil/pdag/chain"
	"github.com/dyfromnil/pdag/chain/blkstorage"
	"github.com/dyfromnil/pdag/chain/blockledger/fileledger"
	"github.com/dyfromnil/pdag/consensus"
	"github.com/dyfromnil/pdag/consensus/solo"
	"github.com/dyfromnil/pdag/globleconfig"
	"github.com/dyfromnil/pdag/msp"
	cb "github.com/dyfromnil/pdag/proto-go/common"
)

//Main for server starting
func Main() {
	log.Printf("Server starting...")

	var idt msp.Identity
	idt.GenRsaKeys()

	if len(os.Args) != 2 {
		log.Panic("Input error")
	}
	nodeID := os.Args[1]
	if _, ok := globleconfig.NodeTable[nodeID]; ok {
		idt = msp.NewIdt(nodeID)
	} else {
		log.Fatal("无此节点编号！")
	}

	//----- config node : ledger and consensus -----
	conf := blkstorage.NewConf("", 0)
	blkstore, _ := blkstorage.NewBlockStore(conf)
	ledger := fileledger.NewFileLedger(blkstore)

	chainSupport := chain.NewSupport(ledger, &idt)

	soloConsensus := solo.New()
	soloChain := soloConsensus.HandleChain(chainSupport)

	node := Node{chain: soloChain}

	//--------- node start ---------
	node.startConsensus()
	envCh := make(chan cb.Envelope, 200)
	go node.startReceiveEnvelope(envCh)

	//------------ listen Envelopes from clients to envCh -----------
	listenEnv := grpc.NewServer()
	sendEnvelopsService := &SendEnvelopsService{
		envCh: envCh,
	}
	cb.RegisterSendEnvelopsServer(listenEnv, sendEnvelopsService)

	lis, err := net.Listen("tcp", globleconfig.LeaderListenEnvelopeAddr)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	go listenEnv.Serve(lis)

	select {}
}

//Node for
type Node struct {
	chain consensus.Chain
}

func (n *Node) startConsensus() {
	n.chain.Start()
}

func (n *Node) startReceiveEnvelope(envCh chan cb.Envelope) {
	log.Printf("receive pipe start!")
	var timer <-chan time.Time
	timer = time.After(time.Duration(time.Second * 10))
	for {
		select {
		case env := <-envCh:
			err := n.chain.Order(&env)
			if err != nil {
				panic("order panic!")
			}
		case <-timer:
			break
		}
	}
}

//SendEnvelopsService for
type SendEnvelopsService struct {
	envCh chan cb.Envelope
}

//Request for
func (s *SendEnvelopsService) Request(stream cb.SendEnvelops_RequestServer) error {
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("Receive Envelope from %s,sendtime:%d", r.GetClientAddr(), r.GetTimestamp())

		s.envCh <- *r

		err = stream.Send(&cb.Response{ResCode: true})
		if err != nil {
			return err
		}
	}
}
