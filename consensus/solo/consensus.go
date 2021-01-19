package solo

import (
	"fmt"
	"time"

	"github.com/dyfromnil/pdag/consensus"
	cfg "github.com/dyfromnil/pdag/globleconfig"
	cb "github.com/dyfromnil/pdag/proto-go/common"
)

type consenter struct{}

type chain struct {
	support  consensus.ConsenterSupport
	sendChan chan *message
	exitChan chan struct{}
}

type message struct {
	configSeq uint64
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
func (ch *chain) Order(env *cb.Envelope, configSeq uint64) error {
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
	var timer <-chan time.Time

	for {
		select {
		case msg := <-ch.sendChan:
			batches, pending := ch.support.BlockCutter().Ordered(msg.normalMsg)

			for _, batch := range batches {
				block := ch.support.CreateNextBlock(batch)
				ch.support.Append(block)
			}

			switch {
			case timer != nil && !pending:
				// Timer is already running but there are no messages pending, stop the timer
				timer = nil
			case timer == nil && pending:
				// Timer is not already running and there are messages pending, so start it
				timer = time.After(time.Duration(time.Second * cfg.BatchTimeOut))
			default:
				// Do nothing when:
				// 1. Timer is already running and there are messages pending
				// 2. Timer is not set and there are no messages pending
			}

		case <-timer:
			//clear the timer
			timer = nil

			batch := ch.support.BlockCutter().Cut()
			if len(batch) == 0 {
				fmt.Println("Batch timer expired with no pending requests, this might indicate a bug")
				continue
			}
			fmt.Println("Batch timer expired, creating block")
			block := ch.support.CreateNextBlock(batch)
			// ch.support.WriteBlock(block, nil)
			ch.support.Append(block)
		case <-ch.exitChan:
			fmt.Println("Exiting")
			return
		}
	}
}
