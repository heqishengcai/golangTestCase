/**

Go 中数组赋值和函数传参都是值复制的
用切片传数组参数，既可以达到节约内存的目的，也可以达到合理处理好共享内存的问题。打印结果第二行就是切片，切片的指针和原来数组的指针是不同的
把第一个大数组传递给函数会消耗很多内存，采用切片的方式传参可以避免上述问题。切片是引用传递，所以它们不需要使用额外的内存并且比使用数组更有效率。


*/
package main

import "fmt"

func main() {
	// 测试 slice 传值
	//testMemAddress()
	// 测试 slice 传指针
	testMemAddressPointer()
}

func testMemAddressPointer() {
	arrayA := [2]int{100, 200}
	testArrayPoint(&arrayA) // 1.传数组指针
	arrayB := arrayA[0:]
	testSlicePoint(&arrayB) // 2.传切片
	arrayC := arrayA
	testArrayPoint(&arrayC) // 3.传数组指针
	fmt.Printf("arrayA : %p , %v\n", &arrayA, arrayA)
}

func testArrayPoint(x *[2]int) {
	fmt.Printf("func Array : %p , %v\n", x, *x)
	(*x)[1] += 200
}
func testSlicePoint(x *[]int) {
	fmt.Printf("func slice : %p , %v\n", x, *x)
	(*x)[1] += 1
}

func testMemAddress() {
	arrayA := [2]int{100, 200}
	var arrayB [2]int

	arrayB = arrayA

	fmt.Printf("arrayA : %p , %v\n", &arrayA, arrayA)
	fmt.Printf("arrayB : %p , %v\n", &arrayB, arrayB)

	testArray(arrayA)
}

func testArray(x [2]int) {
	fmt.Printf("func Array : %p , %v\n", &x, x)
}
