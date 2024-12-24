package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisNode struct {
	Client *redis.Client
}

type RedisLock struct {
	LocalHost []string
	Nodes     []*RedisNode
	Quorum    int           // 至少需要成功加锁的节点
	TTL       time.Duration // 锁的有效时间
}

// NewRedLock 创建 RedisLock 实例
func NewRedLock(LocalHost []string, ttl time.Duration) *RedisLock {
	nodes := make([]*RedisNode, len(LocalHost))
	for i, addr := range LocalHost {
		client := redis.NewClient(&redis.Options{
			Addr: addr,
		})
		nodes[i] = &RedisNode{Client: client}
	}
	return &RedisLock{
		LocalHost: LocalHost,
		Nodes:     nodes,
		Quorum:    (len(nodes) / 2) + 1,
		TTL:       ttl, // 此为key，value的超时时间
	}
}

// Lock 尝试获取分布式锁
func (r *RedisLock) Lock(ctx context.Context, key string) (string, error) {
	lockID := uuid.New().String() // 生成唯一锁 ID
	successCount := 0
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}

	// 遍历所有节点尝试加锁
	for i, node := range r.Nodes {
		wg.Add(1)
		go func(node *RedisNode, i int) {
			defer wg.Done()
			success, err := node.Client.SetNX(ctx, key, lockID, r.TTL).Result()
			mu.Lock()
			defer mu.Unlock()
			if err == nil && success {
				successCount++
				fmt.Printf("localhost:%v lock successful.\n", r.LocalHost[i])
			} else if err != nil {
				fmt.Printf("localhost:%v error: %v\n", r.LocalHost[i], err)
			}
		}(node, i)
	}

	// 等待所有加锁请求完成
	wg.Wait()

	// 如果成功的节点数大于等于 Quorum，则加锁成功
	if successCount >= r.Quorum {
		return lockID, nil
	}

	fmt.Printf("successful locks: %d, required quorum: %d\n", successCount, r.Quorum)

	// 如果未达到 Quorum，清理已加的锁
	for _, node := range r.Nodes {
		node.Client.Eval(ctx, `
			if redis.call("GET", KEYS[1]) == ARGV[1] then
				return redis.call("DEL", KEYS[1])
			else
				return 0
			end
		`, []string{key}, lockID)
	}

	return "", fmt.Errorf("failed to acquire lock")
}

// Unlock 解锁
func (r *RedisLock) Unlock(ctx context.Context, key, lockID string) error {
	var wg sync.WaitGroup
	for _, node := range r.Nodes {
		wg.Add(1)
		go func(node *RedisNode) {
			defer wg.Done()
			node.Client.Eval(ctx, `
				if redis.call("GET", KEYS[1]) == ARGV[1] then
					return redis.call("DEL", KEYS[1])
				else
					return 0
				end
			`, []string{key}, lockID)
		}(node)
	}
	wg.Wait()
	return nil
}

// RenewLock 续期锁
func (r *RedisLock) RenewLock(ctx context.Context, key, lockID string) error {
	var wg sync.WaitGroup
	successCount := 0
	mu := sync.Mutex{}

	for _, node := range r.Nodes {
		wg.Add(1)
		go func(node *RedisNode) {
			defer wg.Done()
			success, err := node.Client.Eval(ctx, `
				if redis.call("GET", KEYS[1]) == ARGV[1] then
					return redis.call("PEXPIRE", KEYS[1], ARGV[2])
				else
					return 0
				end
			`, []string{key}, lockID, int(r.TTL.Milliseconds())).Result()
			mu.Lock()
			defer mu.Unlock()
			if err == nil && success.(int64) == 1 {
				successCount++
			}
		}(node)
	}
	wg.Wait()

	// 如果未达到法定人数，续期失败
	if successCount < r.Quorum {
		return fmt.Errorf("failed to renew lock")
	}
	return nil
}

// Main 用于测试分布式锁
func main() {
	ctx := context.Background()
	redisAddrs := []string{
		"localhost:6380",
		"localhost:6381",
		"localhost:6382",
	}
	// 创建新的锁
	redLock := NewRedLock(redisAddrs, 10*time.Second)
	// 创建锁的key
	key := "my-lock"

	// 尝试加锁
	lockID, err := redLock.Lock(ctx, key)
	if err != nil {
		fmt.Println("加锁失败:", err)
		return
	}
	fmt.Println("加锁成功, LockID:", lockID)

	// 续期锁（防止锁过期）
	go func() {
		for {
			time.Sleep(5 * time.Second)
			err := redLock.RenewLock(ctx, key, lockID)
			if err != nil {
				fmt.Println("续期失败:", err)
				return
			}
			fmt.Println("锁续期成功.")
		}
	}()

	// 模拟业务逻辑
	time.Sleep(15 * time.Second)

	// 解锁
	err = redLock.Unlock(ctx, key, lockID)
	if err != nil {
		fmt.Println("解锁失败:", err)
		return
	}
	fmt.Println("解锁成功.")
}
