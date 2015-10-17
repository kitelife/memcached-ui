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
func parseArrayBody(reader *bufio.Reader, arraylen uint64) (res []interface{}) {
	_, err := reader.ReadString('{')
	if err != nil {
		log.Fatal(err)
	}
	for i:=0; i < int(arraylen); i++ {
		ArrayItem := parseArrayItem(reader)
		res = append(res, ArrayItem)
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
func parseArray(reader *bufio.Reader) []interface{} {
	arraylen := parseLen(reader)
	return parseArrayBody(reader, arraylen)
}
func parseString(reader *bufio.Reader) string {
	strlen := parseLen(reader)
	fmt.Printf("prepare to read string(%d)\n", strlen)
	ReadByteEnsure(reader, '"')

	var p [1000]byte
	n, err := reader.Read(p[:strlen])
	if err != nil {
		log.Fatal(err)
	}
	if n != int(strlen) {
		log.Fatal("want to read ", strlen, ", but read ", n)
	}
	fmt.Printf("%s\n", string(p[:strlen]))
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
	fmt.Printf("must be '%c'\n", c)
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
