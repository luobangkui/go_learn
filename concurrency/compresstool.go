package main

import (
	"os"
	"compress/gzip"
	"io"
	"sync"
	"fmt"
)

func main() {
	var wg sync.WaitGroup
	var i int = -1
	var file string
	for i ,file = range os.Args[1:] {
		wg.Add(1)
		go func(filename string) {
			compress(filename)
			wg.Done()
		}(file)
	}
	wg.Wait()
	fmt.Printf("Compressed %d files\n", i+1)
}


func compress(filename string) error {
	in, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(filename + ".gz")
	if err != nil {
		return err
	}
	defer out.Close()
	gzout := gzip.NewWriter(out)
	_, err = io.Copy(gzout, in)
	gzout.Close()
	return err
}