package main

import (
	"fmt"
	"net"
)

func main() {
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 8081,
	})
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}
	defer func(listen *net.UDPConn) {
		err := listen.Close()
		if err != nil {
			fmt.Println("listen error:", err)
		}
	}(listen)

	for {
		// 接收数据
		var buf [1024]byte
		n, add, err := listen.ReadFromUDP(buf[:])
		if err != nil {
			fmt.Println("read from udp error:", err)
			continue
		}
		fmt.Println("read:", string(buf[:n]))

		// 发回数据
		_, err = listen.WriteToUDP([]byte("hello"), add)

		if err != nil {
			fmt.Println("write to udp error:", err)
			continue
		}
	}
}
