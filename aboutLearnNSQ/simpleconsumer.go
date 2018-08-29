package main

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"math/rand"
	"os"
	"time"
)

type SimpleHandler struct {
}

func (sh *SimpleHandler) HandleMessage(m *nsq.Message) error {
	_, err := os.Stdout.Write(m.Body)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func main() {
	caddr := "127.0.0.1:9987"
	cfg := nsq.NewConfig()
	channel := fmt.Sprintf("tail%06d#ephemeral", rand.Int()%999999)
	c, _ := nsq.NewConsumer("mytest", channel, cfg)
	c.AddHandler(&SimpleHandler{})
	c.ConnectToNSQD(caddr)

	time.Sleep(100 * time.Second)
}
