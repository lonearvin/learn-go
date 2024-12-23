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

// NewRedLock 得到DB的默认值
func NewRedLock(LocalHost []string, ttl time.Duration) *RedisLock {
	Node := make([]*RedisNode, len(LocalHost))
	for i, addr := range LocalHost {
		client := redis.NewClient(&redis.Options{
			Addr: addr,
		})
		Node[i] = &RedisNode{Client: client}
	}
	return &RedisLock{LocalHost: LocalHost,
		Nodes:  Node,
		Quorum: (len(Node) / 2) + 1,
		TTL:    ttl,
	}
}

// Lock 尝试获取分布式锁
func (r *RedisLock) Lock(ctx context.Context, key string) (string, error) {
	lockID := uuid.New().String()
	//startTime := time.Now()

	successCount := 0
	mu := sync.Mutex{}

	// 遍历所有节点并尝试加锁
	for i, node := range r.Nodes {
		go func(node *RedisNode, i int) {
			// 尝试加锁
			success, err := node.Client.SetNX(ctx, key, lockID, r.TTL).Result()
			if err == nil && success {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
			if err != nil {
				fmt.Printf("localhost:%v,err:%v\n", r.LocalHost[i], err)
			} else {
				fmt.Printf("localhost:%v lock successful.\n", r.LocalHost[i])
			}
		}(node, i)
	}

	// 等待所有加锁请求完毕
	time.Sleep(10 * time.Millisecond * time.Duration(len(r.Nodes)))
	if successCount >= r.Quorum {
		return lockID, nil
	}
	fmt.Println("successful:", successCount)
	// 如果未达到预期，释放所有锁
	for _, node := range r.Nodes {
		node.Client.Eval(ctx, `
			if redis.call("GET", KEYS[1]) == ARGV[1] then
				return redis.call("DEL", KEYS[1])
			else
				return 0
			end
		`, []string{key}, lockID)
	}
	return "", fmt.Errorf("fail to acquire lock")
}

// Unlock 解锁
func (r *RedisLock) Unlock(ctx context.Context, key, lockID string) error {
	var wg sync.WaitGroup // 保证原子性
	// 遍历所有节点并尝试解锁
	for _, node := range r.Nodes {
		wg.Add(1)
		go func(node *RedisNode) {
			defer wg.Done()
			node.Client.Eval(ctx, `
				if redis.call("GET", KEYS[1]) == ARGV[1] then
					rerturn redis.call("DEL", KEYS[1])
				else
					return 0
				end
			`, []string{key}, lockID)
		}(node)
	}
	wg.Wait()
	return nil
}
