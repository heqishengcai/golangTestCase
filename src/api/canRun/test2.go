package main

import "fmt"

//func main() {
//
//	// Declare variable of type int with a value of 10.
//	count := 10
//
//	// Display the "value of" and "address of" count.
//	println("count:\tValue Of[", count, "]\tAddr Of[", &count, "]")
//
//	// Pass the "value of" the count.
//	increment(count)
//
//	println("count:\tValue Of[", count, "]\tAddr Of[", &count, "]")
//}
//
////go:noinline
//func increment(inc int) {
//
//	// Increment the "value of" inc.
//	inc++
//	println("inc:\tValue Of[", inc, "]\tAddr Of[", &inc, "]")
//}

func main() {
	a := 0
	times := 10000
	c := make(chan bool)

	for i := 0; i < times; i++ {
		go func() {
			a++
			c <- true
		}()
	}

	for i := 0; i < times; i++ {
		<-c
	}
	fmt.Printf("a = %d \n ", a)
}
