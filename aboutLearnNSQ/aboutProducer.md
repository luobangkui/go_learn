## nsq里的Producer的数据流动
以TCP的方式为例子来看一下producer是如何工作的
```
	cfg := nsq.NewConfig()
	r := []byte("hello consumer")
	addr := "127.0.0.1:9987"
	p, err := nsq.NewProducer(addr, cfg)
	if err != nil {
		log.Print(err)
	}
	err = p.Publish("mytest", r)
	if err != nil {
		log.Println(err)
	}
```
首先使用producer需要调用```nsq.NewProducer```方法新建producer实例，方法接受一个地址和```Config```类型的指针,然后调用```producer.Publish```方法，```Publish```接受string类型的topic和[]byte类型的消息实体。

具体内部是如何实现逻辑的，跟踪进入Publish方法
```
func (w *Producer) Publish(topic string, body []byte) error {
	return w.sendCommand(Publish(topic, body))
}
```
函数直接调用了```sendCommand```方法,接受类型是一个Command指针。包裹了的一层```Publish(topic,body)```方法作为传递，这个方法将topic和body数据做了封装
```
func Publish(topic string, body []byte) *Command {
	var params = [][]byte{[]byte(topic)}
	return &Command{[]byte("PUB"), params, body}
}
```

看一下sendCommand方法的内部
```
func (w *Producer) sendCommand(cmd *Command) error {
	doneChan := make(chan *ProducerTransaction)
	err := w.sendCommandAsync(cmd, doneChan, nil)
	if err != nil {
		close(doneChan)
		return err
	}
	t := <-doneChan
	return t.Error
}
```

这个方法主要是创建了一个doneChan，调用了```
w.sendCommandAsync(cmd, doneChan, nil)```方法，这个方法一会儿说，关注一下doneChan和```ProducerTransaction```这个结构类型.
```
type ProducerTransaction struct {
	cmd      *Command
	doneChan chan *ProducerTransaction
	Error    error         // the error (or nil) of the publish command
	Args     []interface{} // the slice of variadic arguments passed to PublishAsync or MultiPublishAsync
}
```
这个结构体持有了一个包含自己类型指针的chan。doneChan作为参数传入```w.sendCommandAsync```方法里面,然后会塞进一个```ProducerTransaction```里面，再把这个```ProducerTransaction```塞给Producer,由```transactionChan chan *ProducerTransaction```存着，在做cleanup的时候，producer持有的这个transactionChan取出数据,然后调用```finish```方法，finishi把自己的指针放进自己的chan里面:)好他喵诡异的写法。这个t.doneChan也就是```sendCommand```里面的doneChan。
```
func (t *ProducerTransaction) finish() {
	if t.doneChan != nil {
		t.doneChan <- t
	}
}
```
具体doneChan是怎么传递的，进入```sendCommandAsync```里面,这里面建立了一个ProducerTransaction的指针,然后传递给```w.transactionChan```管道里面。取出这个```w.transactionChan```则是在w.connect方法里的router方法进行的。```sendCommandAsync```主要是使用原子方法，记录当前并发producer的数量,确保能最后把所有producer清理
```
func (w *Producer) sendCommandAsync(cmd *Command, doneChan chan *ProducerTransaction,
	args []interface{}) error {
	// keep track of how many outstanding producers we're dealing with
	// in order to later ensure that we clean them all up...
	atomic.AddInt32(&w.concurrentProducers, 1)
	defer atomic.AddInt32(&w.concurrentProducers, -1)

	if atomic.LoadInt32(&w.state) != StateConnected {
		err := w.connect()
		if err != nil {
			return err
		}
	}

	t := &ProducerTransaction{
		cmd:      cmd,
		doneChan: doneChan,
		Args:     args,
	}

	select {
	case w.transactionChan <- t:
	case <-w.exitChan:
		return ErrStopped
	}

	return nil
}
```
具体的数据处理和流动都在w.connect()方法里面做处理，可以看到，首先验证当前producer的状态,如果已连接则返回nil，默认返回未连接错误，如果是初始化，后面会新建立连机器,然后开一个协程处理router()方法,连接器做了很多层的封装，暂且不必关注，主要看router()方法
```
func (w *Producer) connect() error {
	w.guard.Lock()
	defer w.guard.Unlock()

	if atomic.LoadInt32(&w.stopFlag) == 1 {
		return ErrStopped
	}

	switch state := atomic.LoadInt32(&w.state); state {
	case StateInit:
	case StateConnected:
		return nil
	default:
		return ErrNotConnected
	}

	w.log(LogLevelInfo, "(%s) connecting to nsqd", w.addr)

	logger, logLvl := w.getLogger()

	w.conn = NewConn(w.addr, &w.config, &producerConnDelegate{w})
	w.conn.SetLogger(logger, logLvl, fmt.Sprintf("%3d (%%s)", w.id))

	_, err := w.conn.Connect()
	if err != nil {
		w.conn.Close()
		w.log(LogLevelError, "(%s) error connecting to nsqd - %s", w.addr, err)
		return err
	}
	atomic.StoreInt32(&w.state, StateConnected)
	w.closeChan = make(chan int)
	w.wg.Add(1)
	go w.router()

	return nil
}
```

router方法里面开了一个for循环来监听producer的数据或者是信号流入，如果是```w.transactionChan```，则把这个transaction塞进producer的transactions切片里，如果是收到了返回信号或者是错误信号，则会弹出一个transaction。如果收到关闭或者是退出信号，则到exit里面，清理所有transaction，并退出。
```
func (w *Producer) router() {
	for {
		select {
		case t := <-w.transactionChan:
			w.transactions = append(w.transactions, t)
			err := w.conn.WriteCommand(t.cmd)
			if err != nil {
				w.log(LogLevelError, "(%s) sending command - %s", w.conn.String(), err)
				w.close()
			}
		case data := <-w.responseChan:
			w.popTransaction(FrameTypeResponse, data)
		case data := <-w.errorChan:
			w.popTransaction(FrameTypeError, data)
		case <-w.closeChan:
			goto exit
		case <-w.exitChan:
			goto exit
		}
	}

exit:
	w.transactionCleanup()
	w.wg.Done()
	w.log(LogLevelInfo, "exiting router")
}
```

弹出transaction方法很简单，他会把第一个transaction弹出，调用一下它的finish()方法。这个finishi方法就是前面说到的，它把错误数据进行了回传。
```
func (w *Producer) popTransaction(frameType int32, data []byte) {
	t := w.transactions[0]
	w.transactions = w.transactions[1:]
	if frameType == FrameTypeError {
		t.Error = ErrProtocol{string(data)}
	}
	t.finish()
}
```
这样一个框架完成了producer的publish流程。具体如何触发退出或者是其他信号，都可以在源码里看到


