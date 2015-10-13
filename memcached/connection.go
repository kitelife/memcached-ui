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

func (c *Connection) Open() error {
	targetTCPAddress, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, targetTCPAddress)
	if err != nil {
		return err
	}
	err = conn.SetKeepAlive(true)
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}

func (c *Connection) Send(cmd ...string) (resp []byte, err error) {
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
		_, bufferWriteErr := respBuffer.Write(oneRead[:readLength])
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
