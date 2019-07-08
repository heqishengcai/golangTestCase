package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		time.Sleep(1 * time.Hour)
	}()
	c := make(chan int)
	go func() {
		for i := 0; i < 10; i = i + 1 {
			c <- i
		}
		close(c) //注释后，程序会一直阻塞
	}()
	for i := range c { //一直迭代直到channel被关闭
		fmt.Println(i)
	}
	fmt.Println("Finished")
}
