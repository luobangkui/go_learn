package main

import "log"

type Object func(int) int

func LogDecorate(fn Object) Object {
	return func(i int) int {
		log.Println("starting the execution with the integer",i)

		result := fn(i)

		log.Println("this is complete with result",result)

		return result
	}
}

func Double(n int) int {
	return n*n
}

func main() {
	f := LogDecorate(Double)
	f(5)
}


