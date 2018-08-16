package main

import (
	"strconv"
	"fmt"
)

/**
Given a 32-bit signed integer, reverse digits of an integer.

Example 1:

Input: 123
Output: 321
Example 2:

Input: -123
Output: -321
Example 3:

Input: 120
Output: 21
Note:
Assume we are dealing with an environment which could only store integers within the 32-bit signed integer range: [−231,  231 − 1]. For the purpose of this problem, assume
that your function returns 0 when the reversed integer overflows.
 */

func reverse(x int) int {
	a := strconv.Itoa(x)
	var isNegtive bool = false
	if a[0] == '-' {
		isNegtive=true
		a = a[1:]
	}
	lenx := len(a)
	var b = make([]byte,len(a))
	for i:=0;i<lenx;i++{
		b[lenx-i-1] = a[i]
	}
	var res int
	if isNegtive {
		res,_ = strconv.Atoi("-"+string(b))
	}else {
		res,_ = strconv.Atoi(string(b))
	}
	if x < -(1<<31) || x > 1<<31-1 || res < -(1<<31) || res > 1<<31-1 {
		return 0
	}
	return res
}

func main() {
	fmt.Println(reverse(-123))
}
