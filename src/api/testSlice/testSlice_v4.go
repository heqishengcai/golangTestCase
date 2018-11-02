package main

import (
	"fmt"
	"unsafe"
)

func main() {

}
s
func printPointer() {
	s := make([]byte, 200)
	ptr := unsafe.Pointer(&s[0])
	fmt.Println("内存地址:", ptr)
	fmt.Printf("内存地址是：%p\n", ptr)
}
