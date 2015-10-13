package memcached

import (
	"os"
	"testing"
)

var m Memcached = Memcached{}

func Test_FlushAll(t *testing.T) {
	err := m.FlushAll()
	if err == nil {
		t.Log("Success!")
	} else {
		t.Error("Failure!", err)
	}
}

func Test_Set(t *testing.T) {
	err := m.Set(map[string]interface{}{"key": "hello", "value": "world"})
	if err == nil {
		t.Log("Success!")
	} else {
		t.Error("Failure!", err)
	}
	err = m.Set(map[string]interface{}{"key": "xiayf", "value": "youngsterxyf", "expire_time": 3600})
	if err == nil {
		t.Log("Success!")
	} else {
		t.Error("Failure!", err)
	}
}

func Test_Get(t *testing.T) {
	value, err := m.Get("hello")
	if err != nil {
		t.Error("Failure!", err)
	} else {
		t.Log("Success!", value)
	}
}

func Test_Gets(t *testing.T) {
	mapper, err := m.Gets("hello", "xiayf")
	if err != nil {
		t.Error("Failure!", err)
	} else {
		for k, v := range mapper {
			t.Log(k, " --> ", v)
		}
	}

}

func Test_Delete(t *testing.T) {
	err := m.Delete("hello")
	if err == nil {
		t.Log("Success!")
	} else {
		t.Error("failure!", err)
	}
}

func Test_Stats(t *testing.T) {
	mapper, err := m.Stats()
	if err != nil {
		t.Error("stats failure!", err)
	} else {
		for k, v := range mapper {
			t.Log(k, " --> ", v)
		}
	}
	mapper, err = m.Stats("slabs")
	if err != nil {
		t.Error("stats slabs failure!", err)
	} else {
		for k, v := range mapper {
			t.Log(k, " --> ", v)
		}
	}
	mapper, err = m.Stats("items")
	if err != nil {
		t.Error("stats items failure!", err)
	} else {
		for k, v := range mapper {
			t.Log(k, " --> ", v)
		}
	}
}

func TestMain(tm *testing.M) {
	m.New("127.0.0.1", 11211)
	exitCode := tm.Run()
	m.Close()
	os.Exit(exitCode)
}
