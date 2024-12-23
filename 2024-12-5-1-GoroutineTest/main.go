package main

import (
	"fmt"
	"time"
)

// 测试主协程退出了，其他还在运行吗

func main() {
	go func() {
		i := 0
		for {
			i++
			fmt.Printf("new goroutine: i=%d\n", i)
			time.Sleep(time.Second)
		}
	}()

	i := 0
	for {
		i++
		fmt.Printf("main goroutline :i=%d\n", i)
		time.Sleep(time.Second)
		if i == 2 {
			break
		}
	}
}

// 结论当然不会了
