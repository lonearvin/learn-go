package common

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

func Encode(message string) ([]byte, error) {
	// 读取消息的长度，转转成int32类型(占用4个字节)
	// 编码
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// 写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, []byte(message))

	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

func Decode(reader *bufio.Reader) (string, error) {
	//lengthByte, err := reader.ReadByte()
	// 读取消息长度
	lengthByte, _ := reader.Peek(4) // 取出4个字节
	// 创建这么长的空间
	lengthBuffer := bytes.NewBuffer(lengthByte)
	var length int32
	err := binary.Read(lengthBuffer, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}

	// buffered返回缓冲区中现有的可读取的字节数
	if int32(reader.Buffered()) < length+4 {
		return "", err
	}
	// 读取真正的消息体
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return "", err
	}
	return string(pack[4:]), nil
}
