package main

import (
	"fmt"
	"net"
)

func main() {
	// 创建对象
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(12, 0, 0, 1),
		Port: 8081,
	})
	// 连接判断
	if err != nil {
		fmt.Println("net.Conn err:", err)
		return
	}
	// 关闭连接
	defer func(socket *net.UDPConn) {
		err := socket.Close()
		if err != nil {
			fmt.Println("socket.Close err:", err)
		}
	}(socket)
	// 发送数据
	sendData := []byte("hello world")
	_, err = socket.Write(sendData)
	if err != nil {
		fmt.Println("socket.Write err:", err)
		return
	}
	// 接收数据
	data := make([]byte, 4090)
	n, remoteAddr, err := socket.ReadFromUDP(data)
	if err != nil {
		fmt.Println("socket.ReadFromUDP err:", err)
	}
	fmt.Printf("receive:%v, addr:%v, count:%v\n\n", string(data[:n]), remoteAddr, n)
}
