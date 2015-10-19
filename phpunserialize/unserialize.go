// write by: xiaochi

// Package phpunserialize make it possible to unserialize PHP data
// to Go.
// It support null, bool, float, int, string and array.
// Objects are not supported yet.
//
// Usage: str = Parse(*reader)
package phpunserialize

import (
		"log"
		"fmt"
		"bufio"
		"strings"
		"strconv"
)

type ArrayItem struct {
	key interface{}
	value interface{}
}
func parseArrayItem(reader *bufio.Reader) (ai ArrayItem) {
	ai.key = Parse(reader)
	ai.value = Parse(reader)
	fmt.Printf("%q => %q\n", ai.key, ai.value)
	return ai
}
func parseArrayBody(reader *bufio.Reader, arraylen uint64) (res interface{}) {
	_, err := reader.ReadString('{')
	if err != nil {
		log.Fatal(err)
	}
	item := parseArrayItem(reader)
	var t interface{}
	t = item.key
	var arr []interface{}
	m := map[string]interface{}{}
	switch t := t.(type) {
		default:
			fmt.Printf("unexpected type %T\n", t)     // %T prints whatever type t has
			log.Fatal("unexpected type ", t)
		case int64:
			fmt.Printf("int64 %t\n", t)             // t has type int64
			if (int(item.key.(int64)) != 0) {
				log.Fatal("we do not support array not start with 0, but ", item.key)
			}
			arr = append(arr, item.value)
			for i:=1; i < int(arraylen); i++ {
				item = parseArrayItem(reader)
				if int(item.key.(int64)) != i {
					log.Fatal("we do not support that type array ", item.key, i)
				}
				arr = append(arr, item.value)
			}
			return arr
		case string:
			fmt.Printf("string %d\n", t)             // t has type string
			m[item.key.(string)] = item.value
			for i:=1; i < int(arraylen); i++ {
				item = parseArrayItem(reader)
				k := item.key.(string)
				m[k] = item.value
			}
			return m
	}

	_, err = reader.ReadString('}')
	if err != nil {
		log.Fatal(err)
	}
	return res
}
func parseLen(reader *bufio.Reader) uint64 {
	rawlen, err := reader.ReadString(':')
	if err != nil {
		log.Fatal(err)
	}
	pure := strings.TrimSuffix(rawlen, ":")
	ilen, err := strconv.ParseUint(pure, 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	return ilen
}
func parseArray(reader *bufio.Reader) (interface{}) {
	arraylen := parseLen(reader)
	var res []interface{}
	if arraylen == 0 {
		return res
	}
	return parseArrayBody(reader, arraylen)
}
func parseString(reader *bufio.Reader) string {
	strlen := parseLen(reader)
	// fmt.Printf("prepare to read string(%d)\n", strlen)
	ReadByteEnsure(reader, '"')

	var p [1000]byte
	n, err := reader.Read(p[:strlen])
	if err != nil {
		log.Fatal(err)
	}
	if n != int(strlen) {
		log.Fatal("want to read ", strlen, ", but read ", n)
	}
	// fmt.Printf("%s\n", string(p[:strlen]))
	ReadByteEnsure(reader, '"')
	ReadByteEnsure(reader, ';')
	return string(p[0:strlen])
}
func ReadStringTo(reader *bufio.Reader) string {
	raw, err := reader.ReadString(';')
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSuffix(raw, ";")
}
func parseInt(reader *bufio.Reader) int64 {
	i, err := strconv.ParseInt(ReadStringTo(reader), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("read int %d\n", i)
	return i
}
func parseFloat(reader *bufio.Reader) float64 {
	f, err := strconv.ParseFloat(ReadStringTo(reader), 64)
	if err != nil {
		log.Fatal(err)
	}
	return f
}
func parseBool(reader *bufio.Reader) bool {
	i, err := strconv.ParseUint(ReadStringTo(reader), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	if i == 0 {
		return false
	}
	if i == 1 {
		return true
	}
	log.Fatal("bool can not be ", i)
	return false
}
func ReadByteEnsure(reader *bufio.Reader, c byte) (error) {
	// fmt.Printf("must be '%c'\n", c)
	s, err := reader.ReadByte()
	if err != nil {
		return err
	}
	if s != c {
		log.Fatal("not '", string(c), "', but ", string(s))
	}
	return nil;
}
func Parse(reader *bufio.Reader) (i interface{}) {
	t, err := reader.ReadByte()
	if err != nil {
		log.Fatal(err)
	}
	if t == 'N' {
		fmt.Printf("prepare to read Null\n")
		err = ReadByteEnsure(reader, ';')
		if err != nil {
			log.Fatal(err)
		}
		return nil
	} else {
		fmt.Printf("type is %c\n", t)
		err = ReadByteEnsure(reader, ':')
		if err != nil {
			log.Fatal(err)
		}
		switch t {
		case 'd':
			i = parseFloat(reader)
		case 'a':
			i = parseArray(reader)
		case 's':
			i = parseString(reader)
		case 'i':
			i = parseInt(reader)
		case 'b':
			i = parseBool(reader)
		default:
			log.Fatal("unknow type ", t)
		}
	}
	return i
}
