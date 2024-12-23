package main

import (
	"fmt"
	common "learngo/2024-12-4-10-protoEncodeDecode"
	"net"
)

func main() {

	// 创建连接句柄
	conn, err := net.Dial("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("dial failed, err:", err)
		return
	}
	// 关闭
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("close conn failed, err:", err)
		}
	}(conn)

	for i := 0; i < 20; i++ {
		msg := "hello world"
		data, err := common.Encode(msg)
		if err != nil {
			fmt.Println("encode failed, err:", err)
		}
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("write failed, err:", err)
			return
		}
	}
}
