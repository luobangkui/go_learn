package reactor

import (
	"fmt"
	"errors"
)

type EventType int32

const    (
	LOGGING_EVENT EventType = iota
	TIMER_EVENT
)


//selector notifies identified handler to excute a handle fuction
type SynchronousEventDemultiplexer struct {
	handler Handler
}

func (s *SynchronousEventDemultiplexer) Select(i EventType) error  {
	//if we only have one handler
	handler := s.handler
	handler.handle_event(i)
	return nil
}


type InitiationDispatcher interface {
	handle_events() error
	register_handler(h *Handler) error
	remove_handler(h *Handler) error
}



type Initor struct {
	sy SynchronousEventDemultiplexer
}

func (initor *Initor) handle_events (i EventType) error  {
	initor.sy.Select(i)
	return nil
}

func (initor *Initor) register_handler(h Handle) error  {
	initor.sy.handler.handles = append(initor.sy.handler.handles,h)
	return nil
}

func (initor *Initor) remove_handler(h *Handler) error  {
	//TODO
	return nil
}




type EventHandler interface {
	handle_event(EventType) error
	get_handler()
}

type Handler struct {
	handles []Handle
}

func (h Handler) handle_event(e EventType) error  {
	var handle Handle
	for i,_ := range h.handles{
		if h.handles[i].getType() == e{
			handle = h.handles[i]
		}
	}
	if handle == nil{
		fmt.Println("don`t find handle !!")
		return errors.New("don`t find handle ")
	}
	handle.run()
	return nil
}

func (h Handler) get_handle(){
	//TODO
}






