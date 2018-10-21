package main

import "fmt"

// 返回a、b中最大值.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func main() {
	x := 3
	y := 4
	z := 5
	max_xy1 := max(x, y) //调用函数max(x, y)
	max_xz := max(x, z)  //调用函数max(x, z)
	fmt.Printf("max(%d, %d) = %d\n", x, y, max_xy)
	fmt.Printf("max(%d, %d) = %d\n", x, z, max_xz)
	fmt.Printf("max(%d, %d) = %d\n", y, z, max(y, z)) // 也可在这直接调用它
}
