package main

import (
	"fmt"
	"sync"
	"time"
)

type F_b_im struct {
	bt           []*b_im
	btch         chan *b_im
	exitChan     chan int
	wg           sync.WaitGroup
	responseChan chan int
	closeChan    chan int
}

type b_im struct {
	doneChan chan *b_im
	exit     chan int
	bt       *b_im
}

//mock producer simulate the archetecture of nsq`s
// producer to show the direct of data flow
func main() {
	f := &F_b_im{
		btch:         make(chan *b_im),
		exitChan:     make(chan int),
		responseChan: make(chan int),
	}
	f.send()
}

func (f *F_b_im) send() {
	doneChan := make(chan *b_im)
	f.sendAsync(doneChan)
	<-doneChan
	fmt.Println("get donechan")
}

func (f *F_b_im) sendAsync(doneChan chan *b_im) error {
	t := &b_im{
		doneChan: doneChan,
	}

	f.connect()

	select {
	case f.btch <- t:
		fmt.Printf("btch get t\n")
	case <-f.exitChan:
		fmt.Printf("exit")
		return nil
	}
	return nil
}

func (f *F_b_im) connect() {
	f.closeChan = make(chan int)
	go f.router()
}

func (f *F_b_im) router() {
	for {
		select {
		case b := <-f.btch:
			f.bt = append(f.bt, b)
			goto exit
		case data := <-f.responseChan:
			f.popTransaction(data)
		case <-f.closeChan:
			goto exit
		case <-f.exitChan:
			goto exit

		}
	}
exit:
	f.bimCleanup()
	f.wg.Done()
	fmt.Printf("exiting router")
}

func (f *F_b_im) popTransaction(data int) {
	fmt.Printf("pop data %d\n", data)
	t := f.bt[0]
	f.bt = f.bt[:1]
	t.finish()
}

func (f *F_b_im) bimCleanup() {
	for _, t := range f.bt {
		t.finish()
	}
	f.bt = f.bt[:0]
	for {
		select {
		case b := <-f.btch:
			b.finish()
		default:
			time.Sleep(5 * time.Second)
		}
	}
}

func (b *b_im) finish() {
	if b.doneChan != nil {
		b.doneChan <- b
		fmt.Printf("finish bim")
	}
}
