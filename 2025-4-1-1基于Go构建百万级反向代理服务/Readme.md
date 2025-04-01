## 基于Go构建百万级反向代理服务
在现今Web开发中，高效安全地管理海量流量是系统架构设计的核心命题。反向代理作为客户端与后端服务之间的智能调度器，已成为应对高并发场景的利器。

URL:https://mp.weixin.qq.com/s/my1gsoI8pPRm4xaZonJ18Q
### 1. 反向代理
反向代理是一种位于服务器端的代理服务器，它代表后端服务器接收客户端的请求，并将请求转发到内部网络中的实际服务器，最终将服务器的响应返回给客户端。
- 负载均衡：通过轮询/加权算法分发请求至多台后端服务器
- 安全防护盾：隐藏真实服务器IP，有效抵御DDoS攻击（实测可拦截90%的CC攻击）
- 性能加速器：SSL卸载（ssl offloading）使后端CPU负载降低40%，缓存机制使QPS提升3倍
- 数据压缩：减少传输带宽消耗。

```go
package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	target := "http://localhost:8080" // 后端地址
	targetURl, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(targetURl)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// 请求头主注入跟踪
		request.Header.Set("X-Request-ID", generateUUID())
		proxy.ServeHTTP(writer, request)
	})

	log.Printf("代理启动：8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func generateUUID() string {
	return "1"
}
```

### 2. 百万级流量架构设计
**智能负载均衡方案**

问题：当单后端服务无法承受流量压力时，采用轮询策略分发请求，但是简单轮询会导致热点服务器过载

解决方案：
- 动态权重算法：基于实时CPU/内存指标自动调整流量分配
- 一致性哈希：使用stathat.com/c/consistent库实现会话保持或者引入Redis存储session实现粘滞会话
- 健康检查：每30秒TCP探活，自动隔离故障节点

```go
var backends = []string{
    "http://node1.internal:8080",
    "http://node2.internal:8080",
    "http://node3.internal:8080",
}

func getBackend() *url.URL {
// 实际生产环境建议用加权随机算法
    target, _ := url.Parse(backends[atomic.AddUint32(&counter, 1) % 3])
    return target
}

func reverseProxy(w http.ResponseWriter, r *http.Request) {
    proxy := httputil.NewSingleHostReverseProxy(getBackend())
    proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
    // 故障节点自动剔除逻辑
        retryWithNextNode(w, r)
    }
    proxy.ServeHTTP(w, r)
}
```












