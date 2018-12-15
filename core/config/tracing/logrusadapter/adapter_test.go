package logrusadapter_test

import (
	"testing"

	"github.com/npillmayer/gotype/core/config"
	"github.com/npillmayer/gotype/core/config/tracing"
	"github.com/npillmayer/gotype/core/config/tracing/logrusadapter"
)

func Test1(t *testing.T) {
	l := logrusadapter.New()
	l.SetTraceLevel(tracing.LevelDebug)
	l.Debugf("Hello 1")
	l.P("a", "b").Infof("World")
	l.Debugf("Hello 2")
	l.SetTraceLevel(tracing.LevelError)
	l.Debugf("Hello 3")
}

func Test2(t *testing.T) {
	config.InitDefaults()
	config.InitTracing(logrusadapter.GetAdapter())
	tracing.EngineTracer.P("key", "value").Errorf("This is a test")
}

func Test3(t *testing.T) {
	v := []int{1, 2, 3}
	tracing.With(tracing.EngineTracer).Dump("v", v)
}
