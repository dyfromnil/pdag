package solo

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/dyfromnil/pdag/consensus"
	cb "github.com/dyfromnil/pdag/proto-go/common"

	"google.golang.org/grpc"

	cb "github.com/dyfromnil/pdag/proto-go/common"
)

type consenter struct{}

type chain struct {
	support  consensus.ConsenterSupport
	sendChan chan *message
	exitChan chan struct{}
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
		support:  support,
		sendChan: make(chan *message),
		exitChan: make(chan struct{}),
	}
}

func (ch *chain) Start() {
	go ch.main()
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
		return fmt.Errorf("Exiting")
	}
}

// Errored only closes on exit
func (ch *chain) Errored() <-chan struct{} {
	return ch.exitChan
}

func (ch *chain) main() {
	batchCh := make(chan []*cb.Envelope, 200)

	//make batch and send it to batch chan
	go func() {
		for {
			msg := <-ch.sendChan
			batches, _ := ch.support.BlockCutter().Ordered(msg.normalMsg)

			for _, batch := range batches {
				batchCh <- batch
			}
		}
	}()

	go createWorkerPool(1, batchCh, ch)
}

func worker(wg *sync.WaitGroup, batchCh chan []*cb.Envelope, ch *chain) {
	defer wg.Done()
	for batch := range batchCh {
		block, tipsList := ch.support.CreateNextBlock(batch)
		fmt.Println("num of tips:", len(tipsList))
		time.Sleep(time.Millisecond * time.Duration(500))
		ch.support.Append(block, tipsList)
	}
}
func createWorkerPool(numOfWorkers int, batchCh chan []*cb.Envelope, ch *chain) {
	var wg sync.WaitGroup
	for i := 0; i < numOfWorkers; i++ {
		wg.Add(1)
		go worker(&wg, batchCh, ch)
	}
	wg.Wait()
}
