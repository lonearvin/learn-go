package main

import "fmt"

var (
	func1 = func(a int, b int) int {
		return a * b
	}
)

// 入口函数，整个程序的入口，
func main() {
	//fmt.Print("hello world\n")

	res1 := func(a int, b int) int {
		return a + b
	}(10, 20)

	fmt.Printf("res1: %d\n", res1)
	a := func(a int, b int) int {
		return a + b
	}
	res2 := a(10, 20)
	fmt.Printf("res2: %d\n", res2)

	// 使用全局匿名函数的使用
	res3 := func1(10, 20)
	fmt.Printf("res3: %d\n", res3)
}

// 初始化函数，golang 每个包的引用会优先调用
func init() {

}
