package textbuffer

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

/*
func TestBuffer1(t *testing.T) {
	startRecycler()
	reader := strings.NewReader("hello world! the quick brown fox jumps over the lazy dog.")
	cache := NewCache(reader, 3)
	bytebuf := make([]byte, 3)
	_, err := cache.ReadAt(bytebuf, 4)
	if err != nil {
		t.Fail()
	}
	fmt.Printf("read \"%s\" from text cache\n", string(bytebuf))
	cache.Dump()
	time.Sleep(5 * time.Second)
	cache.releaseBuffer(0)
	cache.Dump()
	time.Sleep(5 * time.Second)
	stopRecycler()
	time.Sleep(time.Second)
}
*/

func TestBuffer2(t *testing.T) {
	startRecycler()
	reader := strings.NewReader("hello world! the quick brown fox jumps over the lazy dog.")
	cache := NewCache(reader, 3)
	bytebuf := make([]byte, 22)
	n, err := cache.ReadAt(bytebuf, 6)
	if err != nil {
		t.Fail()
	}
	fmt.Printf("read %d bytes = \"%s\" from text cache\n", n, string(bytebuf))
	if n != len(bytebuf) {
		t.Fail()
	}
	cache.Dump()
	stopRecycler()
	time.Sleep(time.Second)
}

func TestBuffer3(t *testing.T) {
	startRecycler()
	reader := strings.NewReader("hello world! the quick brown fox jumps over the lazy dog.")
	cache := NewCache(reader, 2)
	bytebuf := make([]byte, 22)
	n, err := cache.ReadAt(bytebuf, 6)
	if err != nil {
		t.Fail()
	}
	fmt.Printf("read %d bytes = \"%s\" from text cache\n", n, string(bytebuf))
	if n != len(bytebuf) {
		t.Fail()
	}
	cache.Dump()
	stopRecycler()
	time.Sleep(time.Second)
}
