### 关于nsq的Consumer


nsq里的consumer使用方式十分简单
```
caddr := "127.0.0.1:9987"
cfg := nsq.NewConfig()
channel := fmt.Sprintf("tail%06d#ephemeral", rand.Int()%999999)
c, _ := nsq.NewConsumer("mytest", channel, cfg)
c.AddHandler(&SimpleHandler{})
c.ConnectToNSQD(caddr)
```

只需要新建一个consumer，添加一个实现了```HandleMessage(message *Message) error```方法的结构体即可。但是具体内部是如何实现的，跟踪源码看一下consumer是如何处理接口的方法。

```ConnectToNSQD```方法顾名思义，也就是去连接一个nsqd，看一下代码里面，这里省略了一些代码

```
func (r *Consumer) ConnectToNSQD(addr string) error {
	...
	conn := NewConn(addr, &r.config, &consumerConnDelegate{r})
	...
	resp, err := conn.Connect()
	if err != nil {
		cleanupConnection()
		return err
	}

	if resp != nil {
		if resp.MaxRdyCount < int64(r.getMaxInFlight()) {
			r.log(LogLevelWarning,
				"(%s) max RDY count %d < consumer max in flight %d, truncation possible",
				conn.String(), resp.MaxRdyCount, r.getMaxInFlight())
		}
	}

	cmd := Subscribe(r.topic, r.channel)
	err = conn.WriteCommand(cmd)
	...
	return nil
}
```
可以看到，这里使用了NewConn方法建立了一个连接体。具体的操作方法在 ```resp, err := conn.Connect()```，也就是在这里面调用了HandleMessage的方法，处理消息。


看一下这个函数
```
// Connect dials and bootstraps the nsqd connection
// (including IDENTIFY) and returns the IdentifyResponse
func (c *Conn) Connect() (*IdentifyResponse, error) {
	dialer := &net.Dialer{
		LocalAddr: c.config.LocalAddr,
		Timeout:   c.config.DialTimeout,
	}

	conn, err := dialer.Dial("tcp", c.addr)
	if err != nil {
		return nil, err
	}
	c.conn = conn.(*net.TCPConn)
	c.r = conn
	c.w = conn

	_, err = c.Write(MagicV2)
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("[%s] failed to write magic - %s", c.addr, err)
	}

	resp, err := c.identify()
	if err != nil {
		return nil, err
	}

	if resp != nil && resp.AuthRequired {
		if c.config.AuthSecret == "" {
			c.log(LogLevelError, "Auth Required")
			return nil, errors.New("Auth Required")
		}
		err := c.auth(c.config.AuthSecret)
		if err != nil {
			c.log(LogLevelError, "Auth Failed %s", err)
			return nil, err
		}
	}

	c.wg.Add(2)
	atomic.StoreInt32(&c.readLoopRunning, 1)
	go c.readLoop()
	go c.writeLoop()
	return resp, nil
}
```
这个函数的主要作用是，建立了一个TCP的会话，然后新开两个goroutine分别处理读写的loop。
忽略上面的验证和连接，直接看最下面两个goroutine。首先看```c.readLoop()```，

这个函数，主要是验证心跳监测，以及对message进行解码。这里使用了conn.delegate作为代理，这个代理是为了传递给下面的message，用来给message调用后面的消息处理。下面也使用了conn代理的on方法，on方法是一系列处理消息的方法，接受一个conn和要处理的message。

```
func (c *Conn) readLoop() {
	delegate := &connMessageDelegate{c}
	for {
		if atomic.LoadInt32(&c.closeFlag) == 1 {
			goto exit
		}

		frameType, data, err := ReadUnpackedResponse(c)
		if err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				c.log(LogLevelError, "IO error - %s", err)
				c.delegate.OnIOError(c, err)
			}
			goto exit
		}

		if frameType == FrameTypeResponse && bytes.Equal(data, []byte("_heartbeat_")) {
			c.log(LogLevelDebug, "heartbeat received")
			c.delegate.OnHeartbeat(c)
			err := c.WriteCommand(Nop())
			if err != nil {
				c.log(LogLevelError, "IO error - %s", err)
				c.delegate.OnIOError(c, err)
				goto exit
			}
			continue
		}

		switch frameType {
		case FrameTypeResponse:
			c.delegate.OnResponse(c, data)
		case FrameTypeMessage:
			msg, err := DecodeMessage(data)
			if err != nil {
				c.log(LogLevelError, "IO error - %s", err)
				c.delegate.OnIOError(c, err)
				goto exit
			}
			msg.Delegate = delegate
			msg.NSQDAddress = c.String()

			atomic.AddInt64(&c.rdyCount, -1)
			atomic.AddInt64(&c.messagesInFlight, 1)
			atomic.StoreInt64(&c.lastMsgTimestamp, time.Now().UnixNano())

			c.delegate.OnMessage(c, msg)
		case FrameTypeError:
			c.log(LogLevelError, "protocol error - %s", data)
			c.delegate.OnError(c, data)
		default:
			c.log(LogLevelError, "IO error - %s", err)
			c.delegate.OnIOError(c, fmt.Errorf("unknown frame type %d", frameType))
		}
	}

exit:
	...
}
```

