package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/dyfromnil/pdag/globleconfig"
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
	log.Println("Sending messages...")

	for i := 0; i < client.numOfGorutine; i++ {
		i0 := i
		go func() {
			j := 0
			for {
				envCh <- cb.Envelope{
					Payload:    []byte(fmt.Sprintf("This is User %d, message %d", i0, j)),
					Signature:  []byte(fmt.Sprintf("User %d", i0)),
					Timestamp:  time.Now().UnixNano(),
					ClientAddr: globleconfig.ClientAddr,
				}
				j++
			}
		}()
	}
}

//客户端使用的tcp监听
func clientTCPListen() {
	listen, err := net.Listen("tcp", globleconfig.ClientAddr)
	if err != nil {
		log.Panic(err)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Panic(err)
		}
		b, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Panic(err)
		}
		log.Println(string(b))
	}

}

func main() {
	log.SetFlags(log.Ldate | log.Lshortfile | log.Ltime)
	log.Println("Client start...")

	go clientTCPListen()
	log.Printf("客户端开启监听，地址：%s\n", globleconfig.ClientAddr)

	envCh := make(chan cb.Envelope, 200)

	clt := NewClient(3)
	go clt.SendEnv(envCh)

	conn, err := grpc.Dial(globleconfig.LeaderListenEnvelopeAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	sendClient := cb.NewSendEnvelopsClient(conn)
	stream, err := sendClient.Request(context.Background())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
		panic("Request Error")
	}
	i := 0
	for {
		log.Println(i)
		i++
		env := <-envCh
		err = stream.Send(&env)
		time.Sleep(time.Millisecond * 2)
		if err != nil {
			panic("Send Error")
		}

		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error is:%s", err)
			panic("Receive Error")
		}

		log.Printf("Response:%t", resp.ResCode)
	}
	stream.CloseSend()
}
