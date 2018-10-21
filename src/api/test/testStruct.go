package main

import (
	"fmt"
	"reflect"
)

type Skills []string
type Human struct {
	name   string
	age    int
	weight int
}
type Student struct {
	Human      // 匿名字段，struct
	Skills     // 匿名字段，自定义的类型string slice
	int        // 内置类型作为匿名字段
	speciality string
}

func SmartPrint(i interface{}) {
	var kv = make(map[string]interface{})
	vValue := reflect.ValueOf(i)
	vType := reflect.TypeOf(i)
	for i := 0; i < vValue.NumField(); i++ {
		kv[vType.Field(i).Name] = vValue.Field(i)
	}
	fmt.Println("获取到数据:")
	for k, v := range kv {
		fmt.Print(k)
		fmt.Print(":")
		fmt.Print(v)
		fmt.Println()
	}
}

func main() {
	// 初始化学生Jane
	jane := Student{Human: Human{"Jane", 35, 100}, speciality: "Biology"}
	jane.Skills = append(jane.Skills, "physics", "golang")
	// 现在我们来访问相应的字段
	fmt.Println("jane: ", jane)
	SmartPrint(jane)
	//fmt.Println("jane: ", SmartPrint(jane))

	//fmt.Println("Her name is ", jane.name)
	//fmt.Println("Her age is ", jane.age)
	//fmt.Println("Her weight is ", jane.weight)
	//fmt.Println("Her speciality is ", jane.speciality)
	//// 我们来修改他的skill技能字段
	//jane.Skills = []string{"anatomy"}
	//fmt.Println("Her skills are ", jane.Skills)
	//fmt.Println("She acquired two new ones ")
	//jane.Skills = append(jane.Skills, "physics", "golang")
	//fmt.Println("Her skills now are ", jane.Skills)
	//// 修改匿名内置类型字段
	//jane.int = 3
	//fmt.Println("Her preferred number is", jane.int)
}
