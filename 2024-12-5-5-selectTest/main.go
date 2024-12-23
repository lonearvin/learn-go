package main

import (
	"fmt"
	"time"
)

func test1(ch chan string) {
	time.Sleep(time.Second * 5)
	ch <- "test1"
	close(ch)
}

func test2(ch chan string) {
	time.Sleep(time.Second * 5)
	ch <- "test2"
	close(ch)
}

func main() {
	var ch1 chan string
	var ch2 chan string
	ch1 = make(chan string)
	ch2 = make(chan string)
	go test1(ch1)
	go test2(ch2)
	// 使用select监控
	for {
		select {
		case s1, ok := <-ch1:
			if ok {
				fmt.Println(s1)
			} else {
				ch1 = nil
			}
			//fmt.Println(s1)
		case s2, ok := <-ch2:
			if ok {
				fmt.Println(s2)
			} else {
				ch2 = nil
			}
		}
		if ch1 == nil && ch2 == nil {
			break
		}
	}
	fmt.Println("All tasks completed.")
}
