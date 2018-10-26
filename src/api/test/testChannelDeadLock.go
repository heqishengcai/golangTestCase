package main

import (
	"fmt"
	"time"
)

// channel 阻塞
// 死锁 fatal error: all goroutines are asleep - deadlock!

/*func main() {
	var c1 chan string = make(chan string)
	func() {
		time.Sleep(time.Second)
		c1 <- "1" //push
	}()
	fmt.Println("c1 is", <-c1) // pop
	//push和pop永远不可能同时发生，会deadlock
}*/

//解决channel 阻塞方法
//1 利用go关键字(并行)，创建一个新的协程，让push和pop不在同一个协程中执行就可以避免死锁
//func main() {
//	var c1 chan string = make(chan string)
//	go func() {
//		time.Sleep(time.Second)
//		c1 <- "2"
//	}()
//	fmt.Println("c1 is", <-c1)
//}

//2 也可以给channel加一个buffer，当buffer没有被塞满的时候，channel是不会阻塞的
func main() {
	var c1 chan string = make(chan string, 1) //buffer
	func() {
		time.Sleep(time.Second)
		c1 <- "1"
	}()
	fmt.Println("c1 is", <-c1)
}
