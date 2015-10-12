package memcached

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

func (c *Connection) Send(cmd ...string) (resp []byte, err error) {
	if c.Conn == nil {
		targetTCPAddress, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))
		if err != nil {
			return nil, err
		}
		conn, err := net.DialTCP("tcp", nil, targetTCPAddress)
		if err != nil {
			return nil, err
		}
		err = conn.SetKeepAlive(true)
		if err != nil {
			return nil, err
		}
		c.Conn = conn
	}

	// 写
	for _, cmdPart := range cmd {
		cmdPartBytes := []byte(cmdPart)
		// cmdLength := len(cmdBytes)
		_, err = c.Conn.Write(cmdPartBytes)
		if err != nil {
			return nil, err
		}
	}
	// 读
	var respBuffer bytes.Buffer
	for {
		oneRead := make([]byte, 1024)
		readLength, readErr := c.Conn.Read(oneRead)
		_, bufferWriteErr := respBuffer.Write(oneRead)
		if readLength < len(oneRead) {
			return respBuffer.Bytes(), nil
		}

		if readErr != nil {
			return respBuffer.Bytes(), readErr
		}
		if bufferWriteErr != nil {
			return respBuffer.Bytes(), bufferWriteErr
		}
	}
}
