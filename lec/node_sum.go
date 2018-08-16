package main

import "fmt"

/**
* Definition for singly-linked list.
* type ListNode struct {
*     Val int
*     Next *ListNode
* }
 */

type ListNode struct {
	Val  int
	Next *ListNode
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	res := new(ListNode)
	v1 := l1.Val
	v2 := l2.Val
	//for carry
	carry_node := new(ListNode)
	carry_node.Val = 1
	is_carry := false
	if v1+v2 >= 10 {
		res.Val = (v1 + v2) % 10
		is_carry = true
	} else {
		res.Val = v1 + v2
	}

	// if only one digit
	if l1.Next == nil && l2.Next == nil {
		if is_carry == false {
			return res
		} else {

			res.Next = carry_node
			return res
		}

	}
	//store last node
	last := new(ListNode)
	res.Next = last
	for n1, n2 := l1.Next, l2.Next; ; {
		n1v, n2v := 0, 0
		if n1 == nil {
			n2v = n2.Val
			n2 = n2.Next
		} else if n2 == nil {
			n1v = n1.Val
			n1 = n1.Next
		} else {
			n1v, n2v = n1.Val, n2.Val
			n1 = n1.Next
			n2 = n2.Next
		}
		var cur_val int
		if is_carry {
			cur_val = n1v + n2v + 1

		} else {
			cur_val = n1v + n2v
		}
		if cur_val >= 10 {
			cur_val = cur_val % 10
			is_carry = true
		} else {
			is_carry = false
		}
		last.Val = cur_val
		if n1 == nil && n2 == nil {
			if is_carry {
				last.Next = carry_node
				break
			} else {
				last = nil
				break
			}

		}
		next_curl := new(ListNode)
		last.Next = next_curl
		last = last.Next
	}
	return res
}

func PrintLN(l *ListNode) {
	rev_node := new(ListNode)

	rev_node.Val = l.Val
	rev_node.Next = nil
	for n := l.Next; n != nil; n = n.Next {
		temp_node := new(ListNode)
		temp_node.Val = n.Val
		temp_node.Next = rev_node
		rev_node = temp_node
	}

	for n := rev_node; n != nil; n = n.Next {
		fmt.Print(n.Val, "")
	}

	fmt.Printf("\n")
}

func main() {
	//l1 := &ListNode{2, &ListNode{4, &ListNode{3, nil}}}
	//l2 := &ListNode{5, &ListNode{6, &ListNode{4, nil}}}
	//
	//l3 := &ListNode{0, nil}
	//l4 := &ListNode{0, nil}

	l5 := &ListNode{1, nil}
	l6 := &ListNode{9, &ListNode{9, nil}}
	//t1 := &ListNode{5,&ListNode{4,&ListNode{1,nil}}}
	//PrintLN(l1)
	//PrintLN(t1)
	//PrintLN(addTwoNumbers(l1, l2))
	//PrintLN(addTwoNumbers(l3, l4))
	PrintLN(addTwoNumbers(l5, l6))

}
