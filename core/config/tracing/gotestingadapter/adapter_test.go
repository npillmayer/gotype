package gotestingadapter_test

import (
	"testing"

	"github.com/npillmayer/gotype/core/config"
	"github.com/npillmayer/gotype/core/config/testadapter"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/gotestingadapter"
)

func Test1(t *testing.T) {
	l := gotestingadapter.New()
	l.SetTraceLevel(tracing.LevelDebug)
	l.Debugf("Hello 1")
	l.P("a", "b").Infof("World")
	l.Debugf("Hello 2")
	l.SetTraceLevel(tracing.LevelError)
	l.Debugf("Hello 3")
	// will produce no output
}

func Test2(t *testing.T) {
	config.AddTraceAdapter("test", gotestingadapter.GetAdapter())
	c := testadapter.New()
	c.Set("tracing", "test")
	config.Initialize(c)
	teardown := gotestingadapter.RedirectTracing(t)
	defer teardown()
	tracing.EngineTracer.P("key", "value").Errorf("This is a test")
	// output only seen with -v flag turned on
}
