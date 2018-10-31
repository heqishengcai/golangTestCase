package main

import "fmt"

func main() {
	//test6()
	test7()
}

func test6() {
	var i interface{} = "kk"
	j := i.(int)
	fmt.Printf("%T->%d\n", j, j)
}

func test7() {
	var i interface{} = "123"
	j, b := i.(int) // assert 断言
	if b {
		fmt.Printf("%T->%d\n", j, j)
	} else {
		fmt.Println("类型不匹配")
	}
}
