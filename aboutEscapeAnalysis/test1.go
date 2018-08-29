package main

func foo() *int {
	x := 1
	return &x
}

func main() {
	x := foo()
	println(*x)
}
