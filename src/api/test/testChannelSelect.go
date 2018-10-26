package main

import (
	"fmt"
	"time"
)

/*
	在select中使用发送操作并且有 default可以确保发送不被阻塞！如果没有case，select就会一直阻塞
*/
/*func main() {
	var c1 chan string = make(chan string)
	var c2 chan string = make(chan string)
	time.Sleep(time.Second)
	select {
	case c := <-c1:
		fmt.Println(c)
	case c := <-c2:
		fmt.Println(c)
		//default:
		//	fmt.Println("After one seconds!")
	}
}*/

func askForC(c chan string) {
	fmt.Println("run in other routine")
}
func askForOther() {
	fmt.Println("run in other")
}
func main() {
	var c1 chan string = make(chan string)
	var c2 chan string = make(chan string)
	go askForC(c1)
	go askForC(c2)
	go askForOther()
	select {
	case c := <-c1:
		fmt.Println(c)
	case c := <-c2:
		fmt.Println(c)
	case <-time.After(time.Second):
		fmt.Println("After one second!")
	}
}
