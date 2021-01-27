package node

import (
	"fmt"
	"github.com/dyfromnil/pdag/globleconfig"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/dyfromnil/pdag/chain"
	"github.com/dyfromnil/pdag/chain/blkstorage"
	"github.com/dyfromnil/pdag/chain/blockledger/fileledger"
	"github.com/dyfromnil/pdag/consensus/solo"
	cb "github.com/dyfromnil/pdag/proto-go/common"
)

//SendEnvelopsService for
type SendEnvelopsService struct {
	envCh chan cb.Envelope
}

//Main for server starting
func Main() {
	fmt.Println("Server start...")

	server := grpc.NewServer()
	sendEnvelopsService := &SendEnvelopsService{envCh: make(chan cb.Envelope, 200)}
	cb.RegisterSendEnvelopsServer(server, sendEnvelopsService)

	lis, err := net.Listen("tcp", globleconfig.NodeTable["n0"])
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	go server.Serve(lis)

	conf := blkstorage.NewConf("", 0)
	blkstore, _ := blkstorage.NewBlockStore(conf)
	ledger := fileledger.NewFileLedger(blkstore)

	chainSupport := chain.NewSupport(ledger)

	soloConsensus := solo.New()
	soloChain := soloConsensus.HandleChain(chainSupport)

	var timer <-chan time.Time
	timer = time.After(time.Duration(time.Second * 10))

	soloChain.Start()

	fmt.Println("begining receive....")
	var i int = 0
	for {
		i++
		select {
		case env := <-sendEnvelopsService.envCh:
			err := soloChain.Order(&env)
			if err != nil {
				panic("order panic!")
			}
		case <-timer:
			break
		}
	}
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
