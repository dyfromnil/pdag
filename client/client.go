package client

import (
	"fmt"
	cb "github.com/dyfromnil/pdag/proto-go/common"
)

//Client for
type Client struct {
	numOfGorutine int
}

//NewClient for
func NewClient(n int) *Client {
	return &Client{numOfGorutine: n}
}

//SendEnv for server starting
func (client *Client) SendEnv(envCh chan<- cb.Envelope) {
	fmt.Println("Sending messages...")

	for i := 0; i < client.numOfGorutine; i++ {
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
