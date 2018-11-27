package main

import "fmt"

type newEr interface {
	New()
}

type testInterface interface {
	newEr

	Done() <-chan struct{}
}

type kkTest struct {
	testInterface
}

func NewTest() newEr {

	return kkTest{}

}

func main() {

	kk := NewTest()

	// assert
	i, ok := kk.(testInterface)

	fmt.Println(i, ok)

	//ch := i.Done()

	//fmt.Println(ch)

}
