package main

import "fmt"

type options struct {
	a int64
	b string
	c map[int]string
}

type ServerOption func(*options)

func NewOption(opt ...ServerOption) *options {
	r := new(options)
	for _, o := range opt {
		o(r)
	}
	return r
}

func WriteA(s int64) ServerOption {
	return func(o *options) {
		o.a = s
	}
}

func WriteB(s string) ServerOption {
	return func(o *options) {
		o.b = s
	}
}

func WriteC(s map[int]string) ServerOption {
	return func(o *options) {
		o.c = s
	}
}

func main() {
	opt1 := WriteA(int64(1))
	opt2 := WriteB("test")
	mapC := make(map[int]string, 0)
	mapC[1] = "aaa"
	mapC[10] = "bbbb"
	opt3 := WriteC(mapC)

	op := NewOption(opt1, opt2, opt3)

	//op := new(options)
	//op.WriteA(int64(1)).WriteB("test").WriteC(make(map[int]string, 0))

	fmt.Println(op.a, op.b, op.c)
}
