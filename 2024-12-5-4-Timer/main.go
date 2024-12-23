package main

import (
	"fmt"
	"time"
)

func main() {
	// Timer 是只触发一次
	//timer1 := time.NewTimer(1*time.Second)
	// 使用 Ticker 而不是 Timer
	ticker := time.NewTicker(1 * time.Second) // 创建一个周期性 Ticker，每 1 秒触发一次
	defer ticker.Stop()                       // 确保主程序退出时停止 Ticker

	i := 0
	go func() {
		for t := range ticker.C {
			i++
			fmt.Println(t)
			if i == 5 {
				ticker.Stop()
				return // 退出协程，避免阻塞
			}
		}
	}()

	// 保持主程序运行
	for {
	}
}
