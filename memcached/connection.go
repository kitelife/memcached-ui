package memcached

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"
)

type Connection struct {
	Host string
	Port int
	Conn *net.TCPConn
}

func (c *Connection) checkError(resp []byte) error {
	if bytes.Compare(resp, []byte("ERROR")) == 0 {
		return errors.New("发生错误：ERROR")
	}
	if bytes.HasPrefix(resp, []byte("CLIENT_ERROR ")) {
		return errors.New(fmt.Sprintf("发生错误：%s", string(resp)))
	}
	if bytes.HasPrefix(resp, []byte("SERVER_ERROR ")) {
		return errors.New(fmt.Sprintf("发生错误：%s", string(resp)))
	}
	return nil
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

func (c *Connection) Send(cmd ...string) (err error) {
	// 写
	for _, cmdPart := range cmd {
		cmdPartBytes := []byte(cmdPart)
		cmdPartLength := len(cmdPartBytes)
		allWritenLength := 0
		for {
			if allWritenLength < cmdPartLength {
				writenLength, err := c.Conn.Write(cmdPartBytes[allWritenLength:])
				if err != nil {
					return err
				}
				allWritenLength += writenLength
				continue
			}
			break
		}
	}
	return nil
}

func (c *Connection) Receive(cmd string) (interface{}, error) {
	byteReader := bufio.NewReader(c.Conn)
	switch {
	case cmd == "get" || cmd == "gets":
		mapper := make(map[string]string)
		for {
			line, err := byteReader.ReadBytes('\n')
			if err != nil {
				return mapper, err
			}
			line = bytes.Trim(line, "\r\n")
			if string(line) == "END" {
				return mapper, nil
			}
			lineParts := bytes.Split(line, []byte(" "))
			if len(lineParts) != 4 || string(lineParts[0]) != "VALUE" {
				err = c.checkError(line)
				if err != nil {
					return mapper, err
				}
				return mapper, errors.New("响应数据格式非法")
			}
			valueLength, err := strconv.Atoi(string(lineParts[3]))
			if err != nil {
				return mapper, err
			}
			// 加上\r\n
			valueLength += 2
			value := make([]byte, valueLength)
			bytePosition := 0
			for bytePosition < valueLength {
				oneByte, err := byteReader.ReadByte()
				if err != nil {
					return mapper, err
				}
				value[bytePosition] = oneByte
				bytePosition += 1
			}
			mapper[string(lineParts[1])] = string(bytes.Trim(value, "\r\n"))
		}
		return mapper, nil
	case cmd == "stats":
		mapper := make(map[string]string)
		for {
			line, err := byteReader.ReadBytes('\n')
			if err != nil {
				return mapper, err
			}
			line = bytes.Trim(line, "\r\n")
			if string(line) == "END" {
				return mapper, nil
			}
			lineParts := bytes.Split(line, []byte(" "))
			if len(lineParts) != 3 || string(lineParts[0]) != "STAT" {
				err = c.checkError(line)
				if err != nil {
					return mapper, err
				}
				return mapper, errors.New("响应数据格式非法")
			}
			mapper[string(lineParts[1])] = string(lineParts[2])
		}
		return mapper, nil
	default:
		line, err := byteReader.ReadBytes('\n')
		if err != nil {
			return line, err
		}
		line = bytes.Trim(line, "\r\n")
		err = c.checkError(line)
		if err != nil {
			return nil, err
		}
		return line, nil
	}
}
