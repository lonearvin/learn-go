package main

import (
	"bufio"
	"fmt"
	"net"
)

// 服务端代码

// 处理函数
func process(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection")
		}
	}(conn) // 处理完数据关闭链接
	for {
		reader := bufio.NewReader(conn)
		var buf [1024]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Println("read error:", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("recvStr:", recvStr)
		conn.Write([]byte(recvStr)) // 发回去数据
	}
}
func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		go process(conn)
	}
}
