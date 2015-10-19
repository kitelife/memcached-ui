package phpunserialize

import (
	"os"
	"log"
	// "fmt"
	"bufio"
	"testing"
)

func Assert(real interface{}, should interface{}) {
	if should != real {
		log.Fatal("not good")
	}
}
func TestParse(t *testing.T) {
	file, err := os.Open("serialize.txt")
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)

	arr := Parse(reader)
	// t.Log(a)
	// fmt.Println(arr)
	a, ok := arr.([] interface {})
	Assert(ok, true)

	Assert(len(a), 2)
	_, ok = a[0].(string)
	Assert(ok, true)
	Assert(a[1], nil)

	file, err = os.Open("map.txt")
	if err != nil {
		log.Fatal(err)
	}
	reader = bufio.NewReader(file)
	u := Parse(reader)
	m, ok := u.(map[string] interface {})
	Assert(ok, true)
	Assert(len(m), 1)
	Assert(m["a"], "b")

    // t.Fail()
}
