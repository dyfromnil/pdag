package chain

import (
	"github.com/dyfromnil/pdag/chain/blockcutter"
	"github.com/dyfromnil/pdag/chain/blockledger"
	cb "github.com/dyfromnil/pdag/proto-go/common"
)

//Support for
type Support struct {
	ledger      blockledger.ReadWriter
	blockcutter blockcutter.Receiver
}

//NewSupport for
func NewSupport(lg blockledger.ReadWriter, bc blockcutter.Receiver) *Support {
	return &Support{
		ledger:      lg,
		blockcutter: bc,
	}
}

//CreateNextBlock for
func (s *Support) CreateNextBlock(messages []*cb.Envelope) *cb.Block {
	return s.ledger.CreateNextBlock(messages)
}