conn代理接口

```
type ConnDelegate interface {
	// OnResponse is called when the connection
	// receives a FrameTypeResponse from nsqd
	OnResponse(*Conn, []byte)

	// OnError is called when the connection
	// receives a FrameTypeError from nsqd
	OnError(*Conn, []byte)

	// OnMessage is called when the connection
	// receives a FrameTypeMessage from nsqd
	OnMessage(*Conn, *Message)

	// OnMessageFinished is called when the connection
	// handles a FIN command from a message handler
	OnMessageFinished(*Conn, *Message)

	// OnMessageRequeued is called when the connection
	// handles a REQ command from a message handler
	OnMessageRequeued(*Conn, *Message)

	// OnBackoff is called when the connection triggers a backoff state
	OnBackoff(*Conn)

	// OnContinue is called when the connection finishes a message without adjusting backoff state
	OnContinue(*Conn)

	// OnResume is called when the connection triggers a resume state
	OnResume(*Conn)

	// OnIOError is called when the connection experiences
	// a low-level TCP transport error
	OnIOError(*Conn, error)

	// OnHeartbeat is called when the connection
	// receives a heartbeat from nsqd
	OnHeartbeat(*Conn)

	// OnClose is called when the connection
	// closes, after all cleanup
	OnClose(*Conn)
}
```

回到ConnectToNSQD，找到```conn := NewConn(addr, &r.config, &consumerConnDelegate{r})```,这里的```consumerConnDelegate```是consumer的代理类，持有一个consumer结构体，并且实现了上面的ConnDelegate接口。


找到这个类实现的OnMessage方法
```
func (d *consumerConnDelegate) OnMessage(c *Conn, m *Message)         { d.r.onConnMessage(c, m) }
```
它调用了consumer的onConnMessage方法，也就是在这里面触发了消息channel
```
func (r *Consumer) onConnMessage(c *Conn, msg *Message) {
	atomic.AddInt64(&r.totalRdyCount, -1)
	atomic.AddUint64(&r.messagesReceived, 1)
	// 传递了msg消息体，触发了consumer消费
	r.incomingMessages <- msg
	r.maybeUpdateRDY(c)
}
```

回到我们最开始写的demo，进入```AddHandler```->AddConcurrentHandlers```，可以看到这里跑了一个goroutine去执行一个handlerLoop方法.
```
func (r *Consumer) AddConcurrentHandlers(handler Handler, concurrency int) {
	if atomic.LoadInt32(&r.connectedFlag) == 1 {
		panic("already connected")
	}

	atomic.AddInt32(&r.runningHandlers, int32(concurrency))
	for i := 0; i < concurrency; i++ {
		go r.handlerLoop(handler)
	}
}
```

看一下handlerLoop内部

```
func (r *Consumer) handlerLoop(handler Handler) {
	r.log(LogLevelDebug, "starting Handler")

	for {
	    //接受到msg
		message, ok := <-r.incomingMessages
		if !ok {
			goto exit
		}

		if r.shouldFailMessage(message, handler) {
			message.Finish()
			continue
		}
        //调用HandleMessage处理message
		err := handler.HandleMessage(message)
		if err != nil {
			r.log(LogLevelError, "Handler returned error (%s) for msg %s", err, message.ID)
			if !message.IsAutoResponseDisabled() {
				message.Requeue(-1)
			}
			continue
		}

		if !message.IsAutoResponseDisabled() {
			message.Finish()
		}
	}

exit:
	r.log(LogLevelDebug, "stopping Handler")
	if atomic.AddInt32(&r.runningHandlers, -1) == 0 {
		r.exit()
	}
}
```


