package main

import (
	"log"

	"github.com/dyfromnil/pdag/client"
	"github.com/dyfromnil/pdag/globleconfig"
)

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile | log.Ltime)
	log.Println("Client start...")

	clt := client.NewClient(globleconfig.NumOfClient)
	go clt.ReceiveReplyFromNodes()
	go clt.GenEnv()
	go clt.SendEnv()

	clt.WaitGracefulStop()
}
