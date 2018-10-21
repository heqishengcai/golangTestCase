package main

import "fmt"

/*func ReadWrite() bool {
	file.Open("")
	defer file.Close()
	// 做一些工作
	if failureX {
		file.Close()
		return false
	}
	if failureY {
		file.Close()
		return false
	}
	file.Close()
	return true
}
func main() {
	ReadWrite()
}*/

type testInt func(int) bool // 声明了一个函数类型

func isOdd(integer int) bool {
	if integer%2 == 0 {
		return false
	}
	return true
}
func isEven(integer int) bool {
	if integer%2 == 0 {
		return true
	}
	return false
}

// 声明的函数类型在这个地方当做了一个参数
func filter(slice []int, f testInt) []int {
	var result []int
	for _, value := range slice {
		if f(value) {
			result = append(result, value)
		}
	}
	return result
}
func main() {
	defer println("be deferred")
	slice := []int{1, 2, 3, 4, 5, 7}
	fmt.Println("slice = ", slice)
	odd := filter(slice, isOdd) // 函数当做值来传递了
	fmt.Println("Odd elements of slice are: ", odd)
	even := filter(slice, isEven) // 函数当做值来传递了
	fmt.Println("Even elements of slice are: ", even)
}
