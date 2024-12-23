package main

import (
	"bufio"
	"fmt"
	"io"
	"learngo/2024-12-4-10-protoEncodeDecode"
	"net"
)

/*
粘包解决方案测试
*/
func process(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Print("close conn error:", err)
		}
		fmt.Println("close conn success")
	}(conn)
	reader := bufio.NewReader(conn)
	for {
		// 读取解码
		msg, err := common.Decode(reader)
		if err != nil {
			if err == io.EOF {
				fmt.Println("client close..")
				break
			}
			fmt.Println("decode error:", err)
			break
		}
		fmt.Println("client send data:", msg)
	}
}

func main() {
	// 监听
	listener, err := net.Listen("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}
	// 关闭
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("listener close error:", err)
		}
	}(listener)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		go process(conn)
	}
}
