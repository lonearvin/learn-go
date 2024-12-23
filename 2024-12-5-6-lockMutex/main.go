package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	x      int64
	wg     sync.WaitGroup
	lock   sync.Mutex
	rwLock sync.RWMutex
)

func write() {
	defer wg.Done()
	rwLock.Lock() // 加写锁
	x = x + 1
	time.Sleep(10 * time.Millisecond)
	rwLock.Unlock()
}

func read() {
	defer wg.Done()
	rwLock.RLock()
	time.Sleep(time.Millisecond)
	rwLock.RUnlock()
}

func main() {
	start := time.Now()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go write()
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go read()
	}

	wg.Wait()
	end := time.Now()
	fmt.Println(end.Sub(start))
}
