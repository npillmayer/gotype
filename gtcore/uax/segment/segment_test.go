package segment

import (
	"fmt"
	"strings"
	"testing"
)

//var CT tracing.Trace = tracing.CoreTracer

func TestWhitespace1(t *testing.T) {
	seg := NewSegmenter()
	seg.Init(strings.NewReader("Hello World!"))
	for seg.Next() {
		t.Logf("segment = '%s' with p = %d", seg.Text(), seg.Penalties()[0])
	}
}

func TestWhitespace2(t *testing.T) {
	seg := NewSegmenter()
	seg.Init(strings.NewReader("	for (i=0; i<5; i++)   count += i;"))
	for seg.Next() {
		t.Logf("segment = '%s' with p = %d", seg.Text(), seg.Penalties()[0])
	}
}

func ExampleSegmenter() {
	seg := NewSegmenter()
	seg.Init(strings.NewReader("Hello World!"))
	for seg.Next() {
		fmt.Printf("segment = '%s' with penalty = %d\n", seg.Text(), seg.Penalties()[0])
	}
	// Output:
	// segment = 'Hello' with penalty = 100
	// segment = ' ' with penalty = -100
	// segment = 'World!' with penalty = -1000
}
