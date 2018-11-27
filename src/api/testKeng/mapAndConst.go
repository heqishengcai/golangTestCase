package main

import (
	"fmt"
	"reflect"
)

const t = 1

func main() {

	sliceA := make(map[string]interface{})
	var T uint8
	var TT uint8 = 1
	T = t
	sliceA["T"] = TT

	if sliceA["T"] == T {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}

	if sliceA["T"] == TT {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}

	if sliceA["T"] == t {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}

	if T == TT {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
	fmt.Printf("test type of t : %+v \n", typeof(t))
	fmt.Printf("test type of T : %+v \n", typeof(T))
	fmt.Printf("test type of TT : %+v \n", typeof(TT))
	fmt.Printf("test type of sliceA['T'] : %+v \n", typeof(sliceA["T"]))
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}
