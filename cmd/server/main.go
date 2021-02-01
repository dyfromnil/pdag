package main

import (
	"log"

	"github.com/dyfromnil/pdag/server"
)

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile | log.Ltime)
	server.Main()
}
