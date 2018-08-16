package main

import (
	"github.com/judwhite/go-svc/svc"
	"fmt"
	"syscall"
	"log"
)

type myob struct {

}


func main() {
	prg := &myob{}
	if err := svc.Run(prg, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatal(err)
	}
}

func (mo *myob) Init(env svc.Environment) error {
	fmt.Println("init in svc")
	return nil
}

func (mo *myob) Start() error  {
	fmt.Println("start")
	return nil
}

func (mo *myob) Stop()	error  {
	fmt.Println("stop!")
	return nil
}