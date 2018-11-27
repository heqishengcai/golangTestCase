package main

import "fmt"

type classmate struct {
	class uint32
	grade uint32
	age   uint32
	name  string
}

type classinfo struct {
	total uint32
	grade uint32
}

type info interface {
	coutInfo()
	changeInfo(grade uint32)
}

func (this *classmate) coutInfo() {
	fmt.Println("info is ", this)
}

func (this *classmate) changeInfo(grade uint32) {
	this.grade = grade
	fmt.Println("info is ", this)
}

func (this *classinfo) coutInfo() {
	fmt.Println("info is ", this)
}

func (this *classinfo) changeInfo(grade uint32) {
	this.grade = grade
	fmt.Println("info is ", this)
}

func interTest(test info, grade uint32) {
	test.coutInfo()
	test.changeInfo(grade)
}

/**
相信这个例子可以帮助我们更好的理解函数interTest的第一个输入参数并没有要求参数的具体类型，而是一个接口类型。
*/
func main() {
	mate := &classmate{
		class: 1,
		grade: 1,
		age:   6,
		name:  "jim",
	}
	info := &classinfo{
		total: 30,
		grade: 5,
	}
	interTest(mate, 3)
	interTest(info, 8)
}
