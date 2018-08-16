package main

import (
	"fmt"
	_ "os"
	pb "code.admaster.co/smartserving"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	address := "localhost:55555"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		//TODO
	}
	defer conn.Close()
	c := pb.NewHelloClient(conn)
	name := "Inigo Montoya"
	hr := &pb.HelloRequest{Name: name}
	r, err := c.Say(context.Background(), hr)
	if err != nil {
		//TODO
	}
	fmt.Println(r.Message)
}