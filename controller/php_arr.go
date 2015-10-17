package main

import (
		"os"
		"log"
		"fmt"
		"bufio"
		"strings"
		"strconv"
		// "unicode/utf8"
)

func ptab(tab int) {
	for i := 0; i < tab; i++ {
		fmt.Printf("\t")
	}
}
type ArrayItem struct {
	key interface{}
	value interface{}
}
func parseArrayItem(reader *bufio.Reader) (ai ArrayItem) {
	ai.key = parse(reader)
	ai.value = parse(reader)
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
func parseInt(reader *bufio.Reader) int64 {
	raw, err := reader.ReadString(';')
	if err != nil {
		log.Fatal(err)
	}
	i, err := strconv.ParseInt(strings.TrimSuffix(raw, ";"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("read int %d\n", i)
	return i
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
func parse(reader *bufio.Reader) (i interface{}) {
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
		case 'a':
			i = parseArray(reader)
		case 's':
			i = parseString(reader)
		case 'i':
			i = parseInt(reader)
		// case 'b':
		default:
			log.Fatal("unknow type ", t)
		}
	}
	return i
}
func main() {
	file, err := os.Open("serialize.txt")
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)

	fmt.Printf("%q", parse(reader))
}
