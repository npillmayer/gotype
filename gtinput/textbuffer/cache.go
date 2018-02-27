/*
Package textbuffer provides (as a first rough prototype) handling of cached
byte buffers.

Clients create a *TextBuffer, which caches portions of bytes of an io.ReaderAt.
The text buffer implements io.ReaderAt itself, thus may be used in the same
fashion as the underlying reader.

    reader := strings.NewReader("The quick brown fox jumps over the lazy dog.")
    cache := textbuffer.NewCache(reader, 3)
    bytebuf := make([]byte, 10)
    n, err := cache.ReadAt(bytebuf, 4) // reads "quick..."

*/
package textbuffer

import (
	"container/list"
	"fmt"
	"io"
	"time"
)

// Note: this is just a quick hack / first draft
// TODO polish this up
// TODO use a priority queue instead of a slice of buffers

const BufferSize = 10
const DefaultCacheBufferCount = 3

//const RecycleTimeout = time.Minute
const RecycleTimeout = 2 * time.Second

var makes int
var frees int
var get, trash chan *TextBuffer

type TextBuffer struct {
	lastAccess time.Time
	offset     int64
	buffer     []byte
}

// Create a text buffer. Will use DefaultBufferSize, if buffersize is 0.
func NewBuffer() *TextBuffer {
	buf := make([]byte, BufferSize)
	buf[0] = 0
	textbuf := &TextBuffer{lastAccess: time.Now(), offset: -1, buffer: buf}
	return textbuf
}

func (tbuf *TextBuffer) String() string {
	since := time.Since(tbuf.lastAccess) / time.Millisecond
	if len(tbuf.buffer) == 0 || tbuf.buffer[0] == 0 {
		return fmt.Sprintf("<%dms %d|empty>", since, tbuf.offset)
	} else {
		return fmt.Sprintf("<%dms %d|\"%s\">", since, tbuf.offset, debugPrefix(string(tbuf.buffer)))
	}
}

type queued struct {
	when    time.Time
	textbuf *TextBuffer
}

func startRecycler() {
	get = make(chan *TextBuffer)
	trash = make(chan *TextBuffer)
	go func() {
		q := new(list.List)
		running := true
		for running {
			if q.Len() == 0 {
				q.PushFront(queued{when: time.Now(), textbuf: NewBuffer()})
				fmt.Println("creating new text buffer")
			}
			e := q.Front()
			timeout := time.NewTimer(RecycleTimeout)
			fmt.Println("--------- setting timeout ------------")
			select {
			case b := <-trash:
				if b == nil {
					fmt.Println("=== STOP RECYCLER =========")
					running = false
					break
				}
				fmt.Println("adding buffer back to free list")
				timeout.Stop()
				q.PushFront(queued{when: time.Now(), textbuf: b})
			case get <- e.Value.(queued).textbuf:
				timeout.Stop()
				q.Remove(e)
			case <-timeout.C:
				fmt.Printf("timeout expired, Q-len = %d\n", q.Len())
				e := q.Front().Next() // skip 1st one (which is always present)
				for e != nil {
					n := e.Next()
					if time.Since(e.Value.(queued).when) > RecycleTimeout {
						q.Remove(e)
						tbuf := e.Value.(queued).textbuf
						fmt.Printf("releaseing buffer = %v for GC\n", tbuf)
						e.Value = nil
					}
					e = n
				}
				fmt.Printf("             now Q-len = %d\n", q.Len())
			}
			//time.Sleep(3 * time.Second) // for debugging
		}

	}()
	return
}

func stopRecycler() {
	trash <- nil
}

func offsetBoundary(offset int64) int64 {
	return (offset / BufferSize) * BufferSize
}

// === TextCache =============================================================

type TextCache struct {
	textlen int64
	buffers []*TextBuffer
	reader  io.ReaderAt
}

func NewCache(reader io.ReaderAt, buffercount int) *TextCache {
	if buffercount == 0 {
		buffercount = DefaultCacheBufferCount
	}
	cache := &TextCache{}
	cache.buffers = make([]*TextBuffer, buffercount)
	cache.reader = reader
	return cache
}

func (cache *TextCache) Dump() {
	for i, buf := range cache.buffers {
		if buf != nil {
			fmt.Printf("[%2d] \"%v\"\n", i, string(buf.buffer))
		} else {
			fmt.Printf("[%2d] \"\"\n", i)
		}
	}
}

