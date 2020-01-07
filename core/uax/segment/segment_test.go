package segment

import (
	"fmt"
	"strings"
	"testing"

	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
)

func TestWhitespace1(t *testing.T) {
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	seg := NewSegmenter()
	seg.Init(strings.NewReader("Hello World!"))
	for seg.Next() {
		t.Logf("segment = '%s' with p = %d", seg.Text(), seg.Penalties()[0])
	}
}

func TestWhitespace2(t *testing.T) {
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	seg := NewSegmenter()
	seg.Init(strings.NewReader("	for (i=0; i<5; i++)   count += i;"))
	for seg.Next() {
		t.Logf("segment = '%s' with p = %d", seg.Text(), seg.Penalties()[0])
	}
}

func ExampleSegmenter() {
	seg := NewSegmenter() // will use a SimpleWordBreaker
	seg.Init(strings.NewReader("Hello World!"))
	for seg.Next() {
		fmt.Printf("segment: penalty = %5d for breaking after '%s'\n", seg.Penalties()[0], seg.Text())
	}
	// Output:
	// segment: penalty =   100 for breaking after 'Hello'
	// segment: penalty =  -100 for breaking after ' '
	// segment: penalty = -1000 for breaking after 'World!'
}
