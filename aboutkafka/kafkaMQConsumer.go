package main

import (
	"fmt"
	cluster "github.com/bsm/sarama-cluster"
	"log"
	"os"
	"os/signal"
)

func main() {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	//init consumer
	brokers := []string{"127.0.0.1:9092"}
	topic := []string{"mytopic", "othertopic"}
	consumer, err := cluster.NewConsumer(brokers, "my-consumer-group", topic, config)
	if err != nil {
		//panic(err)
	}
	defer consumer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		for err := range consumer.Errors() {
			log.Print("Error:%s\n", err.Error())
		}
	}()

	go func() {
		for ntf := range consumer.Notifications() {
			log.Print("Rebalanced:%+v\n", ntf)
		}
	}()

	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				consumer.MarkOffset(msg, "") // mark message as processed
			}
		case <-signals:
			return
		}
	}
}
