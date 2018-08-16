package main

import "fmt"

type Object struct {
	a int
}

func (ob *Object) Do()  {
	fmt.Println("obj do!")
}


type Pool chan *Object

func New(total int) *Pool {
	p := make(Pool,total)
	for i:=0;i< total;i++{}
	p<-new(Object)
	return &p
}

func main() {
	p := New(5)
	for{
		select {
		case obj := <- *p:
			obj.Do()
			*p <- obj
		default:
			return
		}
	}
	return

}


