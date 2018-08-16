package main

import (
	"sync"
	"runtime"
	"fmt"
	_"time"
)

// he use a wrapper to package a func which should be locked like this
//*****************************************************
//type NSQLookupd struct {
//	sync.RWMutex
//	opts         *Options
//	tcpListener  net.Listener
//	httpListener net.Listener
//	waitGroup    util.WaitGroupWrapper
//	DB           *RegistrationDB
//}
//*****************************************************
//type WaitGroupWrapper struct {
//	sync.WaitGroup
//}
//*****************************************************
//func (w *WaitGroupWrapper) Wrap(cb func()) {
//	w.Add(1)
//	go func() {
//		cb()
//		w.Done()
//	}()
//}
//*****************************************************
//l.waitGroup.Wrap(func() {
//	protocol.TCPServer(tcpListener, tcpServer, l.logf)
//})
type WaitGroupWrapper struct {
	sync.WaitGroup
}

func (w *WaitGroupWrapper) Wrap (cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}

type Obj struct {
	waitGroup WaitGroupWrapper
	scala int
}

var (
	count int32
	wg sync.WaitGroup
	mutex sync.Mutex
)

func main() {
	obj := &Obj{}
	wg.Add(4)
	obj.waitGroup.Wrap(func() {
		incCount()
	})
	obj.waitGroup.Wrap(func() {
		incCount()
	})
	go incCount()
	go incCount()
	wg.Wait()
	fmt.Println(count)
}

func incCount() {
	defer func() {
		wg.Done()
	}()
	for i := 0; i < 2; i++ {
		mutex.Lock()
		value := count
		runtime.Gosched()
		value++
		count = value
		mutex.Unlock()
	}

}









