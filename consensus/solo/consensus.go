package solo

import (
	"fmt"
	"sync"
	"time"

	"math/rand"

	"github.com/dyfromnil/pdag/consensus"
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

// New creates a new consenter for the solo consensus scheme.
// The solo consensus scheme is very simple, and allows only one consenter for a given chain (this process).
// It accepts messages being delivered via Order/Configure, orders them, and then uses the blockcutter to form the messages
// into blocks before writing to the given ledger
func New() consensus.Consenter {
	return &consenter{}
}

func (solo *consenter) HandleChain(support consensus.ConsenterSupport) consensus.Chain {
	return newChain(support)
}

func (solo *consenter) Start() {
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
	// var timer <-chan time.Time

	// for {
	// 	select {
	// 	case msg := <-ch.sendChan:
	// 		batches, pending := ch.support.BlockCutter().Ordered(msg.normalMsg)

	// 		for _, batch := range batches {
	// 			block, tipsList := ch.support.CreateNextBlock(batch)
	// 			fmt.Println("num of tips:", len(tipsList))
	// 			time.Sleep(time.Second)
	// 			ch.support.Append(block, tipsList)
	// 		}

	// 		switch {
	// 		case timer != nil && !pending:
	// 			// Timer is already running but there are no messages pending, stop the timer
	// 			timer = nil
	// 		case timer == nil && pending:
	// 			// Timer is not already running and there are messages pending, so start it
	// 			timer = time.After(time.Duration(time.Second * cfg.BatchTimeOut))
	// 		default:
	// 			// Do nothing when:
	// 			// 1. Timer is already running and there are messages pending
	// 			// 2. Timer is not set and there are no messages pending
	// 		}

	// 	case <-timer:
	// 		//clear the timer
	// 		timer = nil

	// 		batch := ch.support.BlockCutter().Cut()
	// 		if len(batch) == 0 {
	// 			fmt.Println("Batch timer expired with no pending requests, this might indicate a bug")
	// 			continue
	// 		}
	// 		fmt.Println("Batch timer expired, creating block")
	// 		block, tipsList := ch.support.CreateNextBlock(batch)
	// 		ch.support.Append(block, tipsList)
	// 	case <-ch.exitChan:
	// 		fmt.Println("Exiting")
	// 		return
	// 	}
	// }
	rand.Seed(time.Now().Unix())
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
		block := ch.support.CreateNextBlock(batch)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)*4+300))
		ch.support.Append(block)
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
