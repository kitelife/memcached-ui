package memcached

import (
	"fmt"
)

const (
	SET_FLAGS = 123456
)

type Memcached struct {
	conn Connection
}

func (m *Memcached) New(host string, port int) {
	m.conn = Connection{
		Host: host,
		Port: port,
	}
}

func (m *Memcached) FlushAll() bool {
	cmd := "flush_all\r\n"
	resp, err := m.conn.Send(cmd)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Printf("response: %s\n", string(resp))
	return true
}

func (m *Memcached) Set(key, value string, expTime int) bool {
	cmd := fmt.Sprintf("set %s %d %d %d\r\n", key, SET_FLAGS, expTime, len(value))
	resp, err := m.conn.Send(cmd, fmt.Sprintf("%s\r\n", value))
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Printf("response: %s", string(resp))
	return true
}

func (m *Memcached) Get(key string) interface{} {
	cmd := fmt.Sprintf("get %s\r\n", key)
	resp, err := m.conn.Send(cmd)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Printf("response: %s\n", string(resp))
	return string(resp)
}

func (m *Memcached) Stats(args ...string) interface{} {
	var cmd string
	if len(args) == 0 {
		cmd = "stats\r\n"
	} else {
		cmd = fmt.Sprintf("stats %s\r\n", args[0])
	}
	resp, err := m.conn.Send(cmd)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Printf("response: %s\n", string(resp))
	return string(resp)
}

func (m *Memcached) Close() {
	err := m.conn.Conn.Close()
	if err != nil {
		fmt.Println(err)
	}
}
