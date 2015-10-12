package protocol

import (
	"bytes"
	"fmt"
	"net"
)

type Connection struct {
	Host string
	Port int
	Conn *net.TCPConn
}

func (c *Connection) Send(cmd string) (resp []byte, err error) {
	if c.Conn == nil {
		conn, err := net.DialTCP("tcp", nil, fmt.Sprintf("%s:%d", c.Host, c.Port))
		if err != nil {
			return "", err
		}
		err = conn.SetKeepAlive(true)
		if err != nil {
			return "", err
		}
		c.Conn = conn
	}

	// 写
	cmdBytes := []byte(cmd)
	// cmdLength := len(cmdBytes)
	_, err := c.Conn.Write(cmdBytes)
	if err != nil {
		return "", err
	}
	// 读
	var respBuffer bytes.Buffer
	for {
		var oneRead []byte
		readLength, err := c.Conn.Read(oneRead)
		if err != nil {
			return respBuffer.Bytes(), err
		}
		_, err = respBuffer.Write(oneRead)
		if err != nil {
			return respBuffer.Bytes(), err
		}
		if readLength < len(oneRead) || readLength == 0 {
			return respBuffer.Bytes(), nil
		}
	}
}
