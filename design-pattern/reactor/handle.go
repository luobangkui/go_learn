package reactor

import "fmt"



type Handle interface {
	run() error
	getType() EventType
}


type TimeHandle struct {
	hType EventType
}

func (h TimeHandle) run () error {
	fmt.Println("receive notify,timer run...")
	return nil
}

func (h TimeHandle) getType () EventType{
	return h.hType
}


type LoggingHandle struct {
	hType EventType
}

func (h LoggingHandle) run () error {
	fmt.Println("receive notify,loger run...")
	return nil
}

func (h LoggingHandle) getType () EventType{
	return h.hType
}


