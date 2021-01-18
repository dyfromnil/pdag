package node

import (
	"fmt"
	cb "github.com/dyfromnil/pdag/proto-go/common"
)

//Main for server starting
func Main() {
	fmt.Println("start...")

	envCh := make(chan cb.Envelope, 200)
	for i := 0; i < 3; i++ {
		i0 := i
		go func() {
			j := 0
			for {
				envCh <- cb.Envelope{
					User:      uint32(i0),
					Payload:   []byte(fmt.Sprintf("This is User %d, message %d", i0, j)),
					Signature: []byte(fmt.Sprintf("User %d", i0)),
				}
				j++
			}
		}()
	}

}
