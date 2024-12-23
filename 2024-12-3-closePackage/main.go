package main

import "fmt"

// AddUpper 返回的是一个匿名函数
// 累加器
//func AddUpper() func(int) int {
//	var n int = 10
//	return func(x int) int {
//		n = n + x
//		return n
//	}
//}
//
//func makeSuffix(suffix string) func(string) string {
//	var str string = suffix
//	return func(name string) string {
//		if strings.HasSuffix(name, ".jpg") {
//			return name
//		}
//		return name + str
//	}
//}

func shuffix(base int) (func(i int) int, func(i int) int) {

	add := func(i int) int {
		base += i
		return base
	}

	shuffix_func := func(i int) int {
		base -= i
		return base
	}
	return add, shuffix_func
}

func main() {

	add, shu := shuffix(10)
	fmt.Println(add(20))
	fmt.Println(shu(20))

}
