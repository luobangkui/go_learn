package main

/**
APPLE 3

A E
PL
P

OUTPUT : AEPLP
*/

func convert(s string, numRows int) string {
	row := numRows
	column := len(s)
	array := make([][]byte, row)
	res :=""
	for i := range array {
		array[i] = make([]byte, column)
	}
	//fill
	for i := 0; i < numRows; i++ {
		for j := 0; j < column; j++ {
			array[i][j] = '*'
		}
	}
	j, item := 0, 0
	//每列遍历
	for i := 0; i < column; {
		if item == len(s) {
			break
		}
		if row==1 || i%(row-1) == 0 {
			for t := 0; t < row; t++ {
				j = t
				array[j][i] = s[item]
				item++
				if item == len(s) {
					break
				}
			}
			i++
		} else {
			for t := j-1; t > 0; t-- {
				j = t
				array[j][i] = s[item]
				i = i + 1
				item++
				if item == len(s) {
					break
				}
			}
		}
	}
	for i:=0;i<row;i++{
		for j:=0;j<column;j++{
			if array[i][j]!='*'{
				res += string(array[i][j])
			}
		}
	}
	//fmt.Println(res)
	return res
}

func main() {
	s := "A"
	convert(s, 1)

}
