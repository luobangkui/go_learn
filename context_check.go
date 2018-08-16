package main

import (
	"sync"
	"fmt"
	"time"
	"context"
)

var (
	wg sync.WaitGroup
	all_contrl int
)


func main() {
	all_contrl = 15
	ctx , cancle := context.WithTimeout(context.Background(),120*time.Microsecond)
	defer cancle()
	start := time.Now()
	wg = sync.WaitGroup{}
	for i := 1;i<10;i++{
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			all_contrl+=3

			time.Sleep(110*time.Microsecond)
			select {
			case <-ctx.Done():
				fmt.Printf("wait time out!!  %d\n",i)
				//os.Exit(0)
			default:
				fmt.Printf("---->%d goroutine run, control = %d\n",i,all_contrl)
			}
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println("elapsed :",elapsed)


}



