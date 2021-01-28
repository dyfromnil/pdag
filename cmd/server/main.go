package main

import (
	"log"

	"github.com/dyfromnil/pdag/node"
)

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile | log.Ltime)
	node.Main()
}
