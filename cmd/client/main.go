package main

import (
	"log"

	"github.com/dyfromnil/pdag/client"
)

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile | log.Ltime)
	log.Println("Client start...")

	clt := client.NewClient(3)
	go clt.ReceiveReplyFromNodes()
	go clt.GenEnv()
	go clt.SendEnv()

	clt.WaitGracefulStop()
}
