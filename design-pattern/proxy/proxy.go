package main

import (
	"fmt"
	_ "net/http/pprof"
	"log"
	"net/http"
)


type IObject interface {
	ObjDo(action string)
}

type Object struct {
	action string
}

func (obj *Object) ObjDo(action string){
	fmt.Printf("i can %s",action)
}


type ProxyObject struct {
	object *Object
}

func (p *ProxyObject) ObjDo (action string)  {
	if p.object == nil {
		p.object = new(Object)
	}
	if action == "Run" || action == "run"{
		p.object.ObjDo(action)
	}
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	p := new(ProxyObject)
	p.ObjDo("Run")
	c := make(chan int)
	<-c
}
