package main

import (
	"testing"
)

func main() {
	//laddr := "127.0.0.1"
	//
	//addr, _ := net.ResolveTCPAddr("tcp", laddr+":0")
	//timeout := time.Duration(5 * time.Second)
	//
	//dialer := &net.Dialer{
	//	LocalAddr: addr,
	//	Timeout:   timeout,
	//}
	//
	//conn, err := dialer.Dial("tcp", "127.0.0.1:9987")
	//if err != nil {
	//	panic(err)
	//}
	//conn.fl
	//fmt.Println(myAtoi("words and 987"))

}

func TestMyatoi(t *testing.T) {
	t1 := "42"
	t2 := "   -42"
	t3 := "4193 with words"
	t4 := "words and 987"
	t5 := "-91283472332"
	t6 := "3.14159"
	t7 := "+1"
	t8 := "  -0012a42"
	t9 := "   +0 123"
	t10 := "9223372036854775808"
	t11 := "2147483648"
	if myAtoi(t1) != 42 {
		t.Error("42 is wrong")
	}
	if myAtoi(t2) != -42 {
		t.Errorf("'%s' is wrong", t2)
	}
	if myAtoi(t3) != 4193 {
		t.Errorf("'%s' is wrong", t3)
	}
	if myAtoi(t4) != 0 {
		t.Errorf("'%s' is wrong", t4)
	}
	if myAtoi(t5) != -2147483648 {
		t.Errorf("'%s' is wrong expect %d", t5, myAtoi(t5))
	}
	if myAtoi(t6) != 3 {
		t.Errorf("'%s' is wrong expect %d", t6, myAtoi(t6))
	}
	if myAtoi(t7) != 1 {
		t.Errorf("'%s' is wrong expect %d", t7, myAtoi(t7))
	}
	if myAtoi(t8) != -12 {
		t.Errorf("'%s' is wrong get %d,expect %d", t8, myAtoi(t8), -12)
	}
	if myAtoi(t9) != 0 {
		t.Errorf("'%s' is wrong get %d,expect %d", t9, myAtoi(t9), 0)
	}
	if myAtoi(t10) != 2147483647 {
		t.Errorf("'%s' is wrong get %d,expect %d", t10, myAtoi(t10), 2147483647)
	}
	if myAtoi(t11) != 2147483647 {
		t.Errorf("'%s' is wrong get %d,expect %d", t11, myAtoi(t11), 2147483647)
	}

}

func myAtoi(str string) int {
	isNum := func(a byte) bool {
		return a >= '0' && a <= '9'
	}
	bytesarry := []byte(str)
	isNagtive := false
	result := 0
	getfirst := false
	getEnd := false
	for i := 0; i < len(bytesarry); i++ {

		if isNagtive {
			if 0-result < 0-(1<<31) {
				return 0 - (1 << 31)
			}
		} else {
			if result > 1<<31-1 {
				return 1<<31 - 1
			}
		}

		if getfirst && !getEnd {
			if isNum(bytesarry[i]) {
				bytesarry[i] -= '0'
				result = result*10 + (int)(bytesarry[i])
			} else {
				break
			}
		}
		if bytesarry[i] == ' ' {
			if !getfirst {
				continue
			}
			break
		} else {
			if !getfirst {
				getfirst = true
				if !isNum(bytesarry[i]) && bytesarry[i] != '-' && bytesarry[i] != '+' {
					return 0
				}
				if bytesarry[i] == '-' {
					isNagtive = true
				} else if bytesarry[i] == '+' {
					isNagtive = false
				} else {
					bytesarry[i] -= '0'
					result = result*10 + (int)(bytesarry[i])
				}
				continue
			}
		}
	}
	if isNagtive {
		result = 0 - result
	}
	if result > 1<<31-1 {
		return 1<<31 - 1
	}
	if result < 0-(1<<31) {
		return 0 - (1 << 31)
	}
	return result
}
