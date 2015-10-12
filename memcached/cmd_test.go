package memcached

import (
	"os"
	"testing"
)

var m Memcached = Memcached{}

func Test_FlushAll(t *testing.T) {
	status := m.FlushAll()
	if status {
		t.Log("Success!")
	} else {
		t.Error("Failure!")
	}
}

func Test_Get(t *testing.T) {
	value := m.Get("hello")
	if value == nil {
		t.Error("Failure!")
	} else {
		t.Log("Success!")
	}
}

func Test_Set(t *testing.T) {
	status := m.Set("hello", "world", 0)
	if status {
		t.Log("Success!")
	} else {
		t.Error("Failure!")
	}
}

func Test_Stats(t *testing.T) {
	result := m.Stats()
	if result == nil {
		t.Error("stats failure!")
	} else {
		t.Log("stats Success!")
	}
	result = m.Stats("slabs")
	if result == nil {
		t.Error("stats slabs failure!")
	} else {
		t.Log("stats slabs Success!")
	}
	result = m.Stats("items")
	if result == nil {
		t.Error("stats items failure!")
	} else {
		t.Log("stats items Success!")
	}
}

func TestMain(tm *testing.M) {
	m.New("127.0.0.1", 11211)
	exitCode := tm.Run()
	m.Close()
	os.Exit(exitCode)
}
