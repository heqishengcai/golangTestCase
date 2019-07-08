package main

import "fmt"

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
			//fmt.Println(x, y, x+y)
			//case c <- y:

		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}
func main() {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			res := <-c
			fmt.Println(res)
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}
