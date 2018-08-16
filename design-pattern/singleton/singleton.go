package main

import (
	"sync"
	"fmt"
)

type single map[string]string

var (
	once sync.Once

	instance single
)

func New() single  {
	once.Do(func() {
		instance = make(single)
	})
	return instance
}

func main() {
	t := New()
	t["this"] = "that"
	b :=  New()
	b["this"] = "me"
	fmt.Println(t["this"])
}