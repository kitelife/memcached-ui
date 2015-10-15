package memcached

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// https://github.com/memcached/memcached/blob/master/doc/protocol.txt

const (
	// magic number，目前没啥鸟用
	SET_FLAGS = 123456
)

type Memcached struct {
	conn Connection
}

func (m *Memcached) New(host string, port int) error {
	m.conn = Connection{
		Host: host,
		Port: port,
	}
	return m.conn.Open()
}

/*
存储类型命令：set、add、replace、append、prepend、cas
*/

type StorageCmdArgStruct map[string]interface{}

func (m *Memcached) runStorageCmd(cmdName string, args StorageCmdArgStruct) ([]byte, error) {
	// 必须
	var key, value string
	keyI, ok := args["key"]
	if ok == false {
		return nil, errors.New("缺少参数key")
	} else {
		key = keyI.(string)
	}
	valueI, ok := args["value"]
	if ok == false {
		return nil, errors.New("缺少参数value")
	} else {
		value = valueI.(string)
	}

	// 可选
	var flags, expTime string
	flagsI, ok := args["flags"]
	if ok == false {
		flags = strconv.Itoa(SET_FLAGS)
	} else {
		flags = string(flagsI.(int))
	}
	expTimeI, ok := args["expire_time"]
	if ok == false {
		expTime = "0"
	} else {
		expTime = strconv.Itoa(expTimeI.(int))
	}
	argList := []string{key, flags, expTime, strconv.Itoa(len(value))}
	if cmdName == "cas" {
		if casUnique, ok := args["cas_unique"]; ok {
			argList = append(argList, casUnique.(string))
		}
	}

	cmd := fmt.Sprintf("%s %s\r\n", cmdName, strings.Join(argList, " "))
	err := m.conn.Send(cmd, fmt.Sprintf("%s\r\n", value))
	if err != nil {
		return nil, err
	}
	resp, err := m.conn.Receive(cmdName)
	return resp.([]byte), err
}

func (m *Memcached) Set(args StorageCmdArgStruct) ([]byte, error) {
	return m.runStorageCmd("set", args)
}

func (m *Memcached) Add(args StorageCmdArgStruct) ([]byte, error) {
	return m.runStorageCmd("add", args)
}

func (m *Memcached) Replace(args StorageCmdArgStruct) ([]byte, error) {
	return m.runStorageCmd("replace", args)
}

func (m *Memcached) Append(args StorageCmdArgStruct) ([]byte, error) {
	return m.runStorageCmd("append", args)
}

func (m *Memcached) Prepend(args StorageCmdArgStruct) ([]byte, error) {
	return m.runStorageCmd("prepend", args)
}

func (m *Memcached) Cas(args StorageCmdArgStruct) ([]byte, error) {
	return m.runStorageCmd("cas", args)
}

/*
数据获取类型命令：get、gets
*/

func (m *Memcached) runFetchCmd(cmdName, keys string) (map[string]string, error) {
	cmd := fmt.Sprintf("%s %s\r\n", cmdName, keys)
	err := m.conn.Send(cmd)
	if err != nil {
		return nil, err
	}
	resp, err := m.conn.Receive(cmdName)
	return resp.(map[string]string), err
}

func (m *Memcached) Get(key string) (string, error) {
	resp, err := m.runFetchCmd("get", key)
	if err != nil {
		return "", err
	}
	return resp[key], nil
}

func (m *Memcached) Gets(keys ...string) (map[string]string, error) {
	resp, err := m.runFetchCmd("gets", strings.Join(keys, " "))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

/*
其他命令：flush_all、delete、incr、decr、touch、stats
*/

func (m *Memcached) FlushAll() ([]byte, error) {
	cmd := "flush_all\r\n"
	err := m.conn.Send(cmd)
	if err != nil {
		return nil, err
	}
	resp, err := m.conn.Receive("flush_all")
	return resp.([]byte), err
}

func (m *Memcached) Delete(key string) ([]byte, error) {
	cmd := fmt.Sprintf("delete %s\r\n", key)
	err := m.conn.Send(cmd)
	if err != nil {
		return nil, err
	}
	resp, err := m.conn.Receive("delete")
	return resp.([]byte), err
}

func (m *Memcached) Incr(key string, value int64) ([]byte, error) {
	cmd := fmt.Sprintf("incr %s %d\r\n", key, value)
	err := m.conn.Send(cmd)
	if err != nil {
		return nil, err
	}
	resp, err := m.conn.Receive("incr")
	return resp.([]byte), err
}

func (m *Memcached) Decr(key string, value int64) ([]byte, error) {
	cmd := fmt.Sprintf("decr %s %d\r\n", key, value)
	err := m.conn.Send(cmd)
	if err != nil {
		return nil, err
	}
	resp, err := m.conn.Receive("decr")
	return resp.([]byte), err
}

func (m *Memcached) Touch(key string, expTime int) ([]byte, error) {
	cmd := fmt.Sprintf("touch %s %d\r\n", key, expTime)
	err := m.conn.Send(cmd)
	if err != nil {
		return nil, err
	}
	resp, err := m.conn.Receive("touch")
	return resp.([]byte), err
}

func (m *Memcached) Stats(args ...string) (map[string]string, error) {
	var cmd string
	if len(args) == 0 {
		cmd = "stats\r\n"
	} else {
		cmd = fmt.Sprintf("stats %s\r\n", args[0])
	}
	err := m.conn.Send(cmd)
	if err != nil {
		return nil, err
	}
	resp, err := m.conn.Receive("stats")
	return resp.(map[string]string), err
}

// 关闭网络连接

func (m *Memcached) Close() {
	err := m.conn.Conn.Close()
	if err != nil {
		fmt.Println(err)
	}
}
