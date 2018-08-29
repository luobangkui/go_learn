package main

import (
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"sync"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var (
		wg                          sync.WaitGroup
		enqueued, successes, errors int
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			successes++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			log.Print(err)
			errors++
		}
	}()

ProducerLoop:
	for {
		message := &sarama.ProducerMessage{Topic: "mytopic", Value: sarama.StringEncoder("testing 1234")}
		select {
		case producer.Input() <- message:
			enqueued++
		case <-signals:
			producer.AsyncClose()
			break ProducerLoop
		}
	}

	wg.Wait()
	log.Print("Successfully produced:%d\n", successes, errors)
}
