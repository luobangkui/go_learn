package main

import (
	"github.com/nsqio/go-nsq"
	"log"
)

func main() {
	cfg := nsq.NewConfig()
	r := []byte("hello consumer from p1")
	r1 := []byte("hello consumer from p2")
	body := [][]byte{r, r1}
	addr := "127.0.0.1:9987"
	p, err := nsq.NewProducer(addr, cfg)
	if err != nil {
		log.Print(err)
	}
	err = p.MultiPublish("mytest", body)
	if err != nil {
		log.Println(err)
	}
}
