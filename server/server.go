package server

import (
	"io"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/dyfromnil/pdag/chain"
	"github.com/dyfromnil/pdag/chain/blkstorage"
	"github.com/dyfromnil/pdag/chain/blockledger"
	"github.com/dyfromnil/pdag/consensus"

	// "github.com/dyfromnil/pdag/consensus/pbft"
	hotstuffnet "github.com/dyfromnil/pdag/consensus/hotstuff-net"
	"github.com/dyfromnil/pdag/globleconfig"
	"github.com/dyfromnil/pdag/msp"
	cb "github.com/dyfromnil/pdag/proto-go/common"
)

//Main for server starting
func Main() {
	log.Printf("Server starting...")

	// ----------check identity------------
	var idt *msp.Identity

	if len(os.Args) != 2 {
		log.Fatal("Input error")
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
	ledger := blockledger.NewLedger(blkstore)

	chainSupport := chain.NewSupport(ledger, idt)

	//--------- consensus start ---------
	hotStuffServer := hotstuffnet.NewServer(chainSupport)
	hotStuffServer.Start()

	//------------ Leader Node listen Envelopes from clients to envCh -----------
	if nodeID == globleconfig.LeaderNodeID {
		listenEnv := grpc.NewServer()
		sendEnvelopsService := &SendEnvelopsService{
			ch: hotStuffServer.HandleChain(chainSupport),
		}
		cb.RegisterSendEnvelopsServer(listenEnv, sendEnvelopsService)

		lis, err := net.Listen("tcp", globleconfig.LeaderListenEnvelopeAddr)
		if err != nil {
			log.Fatalf("net.Listen err: %v", err)
		}

		go listenEnv.Serve(lis)
		log.Printf("receive pipe started!")
	}

	select {}
}

//SendEnvelopsService for
type SendEnvelopsService struct {
	ch consensus.Chain
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

		// log.Printf("Receive Envelope from %s,sendtime:%d", r.GetClientAddr(), r.GetTimestamp())

		s.ch.Order(r)

		err = stream.Send(&cb.Response{ResCode: true})
		if err != nil {
			return err
		}
	}
}
