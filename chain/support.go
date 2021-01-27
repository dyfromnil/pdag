package chain

import (
	"github.com/dyfromnil/pdag/chain/blockcutter"
	"github.com/dyfromnil/pdag/chain/blockledger"
	"github.com/dyfromnil/pdag/msp"
	cb "github.com/dyfromnil/pdag/proto-go/common"
)

//Support for
type Support struct {
	blockledger.ReadWriter
	blockcutter.Receiver
	msp.IdentityProvider
}

//NewSupport for
func NewSupport(lg blockledger.ReadWriter, ident msp.IdentityProvider) *Support {
	return &Support{
		ReadWriter:       lg,
		Receiver:         blockcutter.NewReceiverImpl(),
		IdentityProvider: ident,
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

//GetIdendity for
func (s *Support) GetIdendity() msp.IdentityProvider {
	return s.IdentityProvider
}
