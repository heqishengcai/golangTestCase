package main

import "fmt"

type Noddles interface { //面条接口
	Description() string
	Price() float32
}

/**
拉面
*/

type Ramen struct {
	name  string
	price float32
}

func (p Ramen) Description() string {
	return p.name
}

func (p Ramen) Price() float32 {
	return p.price
}

/**
加蛋的拉面
*/

type Egg struct {
	noddles Noddles
	name    string
	price   float32
}

func (p Egg) SetNoddles(noddles Noddles) {
	p.noddles = noddles
}

func (p Egg) Description() string {
	return p.noddles.Description() + "+" + p.name
}

func (p Egg) Price() float32 {
	return p.noddles.Price() + p.price
}

/**
香肠的鸡蛋拉面
*/
type Sausage struct {
	noddles Noddles
	name    string
	price   float32
}

func (p Sausage) SetNoddles(noddles Noddles) {
	p.noddles = noddles
}

func (p Sausage) Description() string {
	return p.noddles.Description() + "+" + p.name
}

func (p Sausage) Price() float32 {
	return p.noddles.Price() + p.price
}

func main() {
	hr := "-------------------"
	ramen := Ramen{name: "ramen", price: 10}                    //面
	egg := Egg{noddles: ramen, name: "egg", price: 2}           //面+蛋
	sausage := Sausage{noddles: egg, name: "sausage", price: 3} //面+蛋+肠
	egg2 := Egg{noddles: ramen, name: "egg", price: 4}          //面+双蛋

	fmt.Println(ramen.Description())
	fmt.Println(ramen.Price())
	fmt.Println(hr)
	fmt.Println(egg.Description())
	fmt.Println(egg.Price())
	fmt.Println(hr)
	fmt.Println(sausage.Description())
	fmt.Println(sausage.Price())
	fmt.Println(hr)
	fmt.Println(egg2.Description())
	fmt.Println(egg2.Price())

}

/**
https://github.com/qibin0506/go-designpattern
java/io guofu 可以讲下
*/
