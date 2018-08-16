package main

import (
	"fmt"
	"strings"
)

var b string = "   sss fff sst   \r\n"

func main() {


	//fmt.Println(b) // print the content as 'bytes'


	str := strings.TrimSpace(b)
	//str := string(b) // convert content to a 'string'


	fmt.Printf(str) // print the content as a 'string'
	fmt.Println("---<")

}