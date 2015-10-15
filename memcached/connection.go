package memcached

import (
	"fmt"
	"io"
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
	respBuffer := make([]byte, 1024)
	respLength := 0
	for {
		readLength, readErr := c.Conn.Read(respBuffer[respLength:])
		respLength += readLength
		if (len(respBuffer)-respLength) > 0 || readErr == io.EOF {
			return respBuffer[:respLength], nil
		}
		if readErr != nil {
			return respBuffer[:respLength], readErr
		}
		biggerRespBuffer := make([]byte, respLength*2)
		copy(biggerRespBuffer, respBuffer)
		respBuffer = biggerRespBuffer
	}
}
