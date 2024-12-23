package main

import (
	"fmt"
	"runtime"
)

func main() {
	go func(s string) {
		for i := 0; i < 3; i++ {
			fmt.Println(s)
		}
	}("world")
	// 主协程

	for i := 0; i < 10; i++ {
		runtime.Gosched() // 切一下
		fmt.Println("hello")
	}
}
