package main

import (
	"fmt"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	// Goroutine 1: 写入数据到 ch1
	go func() {
		for i := 0; i < 10; i++ {
			ch1 <- i // 将数字写入 ch1
		}
		close(ch1) // 关闭 ch1，通知读取方没有更多数据
	}()

	// Goroutine 2: 从 ch1 读取数据并写入到 ch2
	go func() {
		for i := range ch1 { // 从 ch1 读取数据
			ch2 <- i // 写入 ch2
		}
		close(ch2) // 关闭 ch2，通知读取方没有更多数据
	}()

	// 主函数从 ch2 读取并打印
	for i := range ch2 { // 从 ch2 读取数据
		fmt.Println(i) // 打印数据
	}
}
