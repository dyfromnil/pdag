package chain

import (
	"github.com/dyfromnil/pdag/chain/blockcutter"
	"github.com/dyfromnil/pdag/chain/blockledger"
	cb "github.com/dyfromnil/pdag/proto-go/common"
)

//Support for
type Support struct {
	blockledger.ReadWriter
	blockcutter.Receiver
}

//NewSupport for
func NewSupport(lg blockledger.ReadWriter) *Support {
	return &Support{
		ReadWriter: lg,
		Receiver:   blockcutter.NewReceiverImpl(),
	}
}

//CreateNextBlock for
func (s *Support) CreateNextBlock(messages []*cb.Envelope) (*cb.Block, []*cb.Block) {
	return s.ReadWriter.CreateNextBlock(messages)
}

//BlockCutter for
func (s *Support) BlockCutter() blockcutter.Receiver {
	return s.Receiver
}
