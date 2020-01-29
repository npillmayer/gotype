package scanner

import (
	"fmt"
	"strings"
	"testing"

	"github.com/npillmayer/gotype/core/config/gtrace"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
)

var inputStrings = []string{
	"1",
	"1+1",
	"Hello World",
	"x=\"mystring\" // commented ",
}

var tokenCounts = []int{1, 3, 2, 3}

func TestScan1(t *testing.T) {
	gtrace.SyntaxTracer = gotestingadapter.New()
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	for i, input := range inputStrings {
		t.Logf("------+-----------------+--------")
		reader := strings.NewReader(input)
		name := fmt.Sprintf("input #%d", i)
		scanner := GoTokenizer(name, reader)
		tokval, token, pos, _ := scanner.NextToken(AnyToken)
		count := 0
		for tokval != EOF {
			t.Logf(" %4d | %15s | @%5d", tokval, Lexeme(token), pos)
			tokval, token, pos, _ = scanner.NextToken(AnyToken)
			count++
		}
		if count != tokenCounts[i] {
			t.Errorf("Expected token count for #%d to be %d, is %d", i, tokenCounts[i], count)
		}
	}
	t.Logf("------+-----------------+--------")
}
