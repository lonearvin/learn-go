package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ctx = context.Background()

func main() {
	// 配置连接 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "10.16.64.250:6379", // Redis 地址
		Password: "",                  // 无密码设置
		DB:       0,                   // 使用默认数据库
	})

	// 尝试 PING 来检查连接
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to redis:", err)
		return
	}

	fmt.Println("Redis connection successful")

	// 示例：设置一个键
	err = rdb.Set(ctx, "key", "value", 10*time.Second).Err()
	if err != nil {
		fmt.Println("Failed to set key:", err)
		return
	}

	// 示例：获取一个键
	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		fmt.Println("Failed to get key:", err)
		return
	}

	fmt.Println("key:", val)

	// 示例：删除一个键
	err = rdb.Del(ctx, "key").Err()
	if err != nil {
		fmt.Println("Failed to delete key:", err)
		return
	}
}
