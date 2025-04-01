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