func (cache *TextCache) ReadAt(buffer []byte, offset int64) (N int, err error) {
	L := len(buffer) // target count, will be decreased incrementally
	buf, inx, n := cache.loadText(offset)
	start := int(offset - offsetBoundary(offset))
	fmt.Printf("reading pos %d in %d|\"%s\"\n", start, inx, string(buf.buffer))
	avail := n - start
	fmt.Printf("L = %d, start = %d, avail = %d\n", L, start, avail)
	if avail >= L { // can fill in one go, should be frequent case
		slice := buf.buffer[start : start+L]
		copy(buffer, slice)
	} else { // fill buffer incrementally
		slice := buf.buffer[start:]
		copy(buffer, slice)
		fmt.Printf("(1) now buffer = '%s'\n", buffer)
		L -= avail
		N += avail
		for L > 0 {
			buf, inx, n = cache.loadText(offset + int64(N))
			avail = n
			fmt.Printf("L = %d, N = %d, avail = %d\n", L, N, avail)
			if avail >= L { // can fill in one go
				slice = buf.buffer[:L]
				fmt.Printf("copy '%s' -> '%s'\n", buffer[N:], slice)
				copy(buffer[N:], slice)
				fmt.Printf("(3) now buffer = '%s'\n", buffer)
				N += L
				L = 0 // done
			} else {
				fmt.Printf("copy '%s' -> '%s'\n", buffer[N:], buf.buffer)
				copy(buffer[N:], buf.buffer)
				fmt.Printf("(2) now buffer = '%s'\n", buffer)
				L -= avail
				N += avail
			}
		}
	}
	return
}

func (cache *TextCache) loadText(offset int64) (buf *TextBuffer, inx int, n int) {
	buf, inx = cache.findBufferForOffset(offset)
	if inx == -1 {
		boundary := offsetBoundary(offset)
		buf, inx = cache.getFreeBuffer()
		var err error
		n, err = cache.reader.ReadAt(buf.buffer, boundary)
		fmt.Printf("read from reader: %d bytes|\"%s\"\n", n, debugPrefix(string(buf.buffer)))
		if err != nil {
			panic(fmt.Sprint("error loading text, %d bytes read: %v", n, err.Error()))
		}
	}
	buf.lastAccess = time.Now()
	return
}

func (cache *TextCache) getFreeBuffer() (*TextBuffer, int) {
	for inx, buf := range cache.buffers {
		if buf == nil {
			buf = <-get
			buf.offset = -1
			cache.buffers[inx] = buf
			fmt.Println("get buffer from recycler")
			return buf, inx
		}
	}
	buf, inx := cache.findBufferLRU()
	buf.offset = -1
	return buf, inx
}

func (cache *TextCache) findBufferLRU() (*TextBuffer, int) {
	var index int = -1
	recent := time.Now()
	for inx, b := range cache.buffers {
		if b != nil {
			if b.lastAccess.Before(recent) {
				recent = b.lastAccess
				index = inx
			}
		}
	}
	if index >= 0 {
		fmt.Printf("LRU buffer = #%d\n", index)
		return cache.buffers[index], index
	}
	return nil, -1
}

func (cache *TextCache) releaseBuffer(inx int) {
	buf := cache.buffers[inx]
	if buf != nil {
		fmt.Printf("release buffer #%d|%s\n", inx, buf.String())
		cache.buffers[inx] = nil
		trash <- buf
	}
}

func (cache *TextCache) findBufferForOffset(offset int64) (*TextBuffer, int) {
	boundary := offsetBoundary(offset)
	for inx, buf := range cache.buffers {
		if buf != nil && buf.offset == boundary {
			return buf, inx
		}
	}
	return nil, -1
}

/*
func main() {
	pool := make([][]byte, 20)
	get, trash := makeRecycler()
	var m runtime.MemStats
	for {
		b := <-get
		i := rand.Intn(len(pool))
		if pool[i] != nil {
			trash <- pool[i]
			frees++
		}
		pool[i] = b
		time.Sleep(time.Second)
		bytes := 0
		for i := 0; i < len(pool); i++ {
			if pool[i] != nil {
				bytes += len(pool[i])
			}
		}
		runtime.ReadMemStats(&m)
		fmt.Printf("%d,%d,%d,%d,%d,%d,%d\n", m.HeapSys, bytes, m.HeapAlloc,
			m.HeapIdle, m.HeapReleased, makes, frees)
	}

}
*/

func debugPrefix(s string) string {
	if len(s) > 8 {
		return fmt.Sprintf("%.5s...", s)
	} else {
		return s
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
