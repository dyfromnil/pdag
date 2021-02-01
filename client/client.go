package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"

	"github.com/dyfromnil/pdag/globleconfig"
	cb "github.com/dyfromnil/pdag/proto-go/common"
)

//Client for
type Client struct {
	numOfGorutine  int
	receivePool    map[string]int
	isChecked      map[string]bool
	receiveNums    int
	lastCheckPoint checkPoint
	envCh          chan *cb.Envelope
	procRepCh      chan *Msg
	stopCh         chan bool
	file           *os.File
}

// Msg : receive consensus result from nodes
type Msg struct {
	Digest     string
	NumOfTranc int
	NodeID     string
}

type checkPoint struct {
	timeStamp int64
	receNums  int
}

//NewClient for
func NewClient(n int) *Client {
	file, err := os.OpenFile("./tps.log", os.O_RDWR|os.O_CREATE|os.O_CREATE, 0660)
	if err != nil {
		log.Fatal("error opening tps log file writer for file tps.log")
	}
	return &Client{
		numOfGorutine:  n,
		receivePool:    make(map[string]int),
		isChecked:      make(map[string]bool),
		receiveNums:    0,
		lastCheckPoint: checkPoint{0, 0},
		envCh:          make(chan *cb.Envelope, 2000),
		procRepCh:      make(chan *Msg, 2000),
		file:           file,
	}
}

// GenEnv for server starting
func (client *Client) GenEnv() {
	log.Println("Sending messages...")

	for i := 0; i < client.numOfGorutine; i++ {
		i0 := i
		go func() {
			j := 0
			for {
				client.envCh <- &cb.Envelope{
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

// SendEnv for
func (client *Client) SendEnv() {
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
		env := <-client.envCh
		err = stream.Send(env)
		time.Sleep(time.Millisecond * 1)
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

// ReceiveReplyFromNodes 客户端使用的tcp监听
func (client *Client) ReceiveReplyFromNodes() {
	listen, err := net.Listen("tcp", globleconfig.ClientAddr)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("客户端开启监听，地址：%s\n", globleconfig.ClientAddr)
	defer listen.Close()

	go client.updateReceiveNumsOrCheckPoint()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Panic(err)
		}
		b, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Panic(err)
		}
		var msg *Msg
		err = json.Unmarshal(b, &msg)
		if err != nil {
			log.Fatal("Unmarshal Error", err)
		}
		client.procRepCh <- msg
	}
}

func (client *Client) updateReceiveNumsOrCheckPoint() {
	ticker := time.NewTicker(time.Duration(time.Second))
	signCh := make(chan os.Signal)
	signal.Notify(signCh, os.Interrupt)
	for {
		select {
		case msg := <-client.procRepCh:
			client.receivePool[msg.Digest]++
			if client.receivePool[msg.Digest] > len(globleconfig.NodeTable)/3*2 && !client.isChecked[msg.Digest] {
				client.receiveNums += msg.NumOfTranc
				client.isChecked[msg.Digest] = true
			}
		case <-ticker.C:
			now := time.Now().UnixNano()
			interval := float64(now-client.lastCheckPoint.timeStamp) / 1e9
			tps := int(float64(client.receiveNums-client.lastCheckPoint.receNums) / interval)
			client.lastCheckPoint.receNums = client.receiveNums
			client.lastCheckPoint.timeStamp = now
			client.file.WriteString(fmt.Sprintln(time.Now(), " The current tps: ", tps))
		case <-signCh:
			client.file.Sync()
			client.file.Close()
			client.stopCh <- true
		}
	}
}

// WaitGracefulStop for
func (client *Client) WaitGracefulStop() {
	<-client.stopCh
}
