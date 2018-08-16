package reactor

import (
	"testing"
	"time"
	"fmt"
	"net"
	"io"
	"strconv"
	"strings"
)

func TestClient(t *testing.T) {
	initor := Initor{sy:SynchronousEventDemultiplexer{}}
	timer := TimeHandle{hType:TIMER_EVENT}
	logger := LoggingHandle{hType:LOGGING_EVENT}

	initor.register_handler(timer)
	initor.register_handler(logger)

	ln ,err := net.Listen("tcp","localhost:9902")
	if err != nil{
		panic(err)
	}

	for	 {
		conn,err := ln.Accept()
		if err != nil{
			fmt.Println("get connection error!")
		}

		for{
			var buf = make([]byte,10)
			conn.Read(buf)
			if err == io.EOF{
				break
			}
			var slide []string = strings.Split(strings.Trim(string(buf),"\r\n"),",")
			e1,err1:=strconv.Atoi(slide[0])
			_,err2:=strconv.Atoi(slide[1])
			fmt.Println(slide[0],err1,err2)
			if  err1!=nil || err2!=nil{
				fmt.Println(err1,err2)
			}

			initor.handle_events(EventType(e1))
			//initor.handle_events(EventType(e2))
		}
	}

	for{
		fmt.Println("loop")
	}



	if err != nil {
		t.Error("handle wrong!!")
	}
	time.Sleep(1*time.Second)

}







