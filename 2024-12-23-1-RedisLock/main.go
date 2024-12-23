package main

import (
	"context"
	"fmt"
	"time"
)

// 该项目模块是进行分布式锁的测试

func main() {
	ctx := context.Background()
	// 创建RedLock实例

	redisAdds := []string{
		"localhost:6380",
		"localhost:6381",
		"localhost:6382",
	}
	RedLock := NewRedLock(redisAdds, 10*time.Second)

	key := "my-lock"
	// 尝试加锁
	LockID, err := RedLock.Lock(ctx, key)
	if err != nil {
		fmt.Println("加锁失败:", err)
		return
	}
	fmt.Println("加锁成功, LockID:", LockID)
	time.Sleep(5 * time.Second)

	// 解锁
	err = RedLock.Unlock(ctx, key, LockID)
	if err != nil {
		fmt.Println("解锁失败.", err)
	}
	fmt.Println("解锁成功.")
}
