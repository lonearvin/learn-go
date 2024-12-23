package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Println(err)
		}
	}() // 关闭连接
	inputReader := bufio.NewReader(os.Stdin)
	for {
		// 输入用户读的内容
		input, _ := inputReader.ReadString('\n')
		inputInfo := strings.Trim(input, "\r\n")
		if strings.ToUpper(inputInfo) == "Q" {
			// 如果这个字符是Q的话就退出
			return
		}
		// 发送数据
		_, err = conn.Write([]byte(inputInfo))
		if err != nil {
			return
		}
		// 接收数据
		buf := make([]byte, 1024)
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(buf[:n]))
	}
}
