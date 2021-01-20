package node

import (
	"fmt"
	"github.com/dyfromnil/pdag/chain"
	"github.com/dyfromnil/pdag/chain/blkstorage"
	"github.com/dyfromnil/pdag/chain/blockledger/fileledger"
	"github.com/dyfromnil/pdag/client"
	"github.com/dyfromnil/pdag/consensus/solo"
	cb "github.com/dyfromnil/pdag/proto-go/common"
	"time"
)

//Main for server starting
func Main() {
	fmt.Println("start...")

	envCh := make(chan cb.Envelope, 200)

	clt := client.NewClient(3)
	clt.SendEnv(envCh)

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
		case env := <-envCh:
			err := soloChain.Order(&env)
			if err != nil {
				panic("order panic!")
			}
		case <-timer:
			break
		}
	}
}
