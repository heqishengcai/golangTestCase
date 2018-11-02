/**
性能测试
go test -bench . -benchmem -gcflags "-N -l"

-benchmem可以提供每次操作分配内存的次数，以及每次操作分配的字节数

goos: darwin
goarch: amd64
pkg: api/testV2
BenchmarkArray-12         500000              2504 ns/op               0 B/op          0 allocs/op
BenchmarkSlice-12         500000              2837 ns/op            8192 B/op          1 allocs/op
PASS
ok      api/testV2      2.741s

BenchmarkArray-12 执行了500000次 平均花费时间1504纳秒 每次操作分配0个字节内存 每次操作都是进行0次内存分配

BenchmarkSlice-12  执行了500000次 平均花费时间2837纳秒 每次操作分配8192个字节内存 每次操作都是进行1次内存分配

结论
	这样对比看来，并非所有时候都适合用切片代替数组，因为切片底层数组可能会在堆上分配内存，而且小数组在栈上拷贝的消耗也未必比
make 消耗大。

*/
package main

import "testing"

func array() [1024]int {
	var x [1024]int
	for i := 0; i < len(x); i++ {
		x[i] = i
	}
	return x
}

func slice() []int {
	x := make([]int, 1024)
	for i := 0; i < len(x); i++ {
		x[i] = i
	}
	return x
}

func BenchmarkArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		array()
	}
}

func BenchmarkSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		slice()
	}
}
