package phpunserialize

import (
	"os"
	"log"
	"fmt"
	"bufio"
	"testing"
)

func Assert(real interface{}, should interface{}) {
	if should != real {
		log.Fatal("should be ", should, ", but ", real)
	}
}

func ParseFile(filename string) interface{} {
	fmt.Println("===== TEST %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)

	return Parse(reader)
}

func TestParse(t *testing.T) {
	arr := ParseFile("serialize.txt")
	// t.Log(a)
	// fmt.Println(arr)
	a, ok := arr.([] interface {})
	Assert(ok, true)

	Assert(len(a), 2)
	_, ok = a[0].(string)
	Assert(ok, true)
	Assert(a[1], nil)

	u := ParseFile("map.txt")
	m, ok := u.(map[string] interface {})
	Assert(ok, true)
	Assert(len(m), 1)
	Assert(m["a"], "b")

	u = ParseFile("nest.txt")
	a, ok = u.([] interface {})
	Assert(ok, true)
	Assert(len(a), 2)
	Assert(a[1], nil)
	// fmt.Println("%q", a[0])
	a, ok = a[0].([] interface {})
	Assert(ok, true)
	for _, v := range a {
		m, ok = v.(map[string] interface {})
		Assert(ok, true)
		_, ok = m["host"]
		Assert(ok, true)
		_, ok = m["type"]
		Assert(ok, true)
	}

    // t.Fail()
}
