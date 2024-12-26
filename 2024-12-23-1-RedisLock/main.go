package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"sync"
	"sync/atomic"
	"time"
)

// 适合生产环境下的分布式锁

// RedNode 结构体表示一个 Redis 节点
type RedNode struct {
	Client *redis.Client
}

// RedLock 结构体实现分布式锁
type RedLock struct {
	LocalHost   []string      // 所有节点的地址
	CanConnect  []string      // 可以沟通的节点
	Nodes       []*RedNode    // 节点的 Redis 客户端
	Quorum      int           // 需要的最小成功节点数量
	SuccessNode int           // 当前可用的 Redis 节点数量
	TTL         time.Duration // 锁的超时时间
	PTL         time.Duration // 锁的续约时间
}

// NewRedLock 创建多个 Redis 客户端
func NewRedLock(ctx context.Context, localhost []string, ttl time.Duration, ptl time.Duration) *RedLock {
	var Nodes []*RedNode
	var canConnect []string
	success := 0 // 记录当前可沟通的节点数量

	for _, host := range localhost {
		client := redis.NewClient(&redis.Options{
			Addr:         host,
			DialTimeout:  3 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		})
		// 测试 Redis 节点是否能连接
		pong, err := client.Ping(ctx).Result()
		if err != nil {
			fmt.Printf("Redis node %s connection failed: %v\n", host, err)
			continue
		}
		fmt.Printf("Redis node %s connected: %s\n", host, pong)

		success++
		Nodes = append(Nodes, &RedNode{Client: client})
		canConnect = append(canConnect, host)
	}

	return &RedLock{
		LocalHost:   localhost,
		CanConnect:  canConnect,
		Nodes:       Nodes,
		Quorum:      len(Nodes)/2 + 1,
		SuccessNode: success,
		TTL:         ttl,
		PTL:         ptl,
	}
}

// Lock 获取分布式锁
func (rl *RedLock) Lock(ctx context.Context, key string) (string, error) {
	// 获取当前客户端锁使用唯一id
	LockID := uuid.New().String()
	wg := sync.WaitGroup{}
	var successCount int32

	for i, node := range rl.Nodes {
		wg.Add(1)
		go func(node *RedNode, i int) {
			defer wg.Done()
			success, err := node.Client.SetNX(ctx, key, LockID, rl.TTL).Result()

			if err == nil && success {
				atomic.AddInt32(&successCount, 1)
				fmt.Printf("Lock acquired on node: %v\n", rl.CanConnect[i])
			} else {
				fmt.Printf("Failed to acquire lock on node: %v, error: %v\n", rl.CanConnect[i], err)
			}
		}(node, i)
	}
	wg.Wait()

	if atomic.LoadInt32(&successCount) >= int32(rl.Quorum) {
		return LockID, nil
	}

	// 清理已加锁的节点
	for _, node := range rl.Nodes {
		result, err := node.Client.Eval(ctx, `
			if redis.call("GET", KEYS[1]) == ARGV[1] then
				return redis.call("DEL", KEYS[1])
			else
				return 0
			end`, []string{key}, LockID).Result()

		if err != nil || result.(int64) != 1 {
			fmt.Printf("Failed to clean up lock on node: %s, error: %v\n", node.Client.Options().Addr, err)
		}
	}

	return "", fmt.Errorf("failed to acquire lock")
}

// Unlock 释放分布式锁
func (rl *RedLock) Unlock(ctx context.Context, key string, lockID string) error {
	var wg sync.WaitGroup

	for _, node := range rl.Nodes {
		wg.Add(1)
		go func(node *RedNode) {
			defer wg.Done()
			result, err := node.Client.Eval(ctx, `
				if redis.call("GET", KEYS[1]) == ARGV[1] then
					return redis.call("DEL", KEYS[1])
				else
					return 0
				end`, []string{key}, lockID).Result()

			if err != nil {
				fmt.Printf("Unlock failed on node: %s, error: %v\n", node.Client.Options().Addr, err)
			} else if result.(int64) != 1 {
				fmt.Printf("Unlock not successful on node: %s\n", node.Client.Options().Addr)
			}
		}(node)
	}
	wg.Wait()
	return nil
}

// RenewLock 续约分布式锁
func (rl *RedLock) RenewLock(ctx context.Context, key string, lockId string) error {
	var wg sync.WaitGroup
	var successCount int32

	for _, node := range rl.Nodes {
		wg.Add(1)
		go func(node *RedNode) {
			defer wg.Done()
			result, err := node.Client.Eval(ctx, `
				if redis.call("GET", KEYS[1]) == ARGV[1] then
					return redis.call("PEXPIRE", KEYS[1], ARGV[2])
				else
					return 0
				end`, []string{key}, lockId, fmt.Sprintf("%d", rl.TTL.Milliseconds())).Result()

			if err == nil && result.(int64) == 1 {
				atomic.AddInt32(&successCount, 1)
			} else if err != nil {
				fmt.Printf("Renew failed on node: %s, error: %v\n", node.Client.Options().Addr, err)
			}
		}(node)
	}
	wg.Wait()

	if atomic.LoadInt32(&successCount) >= int32(rl.Quorum) {
		fmt.Printf("Renew lock success: %v\n", successCount)
		return nil
	}
	return fmt.Errorf("failed to renew lock")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	redisArrs := []string{
		"localhost:6379",
		//"localhost:6381",
		//"localhost:6382",
	}
	redLock := NewRedLock(ctx, redisArrs, 10*time.Second, 10*time.Second/3)

	key := "my-lock"
	lockID, err := redLock.Lock(ctx, key)
	if err != nil {
		fmt.Println("Failed to acquire lock:", err)
		return
	}
	fmt.Printf("Lock success: %v\n", lockID)

	// 启动 Goroutine 进行续约
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Renew Goroutine exited")
				return
			case <-time.After(redLock.PTL):
				err := redLock.RenewLock(ctx, key, lockID)
				if err != nil {
					fmt.Println("Failed to renew lock:", err)
					cancel() // 取消上下文，停止程序
					return
				}
				fmt.Println("Lock renewed successfully.")
			}
		}
	}()

	// 模拟业务逻辑
	time.Sleep(15 * time.Second)

	// 释放锁
	err = redLock.Unlock(ctx, key, lockID)
	if err != nil {
		fmt.Println("Failed to unlock:", err)
		return
	}
	fmt.Printf("Unlock success: %v\n", lockID)
}
